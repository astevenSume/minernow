package main

import (
	"common"
	"flag"
	"fmt"
	"os"
	"usdt"
)

func main() {
	flag.Parse()

	if flag.NArg() < 1 {
		fmt.Printf("Usage : %s <source string>.\n", os.Args[0])
		return
	}

	tmp := flag.Arg(0)
	dst, err := common.EncryptToBase64(tmp, usdt.PriAesKey)
	if err != nil {
		fmt.Printf("EncryptToBase64 failed : %v", err)
		return
	}

	fmt.Println(dst)
}
