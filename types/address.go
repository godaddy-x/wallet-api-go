package types

import (
	"github.com/godaddy-x/freego/node/common"
	"github.com/godaddy-x/freego/ormx/sqlc"
)

//easyjson:json
type ImportAddressItem struct {
	AddrIndex int64  `json:"addrIndex"`
	IsMemo    int64  `json:"isMemo"`
	Memo      string `json:"memo"`
	WatchOnly int64  `json:"watchOnly"`
	PublicKey string `json:"publicKey"`
	HdPath    string `json:"hdPath"`
}

//easyjson:json
type ImportAddressReq struct {
	common.BaseReq
	AccountID string              `json:"accountID"`
	Addresses []ImportAddressItem `json:"addresses"`
}

//easyjson:json
type ImportAddressRes struct {
	Addresses []AddressResult `json:"addresses"`
}

//easyjson:json
type FindAddressByAddressReq struct {
	common.BaseReq
	Address string `json:"address"`
	Symbol  string `json:"symbol"`
}

//easyjson:json
type FindAddressByAddressRes struct {
	Address AddressResult `json:"address"`
}

//easyjson:json
type FindAddressByAccountIDReq struct {
	common.BaseReq
	AccountID string `json:"accountID"`
}

//easyjson:json
type FindAddressByAccountIDRes struct {
	Result []AddressResult `json:"result"`
	Limit  sqlc.Limit      `json:"limit"`
}

//easyjson:json
type VerifyAddressReq struct {
	common.BaseReq
	Symbol  string `json:"symbol"`
	Address string `json:"address"`
}

//easyjson:json
type VerifyAddressRes struct {
	Result bool `json:"result"`
}

//easyjson:json
type GetBalanceByAddressReq struct {
	common.BaseReq
	Address         string `json:"address"`
	Symbol          string `json:"symbol"`
	ContractAddress string `json:"contractAddress"`
}

//easyjson:json
type GetBalanceByAddressRes struct {
	Balance BalanceResult `json:"balance"`
}

//easyjson:json
type GetAddressBalanceListReq struct {
	common.BaseReq
	AccountID       string `json:"accountID"`
	Symbol          string `json:"symbol"`
	ContractAddress string `json:"contractAddress"`
}

//easyjson:json
type GetAddressBalanceListRes struct {
	Result []BalanceResult `json:"Result"`
	Limit  sqlc.Limit      `json:"limit"`
}
