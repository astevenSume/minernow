package models

//auto_models_start
type MonthDividendRecord struct {
	Id             uint64 `orm:"column(id);pk" json:"id,omitempty"`
	Uid            uint64 `orm:"column(uid)" json:"uid,omitempty"`
	Ctime          int64  `orm:"column(ctime)" json:"ctime,omitempty"`
	SelfDividend   int64  `orm:"column(self_dividend)" json:"self_dividend,omitempty"`
	AgentDividend  int64  `orm:"column(agent_dividend)" json:"agent_dividend,omitempty"`
	ResultDividend int64  `orm:"column(result_dividend)" json:"result_dividend,omitempty"`
	ReceiveStatus  int32  `orm:"column(receive_status)" json:"receive_status,omitempty"`
	ReceivedTime   int64  `orm:"column(received_time)" json:"received_time,omitempty"`
	PayStatus      int32  `orm:"column(pay_status)" json:"pay_status,omitempty"`
	Level          uint32 `orm:"column(level)" json:"level,omitempty"`
}

func (this *MonthDividendRecord) TableName() string {
	return "month_dividend_record"
}

//table month_dividend_record name and attributes defination.
const TABLE_MonthDividendRecord = "month_dividend_record"
const COLUMN_MonthDividendRecord_Id = "id"
const COLUMN_MonthDividendRecord_Uid = "uid"
const COLUMN_MonthDividendRecord_Ctime = "ctime"
const COLUMN_MonthDividendRecord_SelfDividend = "self_dividend"
const COLUMN_MonthDividendRecord_AgentDividend = "agent_dividend"
const COLUMN_MonthDividendRecord_ResultDividend = "result_dividend"
const COLUMN_MonthDividendRecord_ReceiveStatus = "receive_status"
const COLUMN_MonthDividendRecord_ReceivedTime = "received_time"
const COLUMN_MonthDividendRecord_PayStatus = "pay_status"
const COLUMN_MonthDividendRecord_Level = "level"
const ATTRIBUTE_MonthDividendRecord_Id = "Id"
const ATTRIBUTE_MonthDividendRecord_Uid = "Uid"
const ATTRIBUTE_MonthDividendRecord_Ctime = "Ctime"
const ATTRIBUTE_MonthDividendRecord_SelfDividend = "SelfDividend"
const ATTRIBUTE_MonthDividendRecord_AgentDividend = "AgentDividend"
const ATTRIBUTE_MonthDividendRecord_ResultDividend = "ResultDividend"
const ATTRIBUTE_MonthDividendRecord_ReceiveStatus = "ReceiveStatus"
const ATTRIBUTE_MonthDividendRecord_ReceivedTime = "ReceivedTime"
const ATTRIBUTE_MonthDividendRecord_PayStatus = "PayStatus"
const ATTRIBUTE_MonthDividendRecord_Level = "Level"

//auto_models_end
