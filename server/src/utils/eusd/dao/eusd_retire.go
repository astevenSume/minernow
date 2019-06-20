package dao


import (
	"common"
	"utils/eusd/models"
)

const (
	EusdStatusRetired  = iota  //已销毁

)

type EusdRetireDao struct {
	common.BaseDao
}

func NewEusdRetireDao(db string) *EusdRetireDao {
	return &EusdRetireDao{
		BaseDao: common.NewBaseDao(db),
	}
}

func (d *EusdRetireDao) Add(from_uid uint64,quantity int64,from string) (id int64, err error) {
	data := &models.EusdRetire{
		From:    from,
		FromUid:     from_uid,
		Quantity: quantity,
		Status:  EusdStatusRetired,
		Ctime:   common.NowInt64MS(),
	}

	id, err = d.Orm.Insert(data)
	if err != nil {
		common.LogFuncError("EusdRetireDao RDBERR:%v", err)
		return
	}

	return
}