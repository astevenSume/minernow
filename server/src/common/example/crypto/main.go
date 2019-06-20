package main

import (
	"common"
	"fmt"
)

var key = [32]byte{'T', 'X', 'W', '2', '4', 'N', '5', 'M',
	'P', '6', 'Q', 'L', '7', 'J', 'B', '8',
	'G', 'H', 'Z', 'Y', 'V', 'K', 'R', '9',
	'A', 'U', 'C', '3', 'D', 'S', 'E', '='}

func main() {
	src := "5JPzYaNNyrUKTo7pZ5rMH1grWHkhSYSD6ieyXPuXGtSKbMZSD2S"
	dst, err := common.Encrypt([]byte(src), key)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(string(dst))

	tmp, err := common.Decrypt(dst, key)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("tmp %s\n", string(tmp))
}
