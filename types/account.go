package types

import (
	"github.com/godaddy-x/freego/node/common"
	"github.com/godaddy-x/freego/ormx/sqlc"
)

//easyjson:json
type CreateAccountReq struct {
	common.BaseReq
	WalletID       string   `json:"walletID"`
	Alias          string   `json:"alias"`
	Symbol         string   `json:"symbol"`
	OtherOwnerKeys []string `json:"otherOwnerKeys"`
	ReqSigs        int64    `json:"reqSigs"`
	PublicKey      string   `json:"publicKey"`
	Password       string   `json:"password"`
	AccountIndex   int64    `json:"accountIndex"`
	AccountID      string   `json:"accountID"`
	HdPath         string   `json:"hdPath"`
	Remark         string   `json:"remark"`
}

//easyjson:json
type CreateAccountRes struct {
	Account AccountResult `json:"account"`
}

//easyjson:json
type FindAccountByAccountIDReq struct {
	common.BaseReq
	AccountID string `json:"accountID"`
}

//easyjson:json
type FindAccountByAccountIDRes struct {
	Account AccountResult `json:"account"`
}

//easyjson:json
type GetBalanceByAccountReq struct {
	common.BaseReq
	AccountID       string `json:"accountID"`
	Symbol          string `json:"symbol"`
	ContractAddress string `json:"contractAddress"`
	UserID          int64  `json:"userID"`
}

//easyjson:json
type GetBalanceByAccountRes struct {
	Balance BalanceResult `json:"balance"`
}

//easyjson:json
type GetAccountBalanceListReq struct {
	common.BaseReq
	AccountID       string `json:"accountID"`
	Symbol          string `json:"symbol"`
	ContractAddress string `json:"contractAddress"`
}

//easyjson:json
type GetAccountBalanceListRes struct {
	Result []BalanceResult `json:"Result"`
	Limit  sqlc.Limit      `json:"limit"`
}

//easyjson:json
type FindAccountByWalletIDReq struct {
	common.BaseReq
	WalletID string `json:"walletID"`
}

//easyjson:json
type FindAccountByWalletIDRes struct {
	Result []AccountResult `json:"Result"`
	Limit  sqlc.Limit      `json:"limit"`
}

//easyjson:json
type AccountResult struct {
	ID             int64    `json:"id"`
	AppID          string   `json:"appID"`
	WalletID       string   `json:"walletID"`
	AccountID      string   `json:"accountID"`
	Alias          string   `json:"alias"`
	MainSymbol     string   `json:"mainSymbol"`
	OtherOwnerKeys []string `json:"otherOwnerKeys"`
	ReqSigs        int64    `json:"reqSigs"`
	IsTrust        int64    `json:"isTrust"`
	PublicKey      string   `json:"publicKey"`
	HdPath         string   `json:"hdPath"`
	AccountIndex   int64    `json:"accountIndex"`
	AddressIndex   int64    `json:"addressIndex"`
	CreateAt       int64    `json:"createAt"`
	Remark         string   `json:"remark"`
}

// AddressResult is an on-chain address under an account.
//
//easyjson:json
type AddressResult struct {
	ID         int64  `json:"id"`
	AppID      string `json:"appID"`
	WalletID   string `json:"walletID"`
	AccountID  string `json:"accountID"`
	Alias      string `json:"alias"`
	MainSymbol string `json:"mainSymbol"`
	AddrIndex  int64  `json:"addrIndex"`
	Address    string `json:"address"`
	IsMemo     int64  `json:"isMemo"`
	Memo       string `json:"memo"`
	WatchOnly  int64  `json:"watchOnly"`
	PublicKey  string `json:"publicKey"`
	HdPath     string `json:"hdPath"`
	CreateAt   int64  `json:"createAt"`
	UpdateAt   int64  `json:"updateAt"`
	Status     int64  `json:"status"` // 1=active 2=maintenance 3=frozen 4=deprecated 5=abnormal
	State      int64  `json:"state"`
}

//easyjson:json
type BalanceResult struct {
	ID               int64  `json:"id"`
	AppID            string `json:"appID"`
	WalletID         string `json:"walletID"`
	AccountID        string `json:"accountID"`
	Address          string `json:"address"`
	MainSymbol       string `json:"mainSymbol"`
	Symbol           string `json:"symbol"`
	ContractID       string `json:"contractID"`
	ContractAddr     string `json:"contractAddr"`
	Balance          string `json:"balance"`
	ConfirmBalance   string `json:"confirmBalance"`
	UnconfirmBalance string `json:"unconfirmBalance"`
	UpdateAt         int64  `json:"updateAt"`
	ContractToken    string `json:"contractToken"`
}
