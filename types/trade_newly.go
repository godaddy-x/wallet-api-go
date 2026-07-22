package types

import (
	"github.com/godaddy-x/freego/node/common"
	"github.com/godaddy-x/freego/ormx/sqlc"
)

// TradeNewlyResult is a pending (unconfirmed scan) trade row (ow_trade_newly).
// Appears when the scanner first indexes the transfer; after confirmation the
// same logical row is promoted to TradeLogResult with the same UniqueHash.
// Field layout matches TradeLogResult plus NeedManualReview. Do not correlate
// pending and confirmed rows by id — use UniqueHash only.
//
//easyjson:json
type TradeNewlyResult struct {
	ID               int64    `json:"id"`
	AppID            string   `json:"appID"`
	WalletID         string   `json:"walletID"`
	AccountID        string   `json:"accountID"`
	Sid              string   `json:"sid"`
	TxID             string   `json:"txID"`
	TxAction         string   `json:"txAction"`
	FlowType         string   `json:"flowType"`
	FromAddress      []string `json:"fromAddress"`
	FromAddressV     []string `json:"fromAddressV"`
	ToAddress        []string `json:"toAddress"`
	ToAddressV       []string `json:"toAddressV"`
	Amount           string   `json:"amount"`
	Fees             string   `json:"fees"`
	MainSymbol       string   `json:"mainSymbol"`
	Symbol           string   `json:"symbol"`
	IsContract       bool     `json:"isContract"`
	BlockHash        string   `json:"blockHash"`
	BlockHeight      int64    `json:"blockHeight"`
	IsMemo           bool     `json:"isMemo"`
	Memo             string   `json:"memo"`
	ApplyTime        int64    `json:"applyTime"`
	DataType         int64    `json:"dataType"`
	BlockTime        int64    `json:"blockTime"`
	Decimals         int64    `json:"decimals"`
	ContractToken    string   `json:"contractToken"`
	ContractAddress  string   `json:"contractAddress"`
	Success          string   `json:"success"`
	OutputIndex      int64    `json:"outputIndex"`
	Signature        string   `json:"signature"`
	CreateAt         int64    `json:"createAt"`
	UpdateAt         int64    `json:"updateAt"`
	UniqueHash       string   `json:"uniqueHash"`       // stable key shared with TradeLogResult for pending→confirmed upsert
	NeedManualReview int64    `json:"needManualReview"` // non-zero when manual approval is required before confirmation
	State            int64    `json:"state"`
}

// FindTradeNewlyReq lists pending trade rows (incremental lastID pagination).
//
//easyjson:json
type FindTradeNewlyReq struct {
	common.BaseReq
}

// FindTradeNewlyRes holds pending scan rows. After confirmation, the same
// transfer appears in FindTradeLogRes; correlate via UniqueHash.
//
//easyjson:json
type FindTradeNewlyRes struct {
	Result []TradeNewlyResult `json:"result"`
	Limit  sqlc.Limit         `json:"limit"`
}
