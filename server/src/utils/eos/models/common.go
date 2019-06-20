package models

import "github.com/astaxie/beego/orm"

func ModelsInit() (err error) {
	orm.RegisterModel(
		new(EosAccountKeys),
	)

	return
}
