package dao

import (
	"common"
	"errors"
	"fmt"
	"github.com/astaxie/beego/orm"
	"github.com/golang/protobuf/proto"
	"pb"
	"strings"
	"utils/admin/dao"
	"utils/agent/models"
	otcmodels "utils/otc/models"
)

type AgentPathDao struct {
	common.BaseDao
}

func NewAgentPathDao(db string) *AgentPathDao {
	return &AgentPathDao{
		BaseDao: common.NewBaseDao(db),
	}
}

var AgentPathDaoEntity *AgentPathDao

// chip
type Chip struct {
}

// create new agent path
func (d *AgentPathDao) Create(uid uint64, parentInviteCode string) (err error) {
	// get parent data
	var (
		parentUid   uint64
		parentLevel uint32
		parentPath  string
		sql         string
	)

	if len(parentInviteCode) > 0 {
		parentInviteCode = strings.ToUpper(parentInviteCode)
		sql = fmt.Sprintf("SELECT %s, %s, %s FROM %s WHERE %s=?",
			models.COLUMN_AgentPath_Uid, models.COLUMN_AgentPath_Level, models.COLUMN_AgentPath_Path, models.TABLE_AgentPath, models.COLUMN_AgentPath_InviteCode)
		err = d.Orm.Raw(sql, parentInviteCode).QueryRow(&parentUid, &parentLevel, &parentPath)
		if err != nil {
			common.LogFuncError("%v", err)
			return
		}
	}

	var (
		level = parentLevel + 1
		sn    uint32
	)

	// get max index of the child level
	sql = fmt.Sprintf("SELECT max(%s) FROM %s WHERE %s=?",
		models.COLUMN_AgentPath_Sn, models.TABLE_AgentPath, models.COLUMN_AgentPath_Level)
	err = d.Orm.Raw(sql, level).QueryRow(&sn)
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}

	// generate
	now := common.NowInt64MS()
	var newChild = models.AgentPath{
		Uid:              uid,
		Level:            level,
		Sn:               sn + 1,
		Path:             d.generatePath(level, sn+1, parentPath),
		Ctime:            now,
		Mtime:            now,
		ParentUid:        parentUid,
		DividendPosition: 0,
	}

	// generate invite code,
	var (
		inviteCode string
		isRelease  bool
		id         uint32
	)
	id, inviteCode, err = InviteCodeDaoEntity.TryGetUnusedOne()
	if err != nil {
		return
	}

	defer func() {
		if isRelease { //if need to release, just release it.
			err = InviteCodeDaoEntity.Release(id)
			if err != nil {
				common.LogFuncError("%v", err)
			}
		}
	}()

	// generate failed.
	if len(inviteCode) <= 0 {
		err = errors.New("generate inviteCode failed")
		return
	}
	newChild.InviteCode = inviteCode
	_, err = d.Orm.Insert(&newChild)
	if err != nil {
		isRelease = true //need to release
		common.LogFuncError("%v", err)
		return
	}

	//set invite code used
	err = InviteCodeDaoEntity.UpdateStatus(id, INVITE_CODE_STATUS_USED)
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}

	// parent invite num increase
	_, err = d.Orm.Raw(fmt.Sprintf("UPDATE %s SET %s=%s+1 WHERE %s=?",
		models.TABLE_AgentPath,
		models.COLUMN_AgentPath_InviteNum,
		models.COLUMN_AgentPath_InviteNum,
		models.COLUMN_AgentPath_Uid), parentUid).Exec()
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}
	err = d.CreateRedisAgentPath(parentUid, uid, level)
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}

	return
}

// get user id by invite code
func (d *AgentPathDao) GetUidByInviteCode(inviteCode string) (uid uint64, err error) {
	var agentPath models.AgentPath
	err = d.Orm.QueryTable(models.TABLE_AgentPath).Filter(models.COLUMN_AgentPath_InviteCode, inviteCode).One(&agentPath, models.COLUMN_AgentPath_Uid)
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}
	uid = agentPath.Uid

	return
}

//// get invite code by user id
//func (d *AgentPathDao) GetInviteCode(uid uint64) (inviteCode string, err error) {
//	err = d.Orm.Raw(fmt.Sprintf("SELECT %s FROM %s WHERE %s=?",
//		models.COLUMN_AgentPath_InviteCode,
//		models.TABLE_AgentPath,
//		models.COLUMN_AgentPath_Uid), uid).QueryRow(&inviteCode)
//	if err != nil {
//		common.LogFuncError("%v", err)
//		return
//	}
//
//	return
//}

