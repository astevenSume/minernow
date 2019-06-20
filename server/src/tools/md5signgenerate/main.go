package main

import (
	"common"
	"fmt"
	otcdao "utils/otc/dao"
)

func main() {
	fmt.Println(common.GenerateDoubleMD5ByParams(map[string]string{
		"name": "",
	}, otcdao.SIGNATURE_SALT, 1557315319))
}
