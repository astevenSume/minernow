package cron

import (
	"math"
	"strconv"
	"time"
	otcDao "utils/otc/dao"
	"utils/otc/models"
)

func EosOtcReport() {
	limit := 200
	date := getDate(-1)
	_, err := otcDao.EosOtcReportDaoEntity.DeleteByDate(date)
	if err != nil {
		return
	}
	//获取总页数,向上取整
	totalNum, err := otcDao.OrdersDaoEntity.GetOrdersNumByDate(date)
	totalPage := math.Ceil(float64(totalNum) / float64(limit))
	totalPageInt := int(totalPage)

	otcReports := make(map[uint64]*models.EosOtcReport)

	for page := 1; page <= totalPageInt; page++ {
		orderInfoList, err := otcDao.OrdersDaoEntity.GetOrdersByDate(date, page, limit)
		if err != nil {
			return
		}
		//遍历订单统计报表数据
		for _, o := range orderInfoList {
			uid := o.Uid
			euid := o.EUid
			if _, ok := otcReports[uid]; !ok {
				otcReports[uid] = new(models.EosOtcReport)
				otcReports[uid].Uid = uid
				otcReports[uid].Date = o.Date
			}
			if _, ok := otcReports[euid]; !ok {
				otcReports[euid] = new(models.EosOtcReport)
				otcReports[euid].Uid = uid
				otcReports[euid].Date = o.Date
			}
			//订单总数量+1
			otcReports[uid].TotalOrderNum += 1
			otcReports[euid].TotalOrderNum += 1

			if o.Status == otcDao.OrderStatusConfirmed {
				// 已确认的订单代表成功，成功订单数量+1
				otcReports[uid].SuccessOrderNum += 1

				//只有已经成功的订单才可以统计金额
				if o.Side == otcDao.SideBuy {
					//此订单是买入的将Amount记录到买入的eusd的金额数
					otcReports[uid].BuyEusdNum += o.Amount
					otcReports[euid].SellEusdNum += o.Amount
				} else if o.Side == otcDao.SideSell {
					//此订单是卖出的将Amount记录到卖出的eusd的金额数
					otcReports[euid].BuyEusdNum += o.Amount
					otcReports[uid].SellEusdNum += o.Amount
				}
			} else {
				//其他状态下的订单代表失败，失败订单数量+1
				otcReports[uid].FailOrderNum += 1
				otcReports[euid].FailOrderNum += 1
			}
		}
	}

	//将报表数据写入数据库
	reports := make([]models.EosOtcReport, 0)
	for _, report := range otcReports {
		reports = append(reports, *report)
	}
	if len(reports) == 0 {
		return
	}
	err = otcDao.EosOtcReportDaoEntity.Create(reports)
	if err != nil {
		return
	}
}

func getDate(day int) int32 {
	nTime := time.Now()
	yesTime := nTime.AddDate(0, 0, day)
	date, _ := strconv.ParseInt(yesTime.Format("20060102"), 10, 64)
	return int32(date)
}