//// check whether invite code exists or no.
//func (d *AgentPathDao) checkInviteCode(inviteCode string) (isExist bool, err error) {
//	user := models.User{
//		InviteCode: inviteCode,
//	}
//	err = d.Orm.Read(user, models.COLUMN_User_InviteCode)
//	if err != nil {
//		if err != orm.ErrNoRows { //no exist
//			common.LogFuncError("%v", err)
//			return
//		}
//		err = nil
//		return
//	}
//
//	return
//}

func (d *AgentPathDao) generatePath(level, sn uint32, parentPath string) string {
	return fmt.Sprintf("%s.%.4x%.4x", parentPath, level, sn)
}

func (d *AgentPathDao) QueryTotal() (total int, err error) {
	err = d.Orm.Raw(fmt.Sprintf("SELECT count(*) FROM %s",
		models.TABLE_AgentPath)).QueryRow(&total)
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}

	return
}

func (d *AgentPathDao) Query(page, limit int) (list []models.AgentPath, err error) {
	_, err = d.Orm.Raw(fmt.Sprintf("SELECT * FROM %s LIMIT ? OFFSET ?", models.TABLE_AgentPath), limit, (page-1)*limit).QueryRows(&list)
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}

	return
}

type AgentInfo struct {
	Uid              string `json:"uid"`
	UMobile          string `json:"u_mobile"`
	Ctime            int64  `json:"ctime"`
	Mtime            int64  `json:"mtime"`
	InviteCode       string `json:"invite_code"`
	WhitelistId      uint32 `json:"whitelist_id"`
	DividendPosition uint32 `json:"dividend_position"`
	InviteNum        uint32 `json:"invite_num"`
	ParentUid        string `json:"parent_uid"`
	PMobile          string `json:"p_mobile"`
	SumSalary        int64  `orm:"column(sum_salary)" json:"sum_salary,omitempty"`
	SumCanWithdraw   int64  `orm:"column(sum_can_withdraw)" json:"sum_can_withdraw,omitempty"`
}

//分页获取代理信息
func (d *AgentPathDao) GetAgentPathByPage(page, limit, wId int, pUid uint64, inviteCode string) (total int64, agentInfo []AgentInfo, err error) {
	var qbTotal orm.QueryBuilder
	var qbQuery orm.QueryBuilder
	qbTotal, err = orm.NewQueryBuilder("mysql")
	if err != nil {
		common.LogFuncError("err:%v", err)
		return
	}
	qbQuery, err = orm.NewQueryBuilder("mysql")
	if err != nil {
		common.LogFuncError("err:%v", err)
		return
	}

	// 构建查询对象
	qbTotal.Select("Count(*) AS total").
		From(models.TABLE_AgentPath).
		Where("1 = 1")
	qbQuery.Select("T1."+models.COLUMN_AgentPath_Uid,
		"T1."+models.COLUMN_AgentPath_Ctime,
		"T1."+models.COLUMN_AgentPath_Mtime,
		"T1."+models.COLUMN_AgentPath_InviteCode,
		"T1."+models.COLUMN_AgentPath_WhiteListId,
		"T1."+models.COLUMN_AgentPath_InviteNum,
		"T1."+models.COLUMN_AgentPath_ParentUid,
		"T1."+models.COLUMN_AgentPath_DividendPosition,
		"T2."+models.COLUMN_Agent_SumCanWithdraw,
		"T2."+models.COLUMN_Agent_SumSalary,
		"T3."+otcmodels.COLUMN_User_Mobile+" AS u_mobile ",
		"T4."+otcmodels.COLUMN_User_Mobile+" AS p_mobile ").
		From("((").Select("*").From(models.TABLE_AgentPath).Where("1=1")
	var param []interface{}
	if wId > 0 {
		qbTotal.And(models.COLUMN_AgentPath_DividendPosition + "=?")
		qbQuery.And(models.COLUMN_AgentPath_DividendPosition + "=?")
		param = append(param, wId)
	}
	if pUid > 0 {
		qbTotal.And(models.COLUMN_AgentPath_ParentUid + "=?")
		qbQuery.And(models.COLUMN_AgentPath_ParentUid + "=?")
		param = append(param, pUid)
	}
	if inviteCode != "" {
		qbTotal.And(models.COLUMN_AgentPath_InviteCode + "=?")
		qbQuery.And(models.COLUMN_AgentPath_InviteCode + "=?")
		param = append(param, inviteCode)
	}
	sqlTotal := qbTotal.String()
	err = d.Orm.Raw(sqlTotal, param...).QueryRow(&total)
	if err != nil {
		common.LogFuncError("mysql error: %v", err)
		return
	}

	qbQuery.Limit(limit).Offset((page - 1) * limit)
	sqlQuery := qbQuery.String()
	sqlQuery = fmt.Sprintf("%s) AS T1 LEFT JOIN %s AS T2 ON T1.%s=T2.%s LEFT JOIN %s AS T3 ON T1.%s=T3.%s LEFT JOIN %s AS T4 ON T1.%s=T4.%s)",
		sqlQuery, models.TABLE_Agent, models.COLUMN_AgentPath_Uid, models.COLUMN_Agent_Uid,
		otcmodels.TABLE_User, models.COLUMN_AgentPath_Uid, otcmodels.COLUMN_User_Uid,
		otcmodels.TABLE_User, models.COLUMN_AgentPath_ParentUid, otcmodels.COLUMN_User_Uid)
	_, err = d.Orm.Raw(sqlQuery, param...).QueryRows(&agentInfo)
	if err != nil {
		common.LogFuncError("mysql error:%v", err)
		return
	}

	return
}

