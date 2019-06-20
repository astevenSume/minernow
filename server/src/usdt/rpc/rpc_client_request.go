package rpc

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
)

// for omnicore Json-api usage.

// Json-rpc request sent by client
type clientRequest struct {
	//the method to be invoked.
	Method string `json:"method"`
	//
	Params []interface{} `json:"params"`
	//
	Id uint64 `json:"id"`
}

// Json-rpc response return to client
type clientResponse struct {
	//
	Result *json.RawMessage `json:"result"`
	Error  interface{}      `json:"error"`
	Id     uint64           `json:"id"`
}

// encode a Json-rpc request
func encodeClientRequest(method string, params []interface{}, id uint64) ([]byte, error) {
	return json.Marshal(&clientRequest{
		Method: method,
		Params: params,
		Id:     id,
	})
}

// decode a Json-rpc response
func decodeClientResponse(r io.Reader, reply interface{}) (err error) {
	var c clientResponse
	if err = json.NewDecoder(r).Decode(&c); err != nil {
		return
	}

	if c.Error != nil {
		return fmt.Errorf("%v", c.Error)
	}

	if c.Result == nil {
		return errors.New("result is nil")
	}

	return json.Unmarshal(*c.Result, reply)
}
