package controllers

import (
	"common"
	"encoding/json"
	"usdt"

	"github.com/astaxie/beego/httplib"
)

const (
	BIT_COIN_FEE_API = "https://bitcoinfees.earn.com/api/v1/fees/list"
)

// SyncRecommendedFees sync usdt recommended fees
func SyncRecommendedFees() {
	syncRecommendedFees()
	return
}

func syncRecommendedFees() {
	type (
		FeesRes struct {
			Fees []struct {
				MinFee     int `json:"minFee"`
				MaxFee     int `json:"maxFee"`
				DayCount   int `json:"dayCount"`
				MemCount   int `json:"memCount"`
				MinDelay   int `json:"minDelay"`
				MaxDelay   int `json:"maxDelay"`
				MinMinutes int `json:"minMinutes"`
				MaxMinutes int `json:"maxMinutes"`
			} `json:"fees"`
		}
	)
	var (
		tmp            = &FeesRes{}
		err            error
		data           []byte
		fee            usdt.FeeMsg
		maxIdx, minIdx int
	)

	if err = httplib.Get(BIT_COIN_FEE_API).ToJSON(tmp); err != nil {
		return
	}

	for i, f := range tmp.Fees {

		if f.MaxMinutes <= 240 {
			if tmp.Fees[maxIdx].MaxMinutes < tmp.Fees[i].MaxMinutes {
				maxIdx = i
			} else if tmp.Fees[maxIdx].MaxMinutes == tmp.Fees[i].MaxMinutes && minFee(tmp.Fees[i].MinFee, tmp.Fees[i].MaxFee) < minFee(tmp.Fees[maxIdx].MinFee, tmp.Fees[maxIdx].MaxFee) {
				maxIdx = i
			}
		}

		if tmp.Fees[minIdx].MaxMinutes > tmp.Fees[i].MaxMinutes {
			minIdx = i
		} else if tmp.Fees[minIdx].MaxMinutes == tmp.Fees[i].MaxMinutes && minFee(tmp.Fees[i].MinFee, tmp.Fees[i].MaxFee) < minFee(tmp.Fees[minIdx].MinFee, tmp.Fees[minIdx].MaxFee) {
			minIdx = i
		}

		if f.MaxMinutes >= 0 && f.MaxMinutes <= 30 {
			if minFee(f.MinFee, f.MaxFee) > 0 && (fee.HalfHourFee == 0 || fee.HalfHourFee > minFee(f.MinFee, f.MaxFee)) {
				fee.HalfHourFee = minFee(f.MinFee, f.MaxFee)
			}
		}

		if f.MaxMinutes > 30 && f.MaxMinutes <= 60 {
			if minFee(f.MinFee, f.MaxFee) > 0 && (fee.HourFee == 0 || fee.HourFee > minFee(f.MinFee, f.MaxFee)) {
				fee.HourFee = minFee(f.MinFee, f.MaxFee)
			}
		}

		if f.MaxMinutes > 60 && f.MaxMinutes <= 120 {
			if minFee(f.MinFee, f.MaxFee) > 0 && (fee.TwoHourFee == 0 || fee.TwoHourFee > minFee(f.MinFee, f.MaxFee)) {
				fee.TwoHourFee = minFee(f.MinFee, f.MaxFee)
			}
		}

		if f.MaxMinutes > 120 && f.MaxMinutes <= 240 {
			if minFee(f.MinFee, f.MaxFee) > 0 && (fee.FourHourFee == 0 || fee.FourHourFee > minFee(f.MinFee, f.MaxFee)) {
				fee.FourHourFee = minFee(f.MinFee, f.MaxFee)
			}
		}

	}

	if fee.FourHourFee == 0 {
		fee.FourHourFee = tmp.Fees[maxIdx].MinFee
	}

	if fee.HalfHourFee == 0 {
		fee.HalfHourFee = tmp.Fees[minIdx].MinFee
	}

	if fee.TwoHourFee == 0 {
		fee.TwoHourFee = fee.FourHourFee
	}

	if fee.HourFee == 0 {
		fee.HourFee = fee.TwoHourFee
	}

	if data, err = json.Marshal(&fee); err != nil {
		return
	}

	common.RedisManger.Set(usdt.REDIS_KEY_RECOMMENDED_FEES, string(data), 0)
	return
}

func minFee(fee1, fee2 int) int {
	if fee1 == 0 {
		return fee2
	}
	if fee1 > fee2 {
		return fee2
	}
	return fee1
}