//设置用户代理档位
func (d *AgentPathDao) SetAgentWhiteList(uid uint64, wId uint32) (err error) {
	agentPath := &models.AgentPath{
		Uid: uid,
	}

	err = d.Orm.Read(agentPath, models.COLUMN_AgentPath_Uid)
	if err != nil {
		common.LogFuncError("mysql error:%v", err)
		return
	}

	if agentPath.WhiteListId > 0 && wId > 0 {
		err = dao.ErrParam
		return
	}

	agentPath.WhiteListId = wId
	_, err = d.Orm.Update(agentPath, models.COLUMN_AgentPath_WhiteListId)
	if err != nil {
		common.LogFuncError("mysql error:%v", err)
		return
	}

	return
}

//删除代理 admin.agent_white_list删除时调用
func (d *AgentPathDao) DelWhiteList(wid uint32) (err error) {
	sql := fmt.Sprintf("UPDATE %s SET %s=0 WHERE %s=?", models.TABLE_AgentPath,
		models.COLUMN_AgentPath_WhiteListId, models.COLUMN_AgentPath_WhiteListId)
	_, err = d.Orm.Raw(sql, wid).Exec()
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}
	return
}

//设置用户月分红档位
func (d *AgentPathDao) SetAgentDividendWhiteList(uid uint64, wId uint32) (err error) {
	agentPath := &models.AgentPath{
		Uid: uid,
	}

	err = d.Orm.Read(agentPath, models.COLUMN_AgentPath_Uid)
	if err != nil {
		common.LogFuncError("mysql error:%v", err)
		return
	}

	if agentPath.WhiteListId > 0 && wId > 0 {
		err = dao.ErrParam
		return
	}

	agentPath.DividendPosition = wId
	_, err = d.Orm.Update(agentPath, models.COLUMN_AgentPath_DividendPosition)
	if err != nil {
		common.LogFuncError("mysql error:%v", err)
		return
	}

	return
}

//删除代理 admin.agent_white_list删除时调用
func (d *AgentPathDao) DelDividendWhiteList(wid uint32) (err error) {
	sql := fmt.Sprintf("UPDATE %s SET %s=0 WHERE %s=?", models.TABLE_AgentPath,
		models.COLUMN_AgentPath_DividendPosition, models.COLUMN_AgentPath_DividendPosition)
	_, err = d.Orm.Raw(sql, wid).Exec()
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}
	return
}

// get user low levels
func (d *AgentPathDao) Info(uid uint64) (agentPath models.AgentPath, err error) {
	agentPath = models.AgentPath{
		Uid: uid,
	}

	err = d.Orm.Read(&agentPath, models.COLUMN_Agent_Uid)
	if err != nil {
		if err == orm.ErrNoRows {
			err = nil
			return
		}
		common.LogFuncError("mysql_err:%v", err)
		return
	}

	return
}

// get user low levels
func (d *AgentPathDao) GetUserByPath(path string) (paths models.AgentPath, err error) {
	err = d.Orm.Raw(fmt.Sprintf("SELECT * FROM %s WHERE %s = ?",
		models.TABLE_AgentPath, models.COLUMN_AgentPath_Path), path+"%").QueryRow(&paths)
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}

	return
}

