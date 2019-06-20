package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"
)

var (
	num = flag.Int("n", 1000000, "number of invite code to generate.")
	//region = flag.Int("r", 20, "number of region to divide.")
	file   = flag.String("f", "invite_code", "file to store sql.")
	length = flag.Int("l", 6, "length of invite code.")
)

var BASE = []byte("TXW24N5MP6QL7JB8GHZYVKR9AUC3DSEYF")

func main() {

	rand.Seed(time.Now().UnixNano())

	//var files []*os.File
	for i := 0; i < 108; i++ {
		f, err := os.OpenFile(*file+fmt.Sprintf("_%d.sql", i), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
		if err != nil {
			log.Println(err)
			return
		}

		files = append(files, f)

		defer f.Close()
	}

	var (
		depth int
		args  []int
	)

	rand.Seed(time.Now().UnixNano())
	for i := 0; i < *length; i++ {
		randNum = append(randNum, rand.Intn(len(BASE)))
	}

	for i := 0; i < len(BASE); i++ {
		generate2(depth, args, (i+randNum[depth])%len(BASE))
	}

}

var files []*os.File
var randNum []int
var counter int

func generate2(depth int, args []int, idx int) {
	depth++
	args = append(args, idx)
	if depth == *length {
		generate1(args...)
		//fmt.Println(s)
	} else {
		for i := 0; i < len(BASE); i++ {
			generate2(depth, args, (i+randNum[depth])%len(BASE))
		}
	}
	return
}

func generate1(args ...int) {
	counter++
	var s string
	for _, arg := range args {
		s += fmt.Sprintf("%c", BASE[arg])
	}

	files[counter%len(files)].WriteString(fmt.Sprintf("INSERT invite_code VALUES ('%s', 0);\n", s))

	return
}
