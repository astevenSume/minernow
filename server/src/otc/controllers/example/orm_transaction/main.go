package main

import (
	. "common"
	"errors"
	"fmt"
	"github.com/astaxie/beego/orm"
	"sync"
	"time"
	otcmodels "utils/otc/models"
)

func main() {
	var (
		err     error
		configs = []Database{
			{
				Name:         "default",
				User:         "otc",
				Password:     "otc",
				Urls:         "127.0.0.1:3306",
				MaxIdleConns: 10,
				MaxOpenConns: 20,
				MaxLiftTime:  10,
			},
			{
				Name:         "otc",
				User:         "otc",
				Password:     "otc",
				Urls:         "127.0.0.1:3306",
				MaxIdleConns: 10,
				MaxOpenConns: 20,
				MaxLiftTime:  10,
			},
			{
				Name:         "otc_admin",
				User:         "otc",
				Password:     "otc",
				Urls:         "127.0.0.1:3306",
				MaxIdleConns: 10,
				MaxOpenConns: 20,
				MaxLiftTime:  10,
			},
		}
	)

	if err = otcmodels.ModelsInit(); err != nil {
		LogFuncError("%v", err)
		return
	}

	for _, c := range configs {
		_, err = RegisterOrm(c)
		if err != nil {
			LogFuncError("RegisterOrm failed : %v", err)
		}
	}

	orm.Debug = true

	//
	go func() {
		DoTransaction("otc", Test)
	}()

	//time.Sleep(time.Second)
	//
	//go func() {
	//	DoTransaction("otc", Test1)
	//}()

	wg := sync.WaitGroup{}
	wg.Add(2)
	wg.Wait()
}

var ErrNoCorrectRowsEffected = errors.New("no correct rows effected")

func Test(o orm.Ormer) (err error) {
	user := otcmodels.User{}

	err = o.Raw(fmt.Sprintf("SELECT * FROM %s WHERE %s=? FOR UPDATE", otcmodels.TABLE_User, otcmodels.COLUMN_User_Uid), 167628358614515712).QueryRow(&user)
	if err != nil {
		if err != orm.ErrNoRows {
			LogFuncError("#0 %v", err)
			return
		}

		return
	}

	LogFuncDebug("#0 %+v", user)

	time.Sleep(time.Second * 10)

	user.Utime = 1
	var n int64
	n, err = o.Update(&user, otcmodels.COLUMN_User_Utime)
	if err != nil {
		LogFuncError("#0 %v", err)
		return
	}

	if n != 1 {
		err = ErrNoCorrectRowsEffected
		LogFuncError("#0 %d rows effected", n)
		return
	}

	LogFuncDebug("#0 update done")

	err = ErrNoCorrectRowsEffected

	return
}

func Test1(o orm.Ormer) (err error) {
	user := otcmodels.User{}

	err = o.Raw(fmt.Sprintf("SELECT * FROM %s WHERE %s=? FOR UPDATE", otcmodels.TABLE_User, otcmodels.COLUMN_User_Uid), 167628358614515712).QueryRow(&user)
	if err != nil {
		if err != orm.ErrNoRows {
			LogFuncError("#1 %v", err)
			return
		}

		return
	}

	LogFuncDebug("#1 %+v", user)

	user.Utime = 2
	var n int64
	n, err = o.Update(&user, otcmodels.COLUMN_User_Utime)
	if err != nil {
		LogFuncError("#1 %v", err)
		return
	}

	if n != 1 {
		err = ErrNoCorrectRowsEffected
		LogFuncError("#1 %d rows effected", n)
		return
	}

	LogFuncDebug("#1 update done")

	return
}

type TxFunc func(o orm.Ormer) error

func DoTransaction(dbName string, txFunc TxFunc) (err error) {
	o := orm.NewOrm()
	err = o.Using(dbName)
	if err != nil {
		goto errHandle
	}

	err = o.Begin()
	if err != nil {
		goto errHandle
	}

	defer o.Rollback()

	err = txFunc(o)
	if err != nil {
		goto errHandle
	}

	err = o.Commit()
	if err != nil {
		goto errHandle
	}

	return

errHandle:
	LogFuncError("%v", err)
	return
}
