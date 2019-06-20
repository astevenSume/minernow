package dao

import (
	"common"
	"utils/admin/models"

	"github.com/astaxie/beego/orm"
)

const (
	AppChannelPrecision = 10000
)

type AppChannelDao struct {
	common.BaseDao
}

func NewAppChannelDao(db string) *AppChannelDao {
	return &AppChannelDao{
		BaseDao: common.NewBaseDao(db),
	}
}

var AppChannelDaoEntity *AppChannelDao

//获取app类型信息
func (d *AppChannelDao) QueryById(id uint32) (*models.AppChannel, error) {
	appChannel := &models.AppChannel{
		Id: id,
	}
	err := d.Orm.Read(appChannel, models.COLUMN_AppChannel_Id)
	if err != nil {
		if err == orm.ErrNoRows {
			return nil, nil
		}
		common.LogFuncError("mysql_err:%v", err)
		return nil, err
	}

	return appChannel, nil
}

//所有渠道类型
func (d *AppChannelDao) QueryAll() (appChannels []models.AppChannel, err error) {
	_, err = d.Orm.QueryTable(models.TABLE_AppChannel).All(&appChannels)
	if err != nil {
		if err == orm.ErrNoRows {
			err = nil
			return
		}
		common.LogFuncError("mysql error:%v", err)
		return
	}

	return
}

func (d *AppChannelDao) Create(isThirdHall int8, id uint32, name, desc, iconUrl string, exchangeRate, precision, profitRate int32) (isNew bool, appChannel models.AppChannel, err error) {
	appChannel = models.AppChannel{
		Id:           id,
		IsThirdHall:  isThirdHall,
		Name:         name,
		Desc:         desc,
		ExchangeRate: exchangeRate,
		Precision:    precision,
		ProfitRate:   profitRate,
		IconUrl:      iconUrl,
		Ctime:        common.NowInt64MS(),
		Utime:        common.NowInt64MS(),
	}

	isNew, _, err = d.Orm.ReadOrCreate(&appChannel, models.COLUMN_AppChannel_Id)
	if err != nil {
		common.LogFuncError("mysql err:%v", err)
		return
	}

	return
}

func (d *AppChannelDao) Update(isThirdHall int8, id uint32, name, desc, iconUrl string, exchangeRate, precision, profitRate int32) (appChannel models.AppChannel, err error) {
	appChannel = models.AppChannel{Id: id}
	err = d.Orm.Read(&appChannel, models.COLUMN_AppChannel_Id)
	if err != nil {
		common.LogFuncError("mysql err:%v", err)
		return
	}

	appChannel.Name = name
	appChannel.Desc = desc
	appChannel.IconUrl = iconUrl
	appChannel.IsThirdHall = isThirdHall
	appChannel.ExchangeRate = exchangeRate
	appChannel.Precision = precision
	appChannel.ProfitRate = profitRate
	appChannel.Utime = common.NowInt64MS()
	_, err = d.Orm.Update(&appChannel, models.COLUMN_AppChannel_Name, models.COLUMN_AppChannel_Desc,
		models.COLUMN_AppChannel_Utime, models.COLUMN_AppChannel_IsThirdHall, models.COLUMN_AppChannel_ExchangeRate,
		models.COLUMN_AppChannel_Precision, models.COLUMN_AppChannel_ProfitRate, models.COLUMN_AppChannel_IconUrl)
	if err != nil {
		common.LogFuncError("mysql_err:%v", err)
		return
	}

	return
}

func (d *AppChannelDao) DelById(id uint32) error {
	appChannel := &models.AppChannel{
		Id: id,
	}

	_, err := d.Orm.Delete(appChannel, models.COLUMN_AppChannel_Id)
	if err != nil {
		common.LogFuncError("mysql error:%v", err)
		return err
	}

	return nil
}
