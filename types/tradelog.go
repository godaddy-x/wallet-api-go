package types

import (
	"github.com/godaddy-x/freego/node/common"
	"github.com/godaddy-x/freego/ormx/sqlc"
)

//easyjson:json
type BalanceLogResult struct {
	ID                 int64  `json:"id"`
	AppID              string `json:"appID"`
	WalletID           string `json:"walletID"`
	AccountID          string `json:"accountID"`
	Scope              string `json:"scope"` // address | account
	Address            string `json:"address"`
	MainSymbol         string `json:"mainSymbol"`
	Symbol             string `json:"symbol"`
	ContractAddress    string `json:"contractAddress"`
	TradeLogID         int64  `json:"tradeLogID"`
	TxID               string `json:"txID"`
	BlockHeight        int64  `json:"blockHeight"`
	NetworkBlockHeight int64  `json:"networkBlockHeight"`
	TxAction           string `json:"txAction"`
	FlowType           string `json:"flowType"`
	BalanceBefore      string `json:"balanceBefore"`
	BalanceChange      string `json:"balanceChange"`
	BalanceAfter       string `json:"balanceAfter"`
	ConfirmBefore      string `json:"confirmBefore"`
	ConfirmChange      string `json:"confirmChange"`
	ConfirmAfter       string `json:"confirmAfter"`
	CreateAt           int64  `json:"createAt"`
}

// TradeLogResult is a confirmed trade log entry.
// TxAction: send | receive | internal | fee. DataType: see TradeDataType (0=none, 1=transfer, 2=summary, 3=contract, 4=batch, 5=deploy, 6=approve, 7=stake, 8=unstake, 9=withdraw unfreeze).
//easyjson:json
type TradeLogResult struct {
	ID              int64    `json:"id"`
	AppID           string   `json:"appID"`
	WalletID        string   `json:"walletID"`
	AccountID       string   `json:"accountID"`
	Sid             string   `json:"sid"`
	TxID            string   `json:"txID"`
	TxAction        string   `json:"txAction"`
	FlowType        string   `json:"flowType"`
	FromAddress     []string `json:"fromAddress"`
	FromAddressV    []string `json:"fromAddressV"`
	ToAddress       []string `json:"toAddress"`
	ToAddressV      []string `json:"toAddressV"`
	Amount          string   `json:"amount"`
	Fees            string   `json:"fees"`
	MainSymbol      string   `json:"mainSymbol"`
	Symbol          string   `json:"symbol"`
	IsContract      bool     `json:"isContract"`
	BlockHash       string   `json:"blockHash"`
	BlockHeight     int64    `json:"blockHeight"`
	IsMemo          bool     `json:"isMemo"`
	Memo            string   `json:"memo"`
	ApplyTime       int64    `json:"applyTime"`
	DataType        int64    `json:"dataType"`
	DataTypeName    string   `json:"dataTypeName"`
	BlockTime       int64    `json:"blockTime"`
	Decimals        int64    `json:"decimals"`
	ContractToken   string   `json:"contractToken"`
	ContractAddress string   `json:"contractAddress"`
	Success         string   `json:"success"`
	OutputIndex     int64    `json:"outputIndex"` // BTC vout; ETH token log index; native -1; fee -2
	Signature       string   `json:"signature"`
	CreateAt        int64    `json:"createAt"`
	UpdateAt        int64    `json:"updateAt"`
	UniqueHash      string   `json:"uniqueHash"`
	State           int64    `json:"state"`
}

//easyjson:json
type FindTradeLogReq struct {
	common.BaseReq
}

//easyjson:json
type FindTradeLogRes struct {
	Result []TradeLogResult `json:"result"`
	Limit  sqlc.Limit       `json:"limit"`
}

//easyjson:json
type FindBalanceLogReq struct {
	common.BaseReq
}

//easyjson:json
type FindBalanceLogRes struct {
	Result []BalanceLogResult `json:"result"`
	Limit  sqlc.Limit         `json:"limit"`
}
