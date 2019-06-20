package models

//auto_models_start
 type ReportStatisticSum struct{
	Id uint64 `orm:"column(id);pk" json:"id,omitempty"`
	ChannelId uint32 `orm:"column(channel_id)" json:"channel_id,omitempty"`
	ChannelPositiveNums uint64 `orm:"column(channel_positive_nums)" json:"channel_positive_nums,omitempty"`
	ChannelSalaryDaily int64 `orm:"column(channel_salary_daily)" json:"channel_salary_daily,omitempty"`
	ChannelRgDividend int64 `orm:"column(channel_rg_dividend)" json:"channel_rg_dividend,omitempty"`
	ChannelWithdrawEusd int64 `orm:"column(channel_withdraw_eusd)" json:"channel_withdraw_eusd,omitempty"`
	ChannelRechargeEusd int64 `orm:"column(channel_recharge_eusd)" json:"channel_recharge_eusd,omitempty"`
	Ctime int64 `orm:"column(ctime)" json:"ctime,omitempty"`
}

func (this *ReportStatisticSum) TableName() string {
    return "report_statistic_sum"
}

//table report_statistic_sum name and attributes defination.
const TABLE_ReportStatisticSum = "report_statistic_sum"
const COLUMN_ReportStatisticSum_Id = "id"
const COLUMN_ReportStatisticSum_ChannelId = "channel_id"
const COLUMN_ReportStatisticSum_ChannelPositiveNums = "channel_positive_nums"
const COLUMN_ReportStatisticSum_ChannelSalaryDaily = "channel_salary_daily"
const COLUMN_ReportStatisticSum_ChannelRgDividend = "channel_rg_dividend"
const COLUMN_ReportStatisticSum_ChannelWithdrawEusd = "channel_withdraw_eusd"
const COLUMN_ReportStatisticSum_ChannelRechargeEusd = "channel_recharge_eusd"
const COLUMN_ReportStatisticSum_Ctime = "ctime"
const ATTRIBUTE_ReportStatisticSum_Id = "Id"
const ATTRIBUTE_ReportStatisticSum_ChannelId = "ChannelId"
const ATTRIBUTE_ReportStatisticSum_ChannelPositiveNums = "ChannelPositiveNums"
const ATTRIBUTE_ReportStatisticSum_ChannelSalaryDaily = "ChannelSalaryDaily"
const ATTRIBUTE_ReportStatisticSum_ChannelRgDividend = "ChannelRgDividend"
const ATTRIBUTE_ReportStatisticSum_ChannelWithdrawEusd = "ChannelWithdrawEusd"
const ATTRIBUTE_ReportStatisticSum_ChannelRechargeEusd = "ChannelRechargeEusd"
const ATTRIBUTE_ReportStatisticSum_Ctime = "Ctime"

//auto_models_end