// get user low levels
func (d *AgentPathDao) GetAllLowAgentUidByPath(path string) (uids []uint64, err error) {
	_, err = d.Orm.Raw(fmt.Sprintf("SELECT uid FROM %s WHERE %s LIKE ?",
		models.TABLE_AgentPath, models.COLUMN_AgentPath_Path), path+"%").QueryRows(&uids)
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}

	return
}

//设置用户代理档位
func (d *AgentPathDao) SetAgentDividendPosition(uid uint64, dID uint32) (err error) {
	agentPath := &models.AgentPath{
		Uid: uid,
	}

	err = d.Orm.Read(agentPath, models.COLUMN_AgentPath_Uid)
	if err != nil {
		common.LogFuncError("mysql error:%v", err)
		return
	}
	agentPath.DividendPosition = dID
	_, err = d.Orm.Update(agentPath, models.COLUMN_AgentPath_DividendPosition)
	if err != nil {
		common.LogFuncError("mysql error:%v", err)
		return
	}

	return
}

func (d *AgentPathDao) InsertForTest(uid uint64, path string, parentUid uint64, level uint32, inviteCode string) (err error) {
	agentPath := &models.AgentPath{
		Uid:        uid,
		Path:       path,
		Level:      level,
		ParentUid:  parentUid,
		InviteCode: inviteCode,
	}
	_, err = d.Orm.Insert(agentPath)
	if err != nil {
		common.LogFuncError("mysql error:%v", err)
		return
	}
	return
}

func (d *AgentPathDao) DeleteForTest(uid []uint64) (err error) {
	sql := fmt.Sprintf("delete from %s", models.TABLE_AgentPath)
	if len(uid) != 0 {
		sql = fmt.Sprintf("delete from %s where %s in ?", models.TABLE_AgentPath, models.COLUMN_Agent_Uid)
	}
	_, err = d.Orm.Raw(sql).Exec()
	return
}

type AgentPosition struct {
	Uid              uint64 `json:"uid"`
	Level            uint32 `json:"level"`
	Ctime            int64  `json:"ctime"`
	DividendPosition uint32 `json:"dividend_position"`
}

func (d *AgentPathDao) GetAgentByPosition(level uint32, position int32, page, limit int) (agentsPosition []*AgentPosition, err error) {
	querySql := fmt.Sprintf("select uid,level,ctime,dividend_position from %s where %s=? AND %s=? ORDER BY %s DESC LIMIT ?,?", models.TABLE_AgentPath, models.COLUMN_AgentPath_Level, models.COLUMN_AgentPath_DividendPosition, models.COLUMN_Agent_Ctime)
	agents := make([]*models.AgentPath, 0)
	_, err = d.Orm.Raw(querySql, level, position, (page-1)*limit, limit).QueryRows(&agents)
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}
	agentsPosition = make([]*AgentPosition, len(agents))
	for _, agent := range agents {
		position := &AgentPosition{
			Uid:              agent.Uid,
			Level:            agent.Level,
			Ctime:            agent.Ctime,
			DividendPosition: agent.DividendPosition,
		}
		agentsPosition = append(agentsPosition, position)
	}
	return
}

func (d *AgentPathDao) QueryByUIds(uIds []uint64) (paths []models.AgentPath, err error) {
	length := len(uIds)
	if length == 0 {
		return
	}
	totalPage := length / QueryPage
	if length%QueryPage > 0 {
		totalPage += 1
	}

	for i := 0; i < totalPage; i++ {
		end := (i + 1) * QueryPage
		if i == totalPage-1 {
			end = length
		}
		subUIds := uIds[i*QueryPage : end]
		var path []models.AgentPath
		_, err = d.Orm.QueryTable(models.TABLE_AgentPath).Filter(models.COLUMN_AgentPath_Uid+"__in", subUIds).All(&path)
		if err != nil {
			common.LogFuncError("error:%v", err)
			return
		}
		paths = append(paths, path...)
	}

	return
}

type MapWhiteList map[uint64]uint32

