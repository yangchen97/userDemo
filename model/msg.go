package model

import (
	"encoding/json"
	"log"
)

type RespMsg struct {
	Code string      `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}




func (resp *RespMsg) JsonBytes() []byte{
	r, err := json.Marshal(resp)
	if err != nil {
		log.Printf("json marshal err: %v", err)
	}

	return r

}