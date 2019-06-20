package utils

import (
	"common"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base32"
	"encoding/base64"
	"fmt"
	"github.com/skip2/go-qrcode"
	"math/rand"
	"strings"
	"time"
)

const base64Table = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"

//生成32位随机序列
var (
	codes       = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	codeLen     = len(codes)
	secretIdLen = 32
	coder       = base64.NewEncoding(base64Table)
)

func RandSecretId(len int) (secretId string) {
	data := make([]byte, len)
	rand.Seed(time.Now().UnixNano())

	for i := 0; i < len; i++ {
		idx := rand.Intn(codeLen)
		data[i] = byte(codes[idx])
	}
	secretId = string(data)

	return
}

func Base64Encode(encodeByte []byte) []byte {
	return []byte(coder.EncodeToString(encodeByte))
}

//创建二维码
func CreateGoogleAuthQrCode(email string) (secretId string, base64FileData string, err error) {
	var encodeData []byte
	secretId = RandSecretId(secretIdLen)
	url := fmt.Sprintf("otpauth://totp/%s?secret=%s&issuer=qdd", email, secretId)
	encodeData, err = qrcode.Encode(url, qrcode.Medium, 256)
	if err != nil {
		common.LogFuncError("error:%v", err)
		return
	}
	base64FileData = "data:image/png;base64," + string(Base64Encode(encodeData))

	return
}

func toBytes(value int64) []byte {
	var result []byte
	mask := int64(0xFF)
	shifts := [8]uint16{56, 48, 40, 32, 24, 16, 8, 0}
	for _, shift := range shifts {
		result = append(result, byte((value>>shift)&mask))
	}
	return result
}

func toUint32(bytes []byte) uint32 {
	return (uint32(bytes[0]) << 24) + (uint32(bytes[1]) << 16) +
		(uint32(bytes[2]) << 8) + uint32(bytes[3])
}

func oneTimePassword(key []byte, value []byte) uint32 {
	// sign the value using HMAC-SHA1
	hmacSha1 := hmac.New(sha1.New, key)
	hmacSha1.Write(value)
	hash := hmacSha1.Sum(nil)

	// We're going to use a subset of the generated hash.
	// Using the last nibble (half-byte) to choose the index to start from.
	// This number is always appropriate as it's maximum decimal 15, the hash will
	// have the maximum index 19 (20 bytes of SHA1) and we need 4 bytes.
	offset := hash[len(hash)-1] & 0x0F

	// get a 32-bit (4-byte) chunk from the hash starting at offset
	hashParts := hash[offset : offset+4]

	// ignore the most significant bit as per RFC 4226
	hashParts[0] = hashParts[0] & 0x7F

	number := toUint32(hashParts)

	// size to 6 digits
	// one million is the first number with 7 digits so the remainder
	// of the division will always return < 7 digits
	pwd := number % 1000000

	return pwd
}

func googleCodeRedisKey(email string) string {
	return "google_code" + email
}

func VerityGoogleCode(email, secretId, code string) (result bool, err error) {
	// decode the key from the first argument
	inputNoSpaces := strings.Replace(secretId, " ", "", -1)
	inputNoSpacesUpper := strings.ToUpper(inputNoSpaces)

	var key []byte
	key, err = base32.StdEncoding.DecodeString(inputNoSpacesUpper)
	if err != nil {
		common.LogFuncError("error:%v", err)
		return
	}

	// generate a one-time password using the time at 30-second intervals
	epochSeconds := time.Now().Unix()
	verifyCode := fmt.Sprintf("%06d", oneTimePassword(key, toBytes(epochSeconds/30)))
	if verifyCode != code {
		common.LogFuncError("verifyCode:%v != code:%v", verifyCode, code)
		return
	}

	//key := googleCodeRedisKey(email)
	//secondsRemaining := 30 - (epochSeconds % 30)
	//fmt.Printf("%06d-%02d\n", code, secondsRemaining)
	result = true

	return
}