func (d *AgentPathDao) GetAllWhiteList() (mapWhiteList MapWhiteList, err error) {
	mapWhiteList = make(MapWhiteList)
	querySql := fmt.Sprintf("SELECT %s, %s FROM %s WHERE %s > 0", models.COLUMN_AgentPath_Uid,
		models.COLUMN_AgentPath_WhiteListId, models.TABLE_AgentPath, models.COLUMN_AgentPath_WhiteListId)

	var path []models.AgentPath
	_, err = d.Orm.Raw(querySql).QueryRows(&path)
	if err != nil {
		if err == orm.ErrNoRows {
			err = nil
			return
		}
		common.LogFuncError("error:%v", err)
		return
	}
	for _, item := range path {
		mapWhiteList[item.Uid] = item.WhiteListId
	}

	return
}
func (d *AgentPathDao) FindByUids(uids []uint64) (agentPathsMap map[uint64]*models.AgentPath, err error) {
	agentPaths := make([]*models.AgentPath, 0)
	qs := d.Orm.QueryTable(models.TABLE_AgentPath)
	_, err = qs.Filter(models.COLUMN_AgentPath_Uid+"__in", uids).All(&agentPaths)
	if err != nil {
		common.LogFuncError("data err %v", err)
	}
	agentPathsMap = make(map[uint64]*models.AgentPath, 0)
	for _, agentPath := range agentPaths {
		agentPathsMap[agentPath.Uid] = agentPath
	}
	return
}

func (d *AgentPathDao) GetByUid(uid uint64) (agentPath *models.AgentPath, err error) {
	querySql := fmt.Sprintf("select * from %s where %s=?", models.TABLE_AgentPath, models.COLUMN_AgentPath_Uid)
	err = d.Orm.Raw(querySql, uid).QueryRow(&agentPath)
	if err != nil {
		common.LogFuncError("mysql error:%v", err)
		return
	}
	return
}

func (d *AgentPathDao) GetSubLevelByUid(parentUid uint64) (agentPath []models.AgentPath, err error) {
	_, err = d.Orm.Raw(fmt.Sprintf("SELECT * FROM %s WHERE %s=?", models.TABLE_AgentPath,
		models.COLUMN_AgentPath_ParentUid), parentUid).QueryRows(&agentPath)
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}

	return
}

func (d *AgentPathDao) GetTeamPeopleNum(uid uint64) (total uint32, err error) {
	agentPath := models.AgentPath{
		Uid: uid,
	}
	err = d.Orm.Read(&agentPath, models.COLUMN_Agent_Uid)
	if err != nil {
		common.LogFuncError("error:%v", err)
		return
	}

	err = d.Orm.Raw(fmt.Sprintf("SELECT COUNT(%s) FROM %s WHERE %s LIKE ?", models.COLUMN_AgentPath_Uid,
		models.TABLE_AgentPath, models.COLUMN_AgentPath_Path), agentPath.Path+"%").QueryRow(&total)
	if err != nil {
		common.LogFuncError("error:%v", err)
		return
	}

	return
}

//redis存储
func (d *AgentPathDao) CreateRedisAgentPath(pUid, sonUid uint64, level uint32) (err error) {
	agentPath := pb.AgentPath{
		PUId:   pUid,
		Level:  level,
		SubUid: []uint64{},
	}

	var data []byte
	data, err = proto.Marshal(&agentPath)
	if err != nil {
		common.LogFuncError("err:%v", err)
		return
	}

	//sonUid对应redis保存
	key := fmt.Sprintf("%v_%v", AgentPath, sonUid)
	err = common.RedisManger.Set(key, data, 0).Err()
	if err != nil {
		common.LogFuncError("err:%v", err)
		return
	}
	//common.LogFuncDebug("key:%v,path:%+v", key, agentPath)

	//pUid更新
	if pUid > 0 {
		key = fmt.Sprintf("%v_%v", AgentPath, pUid)
		data, err = common.RedisManger.Get(key).Bytes()
		if err != nil {
			common.LogFuncError("err:%v", err)
			return
		}
		var pAgentPath pb.AgentPath
		err = proto.Unmarshal(data, &pAgentPath)
		if err != nil {
			fmt.Printf("err:%v\n", err)
			return
		}
		pAgentPath.SubUid = append(pAgentPath.SubUid, sonUid)
		err = common.RedisManger.Set(key, data, 0).Err()
		if err != nil {
			common.LogFuncError("err:%v", err)
			return
		}
	}

	return
}

//redis存储
func (d *AgentPathDao) GetRedisAgentPathByUid(uid uint64) (path pb.AgentPath, err error) {
	key := fmt.Sprintf("%v_%v", AgentPath, uid)
	b, err := common.RedisManger.Get(key).Bytes()
	if err != nil {
		common.LogFuncError("err:%v", err)
		return
	}

	err = proto.Unmarshal(b, &path)
	if err != nil {
		fmt.Printf("err:%v\n", err)
		return
	}
	//common.LogFuncDebug("key:%v, path:%+v", key, path)

	return
}
