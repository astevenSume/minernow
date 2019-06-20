package dao

import (
	"common"
	"errors"
	"fmt"
	"github.com/astaxie/beego/orm"
	"utils/admin/models"
)

//用户状态
const (
	SmsTemplateTypeVerifyCode = iota //验证码模板
	SmsTemplateTypeUsdtWarning
	SmsTemplateTypeUsdtSignWarning
)

type SmsTemplateDao struct {
	common.BaseDao
}

func NewSmsTemplateDao(db string) *SmsTemplateDao {
	return &SmsTemplateDao{
		BaseDao: common.NewBaseDao(db),
	}
}

var SmsTemplateDaoEntity *SmsTemplateDao

//短信模板表
func (d *SmsTemplateDao) QuerySmsTemplates(smstemplates *models.Smstemplates, cols ...string) error {
	if smstemplates == nil {
		return errors.New("param error")
	}

	err := d.Orm.Read(smstemplates, cols...)
	if err != nil {
		return err
	}
	return nil
}

//分页查询
func (d *SmsTemplateDao) QueryPageSmsTemplates(page int, perPage int) ([]models.Smstemplates, PageInfo, error) {
	var data []models.Smstemplates
	var info PageInfo
	qs := d.Orm.QueryTable(models.TABLE_Smstemplates)
	if qs == nil {
		common.LogFuncError("mysql_err:Permission fail")
		return data, info, ErrParam
	}

	count, err := qs.Count()
	if err != nil {
		if err == orm.ErrNoRows {
			return data, info, nil
		}
		common.LogFuncError("mysql error:%v", err)
		return data, info, err
	}
	info.Total = int(count)
	info.Page = page
	info.Limit = perPage
	//获取当前页数据的起始位置
	start := (page - 1) * perPage
	if int64(start) > count {
		return data, info, ErrParam
	}
	_, err = qs.Limit(perPage, start).All(&data)
	if err != nil {
		if err == orm.ErrNoRows {
			return data, info, nil
		}
		common.LogFuncError("mysql error:%v", err)
		return data, info, err
	}

	return data, info, nil
}

//添加
func (d *SmsTemplateDao) AddSmsTemplates(smstemplates *models.Smstemplates) error {
	if smstemplates == nil {
		return ErrParam
	}

	id, err := d.Orm.Insert(smstemplates)
	if err != nil {
		common.LogFuncError("mysql error:%v", err)
		return err
	}
	smstemplates.Id = id

	return nil
}

//更新
func (d *SmsTemplateDao) UpdateSmsTemplates(newCom *models.Smstemplates, cols ...string) error {
	if newCom == nil {
		return ErrParam
	}

	smstemplates := &models.Smstemplates{Id: newCom.Id}
	err := d.QuerySmsTemplates(smstemplates, models.COLUMN_Smstemplates_Id)
	if err != nil {
		common.LogFuncError("mysql error:%v", err)
		return err
	}

	smstemplates = newCom
	_, err = d.Orm.Update(smstemplates, cols...)
	if err != nil {
		common.LogFuncError("mysql_err:%v", err)
		return err
	}

	return nil
}

//删除
func (d *SmsTemplateDao) DelSmsTemplates(id int64) error {
	smstemplates := &models.Smstemplates{Id: id}
	_, err := d.Orm.Delete(smstemplates)
	if err != nil {
		common.LogFuncError("mysql_err:%v", err)
		return err
	}

	return nil
}

//获取短信模板
func (d *SmsTemplateDao) GetSmsTemplates(smsType int8) (template string, err error) {
	err = d.Orm.Raw(fmt.Sprintf("SELECT %s from %s WHERE %s=?", models.COLUMN_Smstemplates_Template,
		models.TABLE_Smstemplates, models.COLUMN_Smstemplates_Type), smsType).QueryRow(&template)
	if err != nil {
		common.LogFuncError("mysql error:%v", err)
	}
	return
}
