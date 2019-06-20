package common

import (
	"fmt"
	"github.com/aliyun/oss"
)

func UpFile(accessKeyId, accessKeySecret, endpoint, bucketName, objectName, localFileName string) (string, error) {
	// 创建OSSClient实例。
	client, err := oss.New(endpoint, accessKeyId, accessKeySecret)
	if err != nil {
		LogFuncError("UpFile error:%v", err)
		return "", err
	}

	// 获取存储空间。
	bucket, err := client.Bucket(bucketName)
	if err != nil {
		LogFuncError("UpFile error:%v", err)
		return "", err
	}

	// 上传文件。
	err = bucket.PutObjectFromFile(objectName, localFileName)
	if err != nil {
		LogFuncError("UpFile error:%v", err)
		return "", err
	}
	url := fmt.Sprintf("http://%s.%s/%s", bucketName, endpoint, objectName)

	return url, nil
}
