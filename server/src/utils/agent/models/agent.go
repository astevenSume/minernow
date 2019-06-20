package models

//auto_models_start
type Agent struct {
	Uid            uint64 `orm:"column(uid);pk" json:"uid,omitempty"`
	SumSalary      int64  `orm:"column(sum_salary)" json:"sum_salary,omitempty"`
	SumWithdraw    int64  `orm:"column(sum_withdraw)" json:"sum_withdraw,omitempty"`
	SumCanWithdraw int64  `orm:"column(sum_can_withdraw)" json:"sum_can_withdraw,omitempty"`
	Ctime          int64  `orm:"column(ctime)" json:"ctime,omitempty"`
	Mtime          int64  `orm:"column(mtime)" json:"mtime,omitempty"`
	Pwd            string `orm:"column(pwd);size(64)" json:"pwd,omitempty"`
}

func (this *Agent) TableName() string {
	return "agent"
}

//table agent name and attributes defination.
const TABLE_Agent = "agent"
const COLUMN_Agent_Uid = "uid"
const COLUMN_Agent_SumSalary = "sum_salary"
const COLUMN_Agent_SumWithdraw = "sum_withdraw"
const COLUMN_Agent_SumCanWithdraw = "sum_can_withdraw"
const COLUMN_Agent_Ctime = "ctime"
const COLUMN_Agent_Mtime = "mtime"
const COLUMN_Agent_Pwd = "pwd"
const ATTRIBUTE_Agent_Uid = "Uid"
const ATTRIBUTE_Agent_SumSalary = "SumSalary"
const ATTRIBUTE_Agent_SumWithdraw = "SumWithdraw"
const ATTRIBUTE_Agent_SumCanWithdraw = "SumCanWithdraw"
const ATTRIBUTE_Agent_Ctime = "Ctime"
const ATTRIBUTE_Agent_Mtime = "Mtime"
const ATTRIBUTE_Agent_Pwd = "Pwd"

//auto_models_end
