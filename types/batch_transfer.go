package types

import (
	"github.com/godaddy-x/freego/node/common"
)

// BatchTransferRecipient is one payout line in a batch transfer.
//
//easyjson:json
type BatchTransferRecipient struct {
	Address string `json:"address"`
	Amount  string `json:"amount"` // human-readable decimal; server converts to chain units for ABI
}

// CreateBatchTransferTradeReq builds a pending batch-transfer contract call.
// symbol and contractAddress must refer to a registered batch-transfer deployment.
// Empty tokenAddress uses native batchSendNativeToken (dataType=4); non-empty uses batchSendERC20.
//
//easyjson:json
type CreateBatchTransferTradeReq struct {
	common.BaseReq
	AccountID       string                   `json:"accountID"`
	Sid             string                   `json:"sid"`
	Symbol          string                   `json:"symbol"`
	TxFrom          string                   `json:"txFrom"`
	ContractAddress string                   `json:"contractAddress"`
	TokenAddress    string                   `json:"tokenAddress"`
	Recipients      []BatchTransferRecipient `json:"recipients"`
	FeeRate         string                   `json:"feeRate"`
	ExtParam        string                   `json:"extParam"`
	Memo            string                   `json:"memo"`
	Value           string                   `json:"value"` // optional native value; defaults to sum of recipients
	AwaitResult     bool                     `json:"awaitResult"`
	AwaitTimeout    uint64                   `json:"awaitTimeout"`
}

// CreateBatchTransferApproveTradeReq builds an ERC20 approve for batch transfer (dataType=6).
//
//easyjson:json
type CreateBatchTransferApproveTradeReq struct {
	common.BaseReq
	AccountID       string `json:"accountID"`
	Sid             string `json:"sid"`
	Symbol          string `json:"symbol"`
	TxFrom          string `json:"txFrom"`
	ContractAddress string `json:"contractAddress"`
	TokenAddress    string `json:"tokenAddress"`
	Amount          string `json:"amount"`
	Unlimited       bool   `json:"unlimited"`
	FeeRate         string `json:"feeRate"`
	ExtParam        string `json:"extParam"`
	Memo            string `json:"memo"`
	Value           string `json:"value"`
	AwaitResult     bool   `json:"awaitResult"`
	AwaitTimeout    uint64 `json:"awaitTimeout"`
}

// GetBatchTransferAllowanceReq reads on-chain allowance(owner, spender) for batch transfer.
//
//easyjson:json
type GetBatchTransferAllowanceReq struct {
	common.BaseReq
	AccountID       string `json:"accountID"`
	Symbol          string `json:"symbol"`
	TxFrom          string `json:"txFrom"`
	ContractAddress string `json:"contractAddress"`
	TokenAddress    string `json:"tokenAddress"`
}

// GetBatchTransferAllowanceRes is the on-chain allowance query result.
//
//easyjson:json
type GetBatchTransferAllowanceRes struct {
	Owner        string `json:"owner"`
	Spender      string `json:"spender"`
	TokenAddress string `json:"tokenAddress"`
	Amount       string `json:"amount"`    // human-readable
	AmountWei    string `json:"amountWei"` // smallest-unit decimal string
	Unlimited    bool   `json:"unlimited"` // true when allowance is uint256 max
}

// SpeedUpBatchTransferTradeReq shares fields with SpeedUpTransferTradeReq; separate API name.
type SpeedUpBatchTransferTradeReq = SpeedUpTransferTradeReq

// CancelBatchTransferTradeReq shares fields with CancelTransferTradeReq; separate API name.
type CancelBatchTransferTradeReq = CancelTransferTradeReq
