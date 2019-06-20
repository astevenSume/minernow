package dao

import (
	"common"
	"strings"
	"utils/admin/models"
)

const (
	SETTINGS_PPRFIX_KEY         = "config_system_namespace_"
	SETTINGS_REDIS_EXPIRES_SECS = 24 * 60 * 60 // 1 DAY
	CONFIG_KEY_ALL              = "all"
)

var configActionMap = map[string]int8{
	"trade": 1, //交易配置
	"usdt":  2, //usdt配置
	"eos":   3, //eos配置
	"game":  4, //游戏配置
	"sys":   5, //系统配置
}

type ConfigDao struct {
	common.BaseDao
}

func NewConfigDao(db string) *ConfigDao {
	return &ConfigDao{
		BaseDao: common.NewBaseDao(db),
	}
}

func (d *ConfigDao) GenerateConfigKey(key string) string {
	return SETTINGS_PPRFIX_KEY + key
}

var ConfigDaoEntity *ConfigDao

func (d *ConfigDao) GetActionType(action string) (int8, bool) {
	value, ok := configActionMap[action]
	return value, ok
}

//分页查询系统配置信息
func (d *ConfigDao) QueryPageConfig(action int8, page int, perPage int) (int64, []models.Config, error) {
	var config []models.Config
	qs := d.Orm.QueryTable(models.TABLE_Config)
	if qs == nil {
		common.LogFuncError("mysql_err:Permission fail")
		return 0, config, ErrSql
	}
	if action > 0 {
		qs = qs.Filter(models.COLUMN_Config_Action, action)
	}

	count, err := qs.Count()
	if err != nil {
		common.LogFuncError("mysql_err:Count fail: %v", err)
		return count, config, err
	}

	//获取当前页数据的起始位置
	start := (page - 1) * perPage
	if int64(start) > count {
		return count, config, nil
	}
	if _, err := qs.Limit(perPage, start).All(&config); err != nil {
		return count, config, err
	}

	return count, config, nil
}

//insert config
func (d *ConfigDao) InsertConfig(action int8, key, value, desc string) error {
	if err := d.checkConfig(key, value); err != nil {
		return err
	}

	config := &models.Config{
		Key:    key,
		Action: action,
		Desc:   desc,
		Value:  value,
		Ctime:  common.NowInt64MS(),
	}

	id, err := d.Orm.Insert(config)
	if err != nil {
		common.LogFuncError("mysql_err:%v", err)
		return err
	}

	// update redis data
	if err := common.RedisManger.Set(d.GenerateConfigKey(config.Key), config.Value, 0).Err(); err != nil {
		common.LogFuncError("redis update err: %s:%s,%v", config.Key, config.Value, err)
	}
	config.Id = uint32(id)
	return nil
}

//update config
func (d *ConfigDao) UpdateConfig(id uint32, action int8, key, value, desc string) error {
	if err := d.checkConfig(key, value); err != nil {
		return err
	}
	config := &models.Config{
		Id:     id,
		Key:    key,
		Action: action,
		Desc:   desc,
		Value:  value,
	}

	var tmpKey = config.Key
	var tmpAction = config.Action
	if tmpConfig, err := d.QueryConfigById(config.Id); err != nil {
		return err
	} else {
		tmpKey = tmpConfig.Key
		tmpAction = tmpConfig.Action
	}

	if tmpKey != config.Key || tmpAction != config.Action {
		return ErrParam
	}

	_, err := d.Orm.Update(config,
		models.COLUMN_Config_Key,
		models.COLUMN_Config_Value,
		models.COLUMN_Config_Desc,
	)
	if err != nil {
		common.LogFuncError("mysql_err:%v", err)
		return err
	}

	return nil
}

//delete by id
func (d *ConfigDao) DelConfigById(id uint32, action int8) error {
	tmpConfig, err := d.QueryConfigById(id)
	if err != nil {
		return err
	}

	if action != tmpConfig.Action {
		return ErrParam
	}

	config := &models.Config{
		Id: id,
	}
	_, err = d.Orm.Delete(config, models.COLUMN_Config_Id)
	if err != nil {
		common.LogFuncError("mysql error:%s, %v", tmpConfig.Key, err)
		return err
	}
	common.RedisManger.Del(d.GenerateConfigKey(tmpConfig.Key))

	return nil
}

//delete by id
func (d *ConfigDao) DelConfigByKey(key string) error {
	config := &models.Config{
		Key: key,
	}
	if _, err := d.QueryConfigByKey(key); err != nil {
		return err
	}

	_, err := d.Orm.Delete(config, models.COLUMN_Config_Key)
	if err != nil {
		common.LogFuncError("mysql error:%v", err)
		return err
	}

	common.RedisManger.Del(d.GenerateConfigKey(key))

	return nil
}

// search by key field
func (d *ConfigDao) QueryConfigByKey(key string) (models.Config, error) {
	config := models.Config{
		Key: key,
	}
	err := d.Orm.Read(&config, models.COLUMN_Config_Key)
	return config, err
}

// config by id field
func (d *ConfigDao) QueryConfigById(id uint32) (models.Config, error) {
	config := models.Config{
		Id: id,
	}
	err := d.Orm.Read(&config, models.COLUMN_Config_Id)
	return config, err
}

//分页查询系统配置信息
func (d *ConfigDao) DumpConfig(action int8, id uint32) ([]models.Config, error) {
	var config []models.Config
	qs := d.Orm.QueryTable(models.TABLE_Config)
	if action > 0 {
		qs = qs.Filter(models.COLUMN_Config_Action, action)
	}
	if id > 0 {
		qs = qs.Filter(models.COLUMN_Config_Id, id)
	}
	if _, err := qs.All(&config); err != nil {
		return config, err
	}

	return config, nil
}

func (d *ConfigDao) checkConfig(key, value string) error {
	// cannot insert key eq "anonymous",
	if value == "" || key == "" || strings.ToUpper(key) == CONFIG_KEY_ALL {
		return ErrParam
	}

	return nil
}

func (d *ConfigDao) Fetch(args ...string) (res []*models.Config) {
	res = []*models.Config{}
	_, err := d.Orm.QueryTable(models.TABLE_Config).Filter(models.COLUMN_Config_Key+"__in", args).All(&res)
	if err != nil {
		common.LogFuncError("DBERR:%v", err)
		return
	}

	return
}
