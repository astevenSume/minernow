package common

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

const (
	BulkCount = 100
)

var ErrNoRowEffected = errors.New("no row effected")

var (
	DbOrms = make(map[string]orm.Ormer)
)

// get ormer by specific name.
func GetDbOrm(name string) (dbOrm orm.Ormer, err error) {
	if v, ok := DbOrms[name]; ok {
		dbOrm = v
		return
	}

	err = fmt.Errorf("ormer of %s no exists, please check", name)

	return
}

type Database struct {
	Name         string `json:"name"`
	User         string `json:"user"`
	Password     string `json:"password"`
	Urls         string `json:"urls"`
	MaxIdleConns int    `json:"maxidleconns"`
	MaxOpenConns int    `json:"maxopenconns"`
	MaxLiftTime  int    `json:"maxlifttime"`
}

func DbInit() (err error) {
	var dbMap map[string]string
	dbMap, err = beego.AppConfig.GetSection("database")
	if err != nil {
		LogFuncError("get database configuration failed : %v", err)
		return
	}

	// register default first
	if v, ok := dbMap["default"]; ok {
		config := Database{}
		err = json.Unmarshal([]byte(v), &config)
		if err != nil {
			LogFuncError("json.Unmarshal %s failed : %v", v, err)
			return
		}
		config.Name = "default"

		if _, err = RegisterOrm(config); err != nil {
			return
		}
	}

	for k, v := range dbMap {
		if k == "default" {
			continue
		}

		config := Database{}
		err = json.Unmarshal([]byte(v), &config)
		if err != nil {
			LogFuncError("json.Unmarshal %s failed : %v", v, err)
			return
		}
		config.Name = k

		if _, err = RegisterOrm(config); err != nil {
			return
		}

	}

	return
}

func RegisterOrm(config Database) (o orm.Ormer, err error) {
	err = orm.RegisterDriver("mysql", orm.DRMySQL)
	if err != nil {
		LogFuncError("%v", err)
		return
	}
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4", config.User, config.Password, config.Urls, config.Name)
	err = orm.RegisterDataBase(config.Name, "mysql", dsn, config.MaxIdleConns, config.MaxOpenConns)
	if err != nil {
		LogFuncError("RegisterDataBase mysql for %s failed : %v", dsn, err)
		return
	}

	db, err := orm.GetDB(config.Name)
	if err != nil {
		LogFuncError("orm.GetDB(\"mysql\") failed : %v", err)
		return
	}

	if db == nil {
		LogFuncError("orm.GetDB(\"mysql\") is nil")
		return
	}

	db.SetConnMaxLifetime(time.Second * time.Duration(config.MaxLiftTime))

	o = orm.NewOrm()
	_ = o.Using(config.Name)

	DbOrms[config.Name] = o

	return
}

type EntityInitFunc func() error

type BaseDao struct {
	Orm orm.Ormer
	Db  string
}

func NewBaseDao(db string) (d BaseDao) {
	var err error
	d.Orm, err = GetDbOrm(db)
	if err != nil {
		panic(err.Error())
	}
	d.Db = db
	return
}

func (d BaseDao) BusinessTable(table string) string {
	return d.Db + "." + table
}
