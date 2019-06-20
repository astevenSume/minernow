package dao

import "common"

var EosAccountKeysEntity *EosAccountKeysDao

func Init(entityInitFunc common.EntityInitFunc) (err error) {

	const db = "eos"
	EosAccountKeysEntity = NewEosAccountKeysDao(db)

	if entityInitFunc != nil {
		err = entityInitFunc()
		if err != nil {
			return
		}
	}

	return
}
