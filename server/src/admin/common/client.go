package common

const (
	ClientTypeUnkownStr = "unkown"
	ClientTypeWebStr    = "web"
	ClientTypeAppStr    = "app"
)

const (
	ClientTypeUnkown = iota
	ClientTypeWeb
	ClientTypeApp
)

var MapClientTypeStr = map[int]string{
	ClientTypeUnkown: ClientTypeUnkownStr,
	ClientTypeWeb:    ClientTypeWebStr,
	ClientTypeApp:    ClientTypeAppStr,
}

var MapClientType = map[string]int{
	ClientTypeUnkownStr: ClientTypeUnkown,
	ClientTypeWebStr:    ClientTypeWeb,
	ClientTypeAppStr:    ClientTypeApp,
}

// @Title GetClientType
// @Description get app type enum by app type string.
func GetClientType(clientTypeStr string) (clientType int) {
	if t, ok := MapClientType[clientTypeStr]; ok {
		clientType = t
		return
	}

	return ClientTypeUnkown
}

// @Title GetClientTypeStr
// @Description get app type enum by app type int.
func GetClientTypeStr(clientType int) (clientTypeStr string) {
	if t, ok := MapClientTypeStr[clientType]; ok {
		clientTypeStr = t
		return
	}

	return ClientTypeUnkownStr
}
