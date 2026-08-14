package types

import (
	"github.com/godaddy-x/freego/core/query"
	common "github.com/godaddy-x/freego/protocol/dto"
)

//easyjson:json
type CreateWalletReq struct {
	common.BaseReq
	Alias     string `json:"alias"`
	WalletID  string `json:"walletID"`
	RootPath  string `json:"rootPath"`
	Algorithm string `json:"algorithm,omitempty"` // ecdsa | ed25519
}

//easyjson:json
type CreateWalletRes struct {
	Result       bool   `json:"result"`
	WalletID     string `json:"walletID"`
	RootPath     string `json:"rootPath"`
	Alias        string `json:"alias"`
	Algorithm    string `json:"algorithm,omitempty"`
	AccountIndex int64  `json:"accountIndex"`
}

//easyjson:json
type FindWalletByWalletIDReq struct {
	common.BaseReq
	WalletID string `json:"walletID"`
}

//easyjson:json
type FindWalletByWalletIDRes struct {
	Result WalletResult `json:"result"`
}

// FindWalletListReq lists wallets on OPS for the authenticated app.
//
//easyjson:json
type FindWalletListReq struct {
	common.BaseReq
	WalletID string `json:"walletID,omitempty"`
	Alias    string `json:"alias,omitempty"`
}

// FindWalletListRes is the OPS FindWalletList response.
//
//easyjson:json
type FindWalletListRes struct {
	Result []WalletResult `json:"result"`
	Limit  sqlc.Limit     `json:"limit"`
}

//easyjson:json
type WalletResult struct {
	ID           int64  `json:"id"`
	AppID        string `json:"appID"`
	WalletID     string `json:"walletID"`
	RootPath     string `json:"rootPath"`
	Alias        string `json:"alias"`
	Algorithm    string `json:"algorithm,omitempty"` // ecdsa | ed25519
	AccountIndex int64  `json:"accountIndex"`
	CreateAt     int64  `json:"createAt"`
}
