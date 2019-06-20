package common

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/des"
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

//Reference： GO密码学应用 https://blog.lab99.org/post/golang-2017-09-23-video-go-for-crypto-developers.html

// @Description aes gcm encrypt
func Encrypt(data []byte, key [32]byte) ([]byte, error) {
	//	初始化 block cipher
	block, err := aes.NewCipher(key[:])
	if err != nil {
		return nil, err
	}
	//	设置 block cipher mode
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	//	生成随机 nonce
	nonce := make([]byte, gcm.NonceSize())
	_, err = rand.Read(nonce)
	if err != nil {
		return nil, err
	}
	//	封装、返回
	return gcm.Seal(nonce, nonce, data, nil), nil
}

// @Description aes gcm decrypt
func Decrypt(ciphertext []byte, key [32]byte) (plaintext []byte, err error) {
	//	初始化 block cipher
	block, err := aes.NewCipher(key[:])
	if err != nil {
		return nil, err
	}
	//	设置 block cipher mode
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	//	返回解开的包，注意这里的 nonce 是直接取的。
	return gcm.Open(nil,
		ciphertext[:gcm.NonceSize()],
		ciphertext[gcm.NonceSize():],
		nil,
	)
}

// @Description aes gcm encrypt + base64 encode
func EncryptToBase64(src string, key [32]byte) (dst string, err error) {
	var buf []byte
	buf, err = Encrypt([]byte(src), key)
	if err != nil {
		return
	}

	// base64 encode
	dst = base64.StdEncoding.EncodeToString(buf)

	return
}

// @Description base64 decode + aes gcm decrypt
func DecryptFromBase64(tmp string, key [32]byte) (src string, err error) {
	var buf []byte
	buf, err = base64.StdEncoding.DecodeString(tmp)
	if err != nil {
		return
	}

	srcByte, err := Decrypt(buf, key)
	if err != nil {
		fmt.Println(err)
		return
	}

	src = string(srcByte)

	return
}

// DesEncrypt Des加密
func DesEncrypt(origData, key []byte) ([]byte, error) {
	block, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}
	origData = PKCS5Padding(origData, block.BlockSize())
	blockMode := cipher.NewCBCEncrypter(block, key)
	crypted := make([]byte, len(origData))
	// 根据CryptBlocks方法的说明，如下方式初始化crypted也可以
	// crypted := origData
	blockMode.CryptBlocks(crypted, origData)
	return crypted, nil
}

// DesDecrypt Des解密
func DesDecrypt(crypted, key []byte) ([]byte, error) {
	block, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockMode := cipher.NewCBCDecrypter(block, key)
	origData := make([]byte, len(crypted))
	blockMode.CryptBlocks(origData, crypted)
	origData = PKCS5UnPadding(origData)
	return origData, nil
}
