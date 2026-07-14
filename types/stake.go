package types

import "github.com/godaddy-x/freego/node/common"

// CreateStakeTradeReq builds TRX Stake 2.0 freezeBalanceV2 pending sign tx (dataType=7).
//
//easyjson:json
type CreateStakeTradeReq struct {
	common.BaseReq
	AccountID string `json:"accountID"`
	Sid       string `json:"sid"`
	Symbol    string `json:"symbol"`
	TxFrom    string `json:"txFrom"`
	Amount    string `json:"amount"`
	Resource  string `json:"resource"` // ENERGY | BANDWIDTH
	FeeRate   string `json:"feeRate"`
	ExtParam  string `json:"extParam"`
	Memo      string `json:"memo"`
}

// CreateUnstakeTradeReq builds unfreezeBalanceV2 pending sign tx (dataType=8).
//
//easyjson:json
type CreateUnstakeTradeReq struct {
	common.BaseReq
	AccountID string `json:"accountID"`
	Sid       string `json:"sid"`
	Symbol    string `json:"symbol"`
	TxFrom    string `json:"txFrom"`
	Amount    string `json:"amount"`
	Resource  string `json:"resource"`
	FeeRate   string `json:"feeRate"`
	ExtParam  string `json:"extParam"`
	Memo      string `json:"memo"`
}

// CreateWithdrawUnfreezeTradeReq builds withdrawExpireUnfreeze pending sign tx (dataType=9).
//
//easyjson:json
type CreateWithdrawUnfreezeTradeReq struct {
	common.BaseReq
	AccountID string `json:"accountID"`
	Sid       string `json:"sid"`
	Symbol    string `json:"symbol"`
	TxFrom    string `json:"txFrom"`
	FeeRate   string `json:"feeRate"`
	ExtParam  string `json:"extParam"`
	Memo      string `json:"memo"`
}

// GetAccountResourceDetailReq reads TRX energy/bandwidth/frozen balance.
//
//easyjson:json
type GetAccountResourceDetailReq struct {
	common.BaseReq
	AccountID string `json:"accountID"`
	Symbol    string `json:"symbol"`
	Address   string `json:"address"`
}

// GetAccountResourceDetailRes TRX Stake 2.0 resource snapshot.
//
//easyjson:json
type GetAccountResourceDetailRes struct {
	Address                 string `json:"address"`
	EnergyLimit             int64  `json:"energyLimit"`
	EnergyUsed              int64  `json:"energyUsed"`
	EnergyRemaining         int64  `json:"energyRemaining"`
	BandwidthRemaining      int64  `json:"bandwidthRemaining"`
	AvailableBalance        string `json:"availableBalance"`
	FrozenForEnergy         string `json:"frozenForEnergy"`
	FrozenForBandwidth      string `json:"frozenForBandwidth"`
	WithdrawableUnfreezeSUN int64  `json:"withdrawableUnfreezeSUN"`
}
