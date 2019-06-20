package usdt

import "encoding/json"

func ToJson(arg interface{}) string {
	data, _ := json.Marshal(arg)

	return string(data)
}
