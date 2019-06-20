package dao

import (
	"common"
	"fmt"
	"github.com/astaxie/beego/orm"
	"strings"
	agentmodels "utils/agent/models"
	"utils/otc/models"
)

//用户状态
const (
	UserStatusNil = iota
	UserStatusActive
	UserStatusSuspended
	UserStatusDeleted
	UserStatusPending
	UserStatusMax
)
const (
	INVITE_CODE_LENGTH = 5
	MinSearchMobileNum = 2
)

const (
	NotExchanger     = iota //非承兑商
	Exchanger               //承兑商
	ExchangerPending        //申请成为承兑商中
	ExchangerFail           //申请成为承兑商失败
)

type UserDao struct {
	common.BaseDao
}

func NewUserDao(db string) *UserDao {
	return &UserDao{
		BaseDao: common.NewBaseDao(db),
	}
}

var UserDaoEntity *UserDao

//初始化  sys 账号1，1   platform游戏平台1，2  Commission 1,3
func (d *UserDao) InfoByMobile(nationalCode string, mobile string) (user *models.User, err error) {
	user = &models.User{
		NationalCode: nationalCode,
		Mobile:       mobile,
	}

	err = d.Orm.Read(user, models.COLUMN_User_NationalCode, models.COLUMN_User_Mobile)
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

func (d *UserDao) Create(nationalCode string, mobile string, nick string, ip string) (user *models.User, err error) {
	uid, err := common.IdManagerGen(IdTypeUser)
	if err != nil {
		return
	}

	user = &models.User{
		Uid:          uid,
		NationalCode: nationalCode,
		Mobile:       mobile,
		Nick:         nick,
		Ctime:        common.NowInt64MS(),
		Utime:        common.NowInt64MS(),
		Ip:           ip,
		Status:       UserStatusActive,
		SignSalt:     common.AppSignMgr.GenerateSalt(),
	}

	if _, err = d.Orm.Insert(user); err != nil {
		common.LogFuncError("mysql_err:%v", err)
		return
	}

	return
}

//reset user's signature salt
func (d *UserDao) ResetSignSaltByUId(uid uint64) (user *models.User, err error) {

	user = &models.User{
		Uid: uid,
	}

	err = d.Orm.Read(user, models.COLUMN_User_Uid)
	if err != nil {
		common.LogFuncError("mysql_err:%v", err)
		return
	}

	// reset sign salt
	user.SignSalt = common.AppSignMgr.GenerateSalt()
	_, err = d.Orm.Update(user, models.COLUMN_User_SignSalt)
	if err != nil {
		common.LogFuncError("mysql_err:%v", err)
		return
	}

	return
}

//更新用户昵称
func (d *UserDao) EditNickByUId(uid uint64, name string) (user *models.User, err error) {
	if name == "" {
		err = orm.ErrArgs
		return
	}
	user = &models.User{
		Uid:   uid,
		Nick:  name,
		Utime: common.NowInt64MS(),
	}

	user.Utime = common.NowInt64MS()
	_, err = d.Orm.Update(user, models.COLUMN_User_Nick, models.COLUMN_User_Utime)
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

//获取用户信息
func (d *UserDao) InfoByUId(uid uint64) (user *models.User, err error) {
	user = &models.User{
		Uid: uid,
	}

	err = d.Orm.Read(user, models.COLUMN_User_Uid)
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

//获取用户信息
type UserIsExchanger struct {
	Uid         uint64 `json:"uid"`
	IsExchanger int8   `json:"is_exchanger"`
}

func (d *UserDao) InfoByUIdForIsExchanger(uid uint64) (user *UserIsExchanger, err error) {
	querySql := fmt.Sprintf("select %s,%s from %s where %s=?", models.COLUMN_User_Uid, models.COLUMN_User_IsExchanger, models.TABLE_User, models.COLUMN_User_Uid)
	err = d.Orm.Raw(querySql, uid).QueryRow(&user)
	if err != nil {
		common.LogFuncError("mysql_err:%v", err)
		return
	}

	return
}
func (d *UserDao) BecomeExchange(uid uint64) (ok bool, err error) {
	user := &models.User{
		Uid:         uid,
		IsExchanger: Exchanger,
	}
	_, err = d.Orm.Update(user, models.COLUMN_User_IsExchanger)
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}
	ok = true
	return
}

func (d *UserDao) UnBecomeExchange(uid uint64) (ok bool, err error) {
	user := &models.User{
		Uid:         uid,
		IsExchanger: NotExchanger,
	}
	_, err = d.Orm.Update(user, models.COLUMN_User_IsExchanger)
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}
	ok = true
	return
}

//更新用户
func (d *UserDao) UpdateUser(uid uint64, exchanger int8, status int8, nick, mobile string) (models.User, error) {
	user := models.User{Uid: uid}
	err := d.Orm.Read(&user, models.COLUMN_User_Uid)
	if err != nil {
		common.LogFuncError("mysql err:%v", err)
		return user, err
	}

	var cols []string
	if nick != "" {
		user.Nick = nick
		cols = append(cols, models.COLUMN_User_Nick)
	}
	if exchanger > 0 {
		user.IsExchanger = exchanger
		cols = append(cols, models.COLUMN_User_IsExchanger)
	}
	//if mobile != "" {
	//	user.Mobile = mobile
	//	cols = append(cols, models.COLUMN_User_Mobile)
	//}
	if status > UserStatusNil && status < UserStatusMax {
		user.Status = status
		cols = append(cols, models.COLUMN_User_Status)
	}
	if len(cols) == 0 {
		return user, nil
	}
	user.Utime = common.NowInt64MS()
	cols = append(cols, models.COLUMN_User_Utime)

	_, err = d.Orm.Update(&user, cols...)
	if err != nil {
		common.LogFuncError("mysql_err:%v", err)
		return user, err
	}

	return user, nil
}

//更新用户状态
func (d *UserDao) UpdateStatus(id uint64, status int8) error {
	user := &models.User{
		Uid:    id,
		Status: status,
		Utime:  common.NowInt64MS(),
	}

	_, err := d.Orm.Update(user, models.COLUMN_User_Utime, models.COLUMN_User_Status)
	if err != nil {
		common.LogFuncError("mysql error:%v", err)
		return err
	}

	return nil
}

//更新用户手机号
func (d *UserDao) UpdateMobile(uid uint64, mobile string) error {
	user := &models.User{
		Uid:    uid,
		Mobile: mobile,
		Utime:  common.NowInt64MS(),
	}

	_, err := d.Orm.Update(user, models.COLUMN_User_Utime, models.COLUMN_User_Mobile)
	if err != nil {
		common.LogFuncError("mysql error:%v", err)
		return err
	}

	return nil
}

//更新用户登录时间
func (d *UserDao) UpdateLoginTime(id uint64, loginTime int64) error {
	user := &models.User{
		Uid:           id,
		LastLoginTime: loginTime,
	}

	_, err := d.Orm.Update(user, models.COLUMN_User_LastLoginTime)
	if err != nil {
		common.LogFuncError("mysql error:%v", err)
		return err
	}

	return nil
}

//批量删除用户
func (d *UserDao) DelUsers(ids []uint64) error {
	_, err := d.Orm.QueryTable(models.TABLE_User).Filter("uid__in", ids).Update(orm.Params{"Status": UserStatusDeleted, "Utime": common.NowInt64MS()})
	if err != nil {
		common.LogFuncError("mysql error:%v", err)
		return err
	}

	return nil
}

type OtcUserInfo struct {
	Uid           uint64 `json:"uid"`
	Nick          string `json:"nick"`
	Mobile        string `json:"mobile"`
	Status        int    `json:"status"`
	IsExchanger   int    `json:"is_exchanger"`
	ParentUid     uint64 `json:"parent_uid"`
	Ctime         int64  `json:"ctime"`
	Utime         int64  `json:"utime"`
	LastLoginTime int64  `json:"last_login_time"`
}

type OtcUserInfoClient struct {
	Uid           string `json:"uid"`
	Nick          string `json:"name"`
	Mobile        string `json:"mobile"`
	Status        int    `json:"status"`
	IsExchanger   int    `json:"exchanger"`
	ParentUid     string `json:"parent_uid"`
	Ctime         int64  `json:"ctime"`
	Utime         int64  `json:"utime"`
	LastLoginTime int64  `json:"ltime"`
}

func (d *UserDao) GetOtcUserInfo(user *OtcUserInfo) (info OtcUserInfoClient) {
	if user == nil {
		return
	}

	info = OtcUserInfoClient{
		Uid:           fmt.Sprintf("%v", user.Uid),
		Nick:          user.Nick,
		Status:        int(user.Status),
		Mobile:        user.Mobile,
		ParentUid:     fmt.Sprintf("%v", user.ParentUid),
		Ctime:         user.Ctime,
		Utime:         user.Utime,
		LastLoginTime: user.LastLoginTime,
		IsExchanger:   int(user.IsExchanger),
	}

	return
}

//获取用户信息
func (d *UserDao) UserInfoByUId(uid uint64) (*OtcUserInfo, error) {
	sql := fmt.Sprintf("SELECT T1.%s,T1.%s,T1.%s,T1.%s,T1.%s,T1.%s,T1.%s,T1.%s,T2.%s FROM ((SELECT * FROM %s WHERE %s = ?) AS T1 LEFT JOIN %s AS T2 ON T1.%s=T2.%s)",
		models.COLUMN_User_Uid,
		models.COLUMN_User_Nick,
		models.COLUMN_User_Mobile,
		models.COLUMN_User_Status,
		models.COLUMN_User_IsExchanger,
		models.COLUMN_User_Ctime,
		models.COLUMN_User_Utime,
		models.COLUMN_User_LastLoginTime,
		agentmodels.COLUMN_AgentPath_ParentUid,
		models.TABLE_User,
		models.COLUMN_User_Uid,
		agentmodels.TABLE_AgentPath,
		models.COLUMN_User_Uid,
		agentmodels.COLUMN_AgentPath_Uid)

	otcUserInfo := &OtcUserInfo{}
	err := d.Orm.Raw(sql, uid).QueryRow(otcUserInfo)
	if err != nil {
		if err == orm.ErrNoRows {
			return nil, nil
		}
		common.LogFuncError("mysql_err:%v", err)
		return nil, err
	}

	return otcUserInfo, nil
}

//分页查询
func (d *UserDao) QueryPageUser(name, mobile string, exchanger int8, status int8, page int, limit int) (total int64, otcUserInfo []OtcUserInfo, err error) {
	sqlTotal := fmt.Sprintf("SELECT COUNT(*) AS total FROM %s WHERE 1=1 ", models.TABLE_User)
	sql := fmt.Sprintf("SELECT T1.%s,T1.%s,T1.%s,T1.%s,T1.%s,T1.%s,T1.%s,T1.%s,T2.%s FROM ((SELECT * FROM %s WHERE 1=1 ",
		models.COLUMN_User_Uid,
		models.COLUMN_User_Nick,
		models.COLUMN_User_Mobile,
		models.COLUMN_User_Status,
		models.COLUMN_User_IsExchanger,
		models.COLUMN_User_Ctime,
		models.COLUMN_User_Utime,
		models.COLUMN_User_LastLoginTime,
		agentmodels.COLUMN_AgentPath_ParentUid,
		models.TABLE_User)

	var param []interface{}
	if status > UserStatusNil && status < UserStatusMax {
		sql = fmt.Sprintf("%s AND %s=?", sql, models.COLUMN_User_Status)
		sqlTotal = fmt.Sprintf("%s AND %s=?", sqlTotal, models.COLUMN_User_Status)
		param = append(param, status)
	}

	if exchanger >= NotExchanger {
		sql = fmt.Sprintf(" %s AND %s=?", sql, models.COLUMN_User_IsExchanger)
		sqlTotal = fmt.Sprintf(" %s AND %s=?", sqlTotal, models.COLUMN_User_IsExchanger)
		param = append(param, exchanger)
	}
	if len(name) > 0 {
		sql = fmt.Sprintf("  %s AND %s REGEXP ? ", sql, models.COLUMN_User_Nick)
		sqlTotal = fmt.Sprintf(" %s AND %s REGEXP ? ", sqlTotal, models.COLUMN_User_Nick)
		param = append(param, name)
	}
	if len(mobile) > 0 {
		sql = fmt.Sprintf(" %s AND %s REGEXP ? ", sql, models.COLUMN_User_Mobile)
		sqlTotal = fmt.Sprintf(" %s AND %s REGEXP ? ", sqlTotal, models.COLUMN_User_Mobile)
		param = append(param, mobile)
	}

	err = d.Orm.Raw(sqlTotal, param...).QueryRow(&total)
	if err != nil {
		common.LogFuncError("mysql error: %v", err)
		return
	}

	param = append(param, limit)
	param = append(param, (page-1)*limit)
	sql = fmt.Sprintf("%s LIMIT ? OFFSET ?) AS T1 LEFT JOIN %s AS T2 ON T1.%s=T2.%s)", sql,
		agentmodels.TABLE_AgentPath,
		models.COLUMN_User_Uid,
		agentmodels.COLUMN_AgentPath_Uid)
	_, err = d.Orm.Raw(sql, param...).QueryRows(&otcUserInfo)
	if err != nil {
		common.LogFuncError("mysql error: %v", err)
		return
	}

	return
}

//获取用户信息
func (d *UserDao) GetUserContact(uid uint64, PayId uint64) (string, string) {
	user := &models.User{Uid: uid}
	err := d.Orm.Read(user, models.COLUMN_User_Uid)
	if err != nil {
		if err == orm.ErrNoRows {
			return "", ""
		}
		common.LogFuncError("mysql_err:%v", err)
		return "", ""
	}

	paymentMethod := &models.PaymentMethod{Pmid: PayId}
	err = d.Orm.Read(paymentMethod, models.COLUMN_PaymentMethod_Pmid)
	if err != nil {
		if err == orm.ErrNoRows {
			return user.Mobile, ""
		}
		common.LogFuncError("mysql_err:%v", err)
		return user.Mobile, ""
	}

	return user.Mobile, paymentMethod.QRCode
}

func (d *UserDao) GetUidByMobile(mobile string) (uid uint64, err error) {
	sql := fmt.Sprintf("select %s from %s where %s = ? ", models.COLUMN_User_Uid, models.TABLE_User, models.COLUMN_User_Mobile)
	err = d.Orm.Raw(sql, mobile).QueryRow(&uid)
	if err != nil {
		common.LogFuncError("mysql error:%v", err)
		return
	}

	return
}

func (d *UserDao) GetMobileByUid(uid uint64) (mobile string, err error) {
	sql := fmt.Sprintf("select %s from %s where %s = ? ", models.COLUMN_User_Mobile, models.TABLE_User, models.COLUMN_User_Uid)
	err = d.Orm.Raw(sql, uid).QueryRow(&mobile)
	if err != nil {
		common.LogFuncError("mysql error:%v", err)
		return
	}

	return
}

//统计改名字模糊匹配下的用户数量
func (d *UserDao) CountUserByNick(nick string) (total int64, err error) {
	user := new(models.User)
	qs := d.Orm.QueryTable(user)
	total, err = qs.Filter("nick__icontains", nick).Count()
	if err != nil {
		common.LogFuncError("mysql error:%v", err)
		return
	}
	return
}

//分页查询
func (d *UserDao) FindUserByNick(nick string, page, limit int) (users []models.User, err error) {

	_, err = d.Orm.QueryTable(models.TABLE_User).Filter("nick__icontains", nick).OrderBy("-ctime").Limit(limit, (page-1)*limit).All(&users)
	if err != nil {
		common.LogFuncError("mysql error:%v", err)
		return
	}
	return
}

//分页查询只获取名字和id
type UserNickAndIds struct {
	Uid  uint64 `json:"uid"`
	Nick string `json:"nick"`
}

func (d *UserDao) GetUserByNick(nick string) (user *models.User, err error) {
	querySql := fmt.Sprintf("select * from %s where %s=?", models.TABLE_User, models.COLUMN_User_Nick)
	err = d.Orm.Raw(querySql, nick).QueryRow(&user)
	if err != nil {
		if err == orm.ErrNoRows {
			return
		}
		common.LogFuncError("mysql error:%v", err)
		return
	}
	return
}
func (d *UserDao) FindUserIdByNick(nick string, page, limit int) (users []*models.User, err error) {
	_, err = d.Orm.QueryTable(models.TABLE_User).Filter("nick__icontains", nick).OrderBy("-ctime").Limit(limit, (page-1)*limit).All(&users, models.COLUMN_User_Uid, models.COLUMN_User_Nick)
	if err != nil {
		common.LogFuncError("mysql error:%v", err)
		return
	}
	return
}
func (d *UserDao) GetUserByCondition(nick string) (user *models.User, err error) {
	user = new(models.User)
	sql := fmt.Sprintf("select * from %s where %s like ? ", models.TABLE_User, models.COLUMN_User_Nick)
	err = d.Orm.Raw(sql, nick).QueryRow(&user)
	if err != nil {
		common.LogFuncError("mysql error:%v", err)
		return
	}

	return
}
func (d *UserDao) GetMobileByUIds(uIds []uint64) (user []models.User, err error) {
	if len(uIds) == 0 {
		return
	}
	_, err = d.Orm.QueryTable(models.TABLE_User).Filter(models.COLUMN_User_Uid+"__in", uIds).All(&user,
		models.COLUMN_User_Uid, models.COLUMN_User_Mobile)
	if err != nil {
		common.LogFuncError("error:%v", err)
		return
	}

	return
}

//判断新手机号是否存在
func (d *UserDao) NewMobileIsPresence(mobile string) (num int, err error) {
	sql := fmt.Sprintf("select count(*) from %s where %s = ? ", models.TABLE_User, models.COLUMN_User_Mobile)
	err = d.Orm.Raw(sql, mobile).QueryRow(&num)
	if err != nil {
		common.LogFuncError("mysql error:%v", err)
		return
	}

	return
}

//判断用户是否存在
func (d *UserDao) UserIsPresence(uid string) (num int, err error) {
	sql := fmt.Sprintf("select count(*) from %s where %s = ? ", models.TABLE_User, models.COLUMN_User_Uid)
	err = d.Orm.Raw(sql, uid).QueryRow(&num)
	if err != nil {
		common.LogFuncError("mysql error:%v", err)
		return
	}

	return
}

//统计新增人数
func (d *UserDao) CountByTime(start, end int64) (c uint32) {
	sql := fmt.Sprintf("select count(*) as c from %s where %s>=? and %s<?", models.TABLE_User,
		models.COLUMN_User_Ctime, models.COLUMN_User_Ctime)

	type count struct {
		C uint32 `json:"c"`
	}
	res := &count{}
	if err := d.Orm.Raw(sql, start, end).QueryRow(res); err != nil {
		if err == orm.ErrNoRows {
			err = nil
			return
		}
		common.LogFuncError("DBERR:%v", err)
		return
	}
	c = res.C
	return
}

func (d *UserDao) FindByUids(uids []uint64) (users []*models.User, err error) {
	user := new(models.User)
	qs := d.Orm.QueryTable(user)
	if len(uids) == 0 {
		common.LogFuncWarning("FindByUids no uids %v", err)
	}
	_, err = qs.Filter(models.COLUMN_User_Uid+"__in", uids).All(&users)
	if err != nil {
		common.LogFuncError("data err %v", err)
	}
	return
}

func (d *UserDao) GetByUid(uid uint64) (user *models.User, err error) {
	querySql := fmt.Sprintf("select * from %s where %s=?", models.TABLE_User, models.COLUMN_User_Uid)
	err = d.Orm.Raw(querySql, uid).QueryRow(&user)
	if err != nil {
		common.LogFuncError("mysql error:%v", err)
		return
	}
	return
}
func (d *UserDao) FindMapByUids(uids []uint64) (usersMap map[uint64]*models.User, err error) {
	users := make([]*models.User, 0)
	qs := d.Orm.QueryTable(models.TABLE_User)
	_, err = qs.Filter(models.COLUMN_User_Uid+"__in", uids).All(&users)
	if err != nil {
		common.LogFuncError("data err %v", err)
	}

	usersMap = make(map[uint64]*models.User, 0)
	for _, user := range users {
		usersMap[user.Uid] = user
	}
	return
}

type UserNick struct {
	Uid  uint64 `orm:"column(uid);pk" json:"uid,omitempty"`
	Nick string `orm:"column(nick);size(100)" json:"nick,omitempty"`
}

func (d *UserDao) GetNickByUIds(uIds []string) (user []UserNick, err error) {
	if len(uIds) == 0 {
		return
	}
	sqlQuery := fmt.Sprintf("SELECT %s,%s FROM %s WHERE %s IN(%s) ", models.COLUMN_User_Uid,
		models.COLUMN_User_Nick, models.TABLE_User, models.COLUMN_User_Uid, strings.Join(uIds, ","))
	_, err = d.Orm.Raw(sqlQuery).QueryRows(&user)
	if err != nil {
		common.LogFuncError("mysql error:%v", err)
		return
	}

	return
}
