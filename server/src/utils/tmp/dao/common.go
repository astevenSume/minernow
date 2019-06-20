package dao

import "common"

func Init(entityInitFunc common.EntityInitFunc) (err error) {
	const dbOtc = "otc"
	TmpGameBetersDaoEntity = NewTmpGameBetersDao(dbOtc)

	if entityInitFunc != nil {
		err = entityInitFunc()
		if err != nil {
			return
		}
	}
	return
}
