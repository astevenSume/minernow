package main

import (
	"common"
	"flag"
	"fmt"
	"os"
	"usdt"
	"utils/usdt/dao"

	"github.com/btcsuite/btcd/wire"
)

var (
	count    = *(flag.Int("n", 10000, "how many keys to spawn."))
	regionId = *(flag.Int64("r", 0, "region id of main server."))
	serverId = *(flag.Int64("s", 0, "server id of main server."))
	sql      string
)

func main() {
	flag.Parse()

	err := common.IdManagerInit(dao.IdTypeMax, regionId, serverId)
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}

	for i := 0; i < count; i++ {
		pri, addr, _, err := usdt.GenerateSegWitKey(wire.MainNet)
		if err != nil {
			return
		}

		id, err := common.IdManagerGen(dao.IdTypePriKey)
		if err != nil {
			common.LogFuncError("%v", err)
			return
		}

		var priDb string
		priDb, err = common.EncryptToBase64(pri, usdt.PriAesKey)
		if err != nil {
			common.LogFuncError("%v", err)
			return
		}

		sql += fmt.Sprintf("INSERT INTO usdt_prikey values(%v, '%s', '%s', 0);\n", id, priDb, addr)
	}

	var file *os.File
	file, err = os.OpenFile("pri.sql", os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0666)
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}

	_, err = file.WriteString(sql)
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}
}
