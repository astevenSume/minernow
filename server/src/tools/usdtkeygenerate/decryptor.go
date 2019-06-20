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
		fmt.Printf("Usage : %s <encrypt string>.\n", os.Args[0])
		return
	}

	tmp := flag.Arg(0)
	src, err := common.DecryptFromBase64(tmp, usdt.PriAesKey)
	if err != nil {
		fmt.Printf("DecryptFromBase64 failed : %v", err)
		return
	}

	fmt.Println(src)
}
