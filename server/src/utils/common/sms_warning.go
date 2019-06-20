package common

import (
	"common"
	"fmt"
	admindao "utils/admin/dao"
)

//短信预警 param参数对应模板配置参数
func SmsWarning(wType int8, param map[string]string) (err error) {
	accessKeyID, err := AppConfigMgr.String(SmsAccessAppId)
	if err != nil {
		return
	}

	accessSecret, err := AppConfigMgr.String(SmsAccessSecret)
	if err != nil {
		return
	}

	configWarning, err := admindao.ConfigWarningDaoEntity.GetConfigWarning(wType)
	if err != nil {
		common.LogFuncError("error:%v", err)
		return err
	}

	var msg string
	if len(configWarning) > 0 {
		msg, err = getSmsParam(configWarning[0].SmsType, param)
		if err != nil {
			common.LogFuncDebug("error:%v", err)
		}

		for _, item := range configWarning {
			phoneNum := fmt.Sprintf("%v%v", item.NationalCode, item.Mobile)
			result := common.NewSendSmsOverSea(accessKeyID, accessSecret, phoneNum, msg)
			if result != nil {
				common.LogFuncError("sendalisms fail err:%v", result)
			}
		}
	}

	return
}
