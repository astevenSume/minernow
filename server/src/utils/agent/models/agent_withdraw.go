package models

//auto_models_start
type AgentWithdraw struct {
	Id     uint64 `orm:"column(id);pk" json:"id,omitempty"`
	Uid    uint64 `orm:"column(uid)" json:"uid,omitempty"`
	Amount int64  `orm:"column(amount)" json:"amount,omitempty"`
	Status uint8  `orm:"column(status)" json:"status,omitempty"`
	Ctime  int64  `orm:"column(ctime)" json:"ctime,omitempty"`
	Mtime  int64  `orm:"column(mtime)" json:"mtime,omitempty"`
	Desc   string `orm:"column(desc);null;size(256)" json:"desc,omitempty"`
}

func (this *AgentWithdraw) TableName() string {
	return "agent_withdraw"
}

//table agent_withdraw name and attributes defination.
const TABLE_AgentWithdraw = "agent_withdraw"
const COLUMN_AgentWithdraw_Id = "id"
const COLUMN_AgentWithdraw_Uid = "uid"
const COLUMN_AgentWithdraw_Amount = "amount"
const COLUMN_AgentWithdraw_Status = "status"
const COLUMN_AgentWithdraw_Ctime = "ctime"
const COLUMN_AgentWithdraw_Mtime = "mtime"
const COLUMN_AgentWithdraw_Desc = "desc"
const ATTRIBUTE_AgentWithdraw_Id = "Id"
const ATTRIBUTE_AgentWithdraw_Uid = "Uid"
const ATTRIBUTE_AgentWithdraw_Amount = "Amount"
const ATTRIBUTE_AgentWithdraw_Status = "Status"
const ATTRIBUTE_AgentWithdraw_Ctime = "Ctime"
const ATTRIBUTE_AgentWithdraw_Mtime = "Mtime"
const ATTRIBUTE_AgentWithdraw_Desc = "Desc"

//auto_models_end
