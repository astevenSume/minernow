package common

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/skip2/go-qrcode"
	"os/exec"
	"strings"
	"time"
)

var (
	ErrFile = errors.New("image is not QrCode")
)

//生成二维码
func GenQrCode(content string, uid uint64, prefix string) (imageFileName string, err error) {
	imageFileName = fmt.Sprintf("%v_%v_%v.png", prefix, uid, time.Now().Unix())
	err = qrcode.WriteFile(content, qrcode.Medium, 256, imageFileName)
	if err != nil {
		LogFuncError("error:%v", err)
		return
	}

	return
}

//二维码识别
func ZBarImgDecode(imageFile string) (decodeData string, err error) {
	var out bytes.Buffer
	cmd := exec.Command("zbarimg", imageFile)
	cmd.Stdout = &out
	err = cmd.Run()
	if err != nil {
		LogFuncError("err:%v", err)
		return
	}
	decodeData = strings.TrimSuffix(strings.TrimPrefix(out.String(), "QR-Code:"), "\n")
	if len(decodeData) == 0 {
		err = ErrFile
		return
	}

	return
}
