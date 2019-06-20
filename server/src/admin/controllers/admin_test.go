package controllers

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"testing"
)

func TestGetUrlPermission(t *testing.T) {
	fmt.Println(GetUrlPermission("/v/admin/user/1223", "get"))
	fmt.Println(GetUrlPermission("/v/admin/user/123/role/12", "put"))
	fmt.Println(GetUrlPermission("/v/admin/user", "post"))
	fmt.Println(GetUrlPermission("/v/admin/user/role", "delete"))
}

func TestEncodingFile(t *testing.T) {
	image, _ := ioutil.ReadFile("test.jpg")
	imageBase64 := base64.StdEncoding.EncodeToString(image)
	buf := []byte(imageBase64)
	ioutil.WriteFile("./test.jpg.txt", buf, 0666)
}
