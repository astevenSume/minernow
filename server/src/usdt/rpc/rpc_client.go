package rpc

import (
	"bytes"
	"common"
	"fmt"
	"net/http"
	"utils/otc/dao"
)

// for omnicore Json-api usage.

type Client struct {
}

func NewClient() *Client {
	return &Client{}
}

func (c *Client) Call(result interface{}, method string, params ...interface{}) (err error) {

	url := fmt.Sprintf("http://%s/", UsdtConfig.Host)

	var id uint64
	id, err = common.IdManagerGen(dao.IdTypeUsdtReq)
	var message []byte

	message, err = encodeClientRequest(method, params, id)
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}

	var req *http.Request
	req, err = http.NewRequest("POST", url, bytes.NewBuffer(message))
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(UsdtConfig.User, UsdtConfig.Password)

	client := &http.Client{}
	var resp *http.Response
	resp, err = client.Do(req)
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}
	defer resp.Body.Close()

	err = decodeClientResponse(resp.Body, &result)
	if err != nil {
		common.LogFuncError("%v", err)
		return
	}

	return
}
