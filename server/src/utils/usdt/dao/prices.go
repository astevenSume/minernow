package dao

import (
	"common"
	"errors"
	"fmt"
	"math"
	"utils/usdt/models"

	"github.com/astaxie/beego/orm"
)

const (
	PRICE_MARKET_OKEX = iota
	PRICE_MARKET_HBG
	PRICE_MARKET_BIAN
	PRICE_MARKET_MAX
)

type PricesDao struct {
	common.BaseDao
}

func NewPricesDao(db string) *PricesDao {
	return &PricesDao{
		BaseDao: common.NewBaseDao(db),
	}
}

var PricesDaoEntity *PricesDao

type PriceAck struct {
	Id     string `json:"id"`
	Market int8   `json:"market"`
	Price  string `json:"price"`
	Ctime  int64  `json:"ctime"`
}

func (d *PricesDao) ClientPrice(data *models.MarketPrices) (ack PriceAck) {
	if data == nil {
		return
	}
	ack.Id = fmt.Sprintf("%v", data.Id)
	ack.Market = data.Market
	ack.Price = fmt.Sprintf("%.2f", float64(data.PowPrice)/float64(math.Pow10(int(data.Pow))))
	ack.Ctime = data.Ctime

	return
}

func (d *PricesDao) ClientPrices(datas []models.MarketPrices) (acks []PriceAck) {
	for _, data := range datas {
		var ack PriceAck
		ack.Id = fmt.Sprintf("%v", data.Id)
		ack.Market = data.Market
		ack.Price = fmt.Sprintf("%.2f", float64(data.PowPrice)/float64(math.Pow10(int(data.Pow))))
		ack.Ctime = data.Ctime
		acks = append(acks, ack)
	}

	return
}

//query last Pirce
func (d *PricesDao) QueryLastPrice(currency string) (models.Prices, error) {
	price := models.Prices{}

	qs := d.Orm.QueryTable(models.TABLE_Prices)
	if qs == nil {
		return price, errors.New("db error")
	}

	qs = qs.Filter(models.COLUMN_Prices_Currency, currency)

	err := qs.OrderBy("-" + models.COLUMN_Prices_Ctime).
		One(&price)
	if err != nil {
		common.LogFuncError("DB_ERR:%v", err)
		return price, err
	}

	return price, nil
}

//inset marketPirce
func (d *PricesDao) InsertMarketPrice(marketPrice *models.MarketPrices) error {
	marketPrice.Ctime = common.NowInt64MS()

	id, err := d.Orm.Insert(marketPrice)
	if err != nil {
		common.LogFuncError("mysql_err:%v", err)
		return err
	}
	marketPrice.Id = uint64(id)
	return nil
}

//inset marketPirce
func (d *PricesDao) InsertPrice(price *models.Prices) error {
	price.Ctime = common.NowInt64MS()

	id, err := d.Orm.Insert(price)
	if err != nil {
		common.LogFuncError("mysql_err:%v", err)
		return err
	}
	price.Id = uint64(id)
	return nil
}

//获取指导价
func (d *PricesDao) QueryById(id uint64) (marketPrices models.MarketPrices, err error) {
	marketPrices = models.MarketPrices{
		Id: id,
	}

	err = d.Orm.Read(&marketPrices, models.COLUMN_MarketPrices_Id)
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

//分页条件查询
func (d *PricesDao) QueryPageCondition(market int8, page int, perPage int, currency string) (total int, marketPrices []models.MarketPrices, err error) {
	qs := d.Orm.QueryTable(models.TABLE_MarketPrices)
	if qs == nil {
		common.LogFuncError("QueryTable fail")
		return
	}

	qs = qs.Filter(models.COLUMN_MarketPrices_Currency, currency)

	if market >= PRICE_MARKET_OKEX && market < PRICE_MARKET_MAX {
		qs = qs.Filter(models.COLUMN_MarketPrices_Market, market)
	}
	qs = qs.OrderBy("-" + models.COLUMN_MarketPrices_Id)

	var count int64
	count, err = qs.Count()
	if err != nil {
		if err == orm.ErrNoRows {
			err = nil
			return
		}
		common.LogFuncError("mysql error:%v", err)
		return
	}
	total = int(count)

	//获取当前页数据的起始位置
	start := (page - 1) * perPage
	if int64(start) > count {
		err = nil
		return
	}
	_, err = qs.Limit(perPage, start).All(&marketPrices)
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
