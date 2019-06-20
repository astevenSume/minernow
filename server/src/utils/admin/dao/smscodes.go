package dao

import (
	"common"
	"github.com/astaxie/beego/orm"
	"utils/admin/models"
)

//用户状态
const (
	SmsCodeStatusNil = iota
	SmsCodeStatusActive
	SmsCodeStatusUsed
	SmsCodeStatusExpire
	SmsCodeStatusMax
)

type SmsCodeDao struct {
	common.BaseDao
}

func NewSmsCodeDao(db string) *SmsCodeDao {
	return &SmsCodeDao{
		BaseDao: common.NewBaseDao(db),
	}
}

var SmsCodeDaoEntity *SmsCodeDao

//查询验证码
func (d *SmsCodeDao) QuerySmsCode(smscodes *models.Smscodes, cols ...string) error {
	if smscodes == nil {
		return ErrParam
	}

	err := d.Orm.Read(smscodes, cols...)
	if err != nil {
		common.LogFuncError("mysql_err:%v", err)
		return err
	}
	return nil
}

//分页查询
func (d *SmsCodeDao) QueryPageSmsCodes(page int, perPage int, status int8, mobile, action string) (total int64, list []models.Smscodes, err error) {
	qs := d.Orm.QueryTable(models.TABLE_Smscodes)
	if qs == nil {
		common.LogFuncError("mysql_err:TABLE_Smscodes fail")
		return
	}
	if status > SmsCodeStatusNil && status < SmsCodeStatusMax {
		qs = qs.Filter(models.COLUMN_Smscodes_Status, status)
	}
	if len(mobile) > 0 {
		qs = qs.Filter(models.COLUMN_Smscodes_Mobile, mobile)
	}
	if len(action) > 0 {
		qs = qs.Filter(models.COLUMN_Smscodes_Action, action)
	}

	total, err = qs.Count()
	if err != nil {
		common.LogFuncError("mysql error:%v", err)
		return
	}

	//获取当前页数据的起始位置
	start := (page - 1) * perPage
	if int64(start) > total {
		return
	}
	qs = qs.OrderBy("-" + models.COLUMN_Smscodes_Id)
	_, err = qs.Limit(perPage, start).All(&list)
	if err != nil {
		if err == orm.ErrNoRows {
			return
		}
		common.LogFuncError("mysql error:%v", err)
		return
	}

	return
}

//插入短信验证码
func (d *SmsCodeDao) InsertSmsCode(nationalCode, mobile, action, code string, expire int64) (id int64, err error) {
	smsCodes := new(models.Smscodes)
	smsCodes.Code = code
	smsCodes.Action = action
	smsCodes.NationalCode = nationalCode
	smsCodes.Mobile = mobile
	smsCodes.Status = SmsCodeStatusActive
	smsCodes.Ctime = common.NowInt64MS()
	smsCodes.Etime = smsCodes.Ctime + expire
	id, err = d.Orm.Insert(smsCodes)
	if err != nil {
		common.LogFuncError("mysql error:%v", err)
		return
	}
	return
}

//更新短信验证码
func (d *SmsCodeDao) SetSmsCodeUsed(id int64) (err error) {
	smscodes := &models.Smscodes{
		Id: id,
	}
	err = d.Orm.Read(smscodes, models.COLUMN_Smscodes_Id)
	if err != nil {
		common.LogFuncError("mysql_err:%v", err)
		return
	}
	smscodes.Status = SmsCodeStatusUsed
	_, err = d.Orm.Update(smscodes, models.COLUMN_Smscodes_Status)
	if err != nil {
		common.LogFuncDebug("mysql error:%v", err)
	}

	return
}

//获取短信验证码
func (d *SmsCodeDao) GetSmsCode(nationalCode, mobile, action string, cols ...string) (id uint64, code string, err error) {
	return
}
