package types

import (
	common "github.com/godaddy-x/freego/protocol/dto"
)

//easyjson:json
type CreateSubscribeReq struct {
	common.BaseReq
	SubscribeMethod   []string `json:"subscribeMethod"`
	SubscribeContract []string `json:"subscribeContract"`
}

//easyjson:json
type CreateSubscribeRes struct {
	Result bool `json:"result"`
}
