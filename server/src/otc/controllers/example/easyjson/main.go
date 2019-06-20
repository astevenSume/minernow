package main

import (
	//"encoding/json"
	"fmt"
	json "github.com/mailru/easyjson"
)

//easyjson:json
type TestItem struct {
	Id   int    `json:"id"`
	Desc string `json:"desc"`
}

//easyjson:json
type Test struct {
	Id    int        `json:"id"`
	Desc  string     `json:"desc"`
	Items []TestItem `json:"items"`
	SubId []uint8    `json:"sub_ids"`
}

func main() {

	// Marshal
	//t := Test{
	//	Id:   100,
	//	Desc: "desc 100",
	//}
	//
	//t.Items = append(t.Items, TestItem{
	//	Id:   101,
	//	Desc: "desc 101",
	//}, TestItem{
	//	Id:   102,
	//	Desc: "desc 102",
	//})
	//
	//t.SubIds = []int{1, 3, 4}

	//buf, err := json.Marshal(&t)
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//
	//fmt.Print(string(buf))

	//
	t := Test{}
	err := json.Unmarshal([]byte("{\"id\":1, \"desc\":\"heihei\", \"sub_ids\":[1,2,4]}"), &t)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("%+v\n", t)
}
