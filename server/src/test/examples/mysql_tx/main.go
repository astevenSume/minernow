package main

import (
	"common"
	"database/sql"
	"github.com/astaxie/beego/orm"
	"time"
	usdtmodels "utils/usdt/models"
)

func main() {

	err := orm.RegisterDriver("mysql", orm.DRMySQL)
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}

	err = orm.RegisterDataBase("otc", "mysql", "otc:otc@tcp(127.0.0.1:3306)/otc?charset=utf8mb4", 10, 20)
	if err != nil {
		common.LogFuncError("RegisterDataBase failed : %v", err)
		return
	}

	err = orm.RegisterDataBase("default", "mysql", "otc:otc@tcp(127.0.0.1:3306)/default?charset=utf8mb4", 10, 20)
	if err != nil {
		common.LogFuncError("RegisterDataBase failed : %v", err)
		return
	}

	db, err := orm.GetDB("otc")
	if err != nil {
		common.LogFuncError("orm.GetDB(\"mysql\") failed : %v", err)
		return
	}

	if db == nil {
		common.LogFuncError("orm.GetDB(\"mysql\") is nil")
		return
	}

	db.SetConnMaxLifetime(time.Second * time.Duration(1200))

	orm.RegisterModel(&usdtmodels.UsdtAccount{})

	tx(0, db, false)

	go tx(1, db, true)

	forever := make(chan bool)
	<-forever
}

func tx(sn int, db *sql.DB, isCommit bool) (err error) {
	o, err := orm.NewOrmWithDB("mysql", "otc", db)
	if err != nil {
		common.LogFuncError("#%d %v", sn, err)
		return
	}

	o.Begin()
	if isCommit {
		defer o.Commit()
	}

	_, err = o.Raw("update otc.usdt_account set available_integer = available_integer+10 where uid=1 and available_integer=0").Exec()
	if err != nil {
		common.LogFuncError("#%d %v", sn, err)
		return
	}

	common.LogFuncDebug("#%d update done.", sn)

	a := &usdtmodels.UsdtAccount{
		Uaid: 1,
	}
	err = o.Read(a)
	if err != nil {
		common.LogFuncError("#%d %v", sn, err)
		return
	}

	common.LogFuncDebug("#%d %+v.", sn, a)

	return
}
