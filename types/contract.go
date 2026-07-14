package types

import (
	"github.com/godaddy-x/freego/node/common"
	"github.com/godaddy-x/freego/ormx/sqlc"
	adapter "github.com/godaddy-x/wallet-adapter"
)

//easyjson:json
type ContractResult struct {
	ID       int64  `json:"id"`
	Symbol   string `json:"symbol"`
	Name     string `json:"name"`
	Icon     string `json:"icon"`
	Decimals int64  `json:"decimals"`
	Address  string `json:"address"`
	Token    string `json:"token"`
	Protocol string `json:"protocol"`
	CreateAt int64  `json:"createAt"`
}

//easyjson:json
type GetContractsReq struct {
	common.BaseReq
}

//easyjson:json
type GetContractsRes struct {
	Result []ContractResult `json:"result"`
	Limit  sqlc.Limit       `json:"limit"`
}

//easyjson:json
type CreateSmartContractTradeReq struct {
	common.BaseReq
	AccountID    string            `json:"accountID"`
	Sid          string            `json:"sid"`
	Coin         CoinInfo          `json:"coin"`
	FeeRate      string            `json:"feeRate"`
	ExtParam     string            `json:"extParam"`
	Memo         string            `json:"memo"`
	To           map[string]string `json:"to"`
	Raw          string            `json:"raw"`
	RawType      uint64            `json:"rawType"`
	ABIParam     []string          `json:"abiParam"`
	Value        string            `json:"value"`
	AwaitResult  bool              `json:"awaitResult"`
	AwaitTimeout uint64            `json:"awaitTimeout"`
	TxFrom       string            `json:"txFrom"`
	TxTo         string            `json:"txTo"`
	DataType     int64             `json:"dataType"` // deprecated: use CreateBatchTransferTrade / DeployContract / CreateBatchTransferApproveTrade
}

//easyjson:json
type SubmitSmartContractTradeReq struct {
	common.BaseReq
	PendingSignTx *adapter.PendingSignTx `json:"pendingSignTx"`
}

//easyjson:json
type SubmitSmartContractTradeRes struct {
	Receipt *adapter.SmartContractReceipt `json:"receipt,omitempty"`
}

// CallSmartContractABIReq is a read-only eth_call; ABI is resolved from contract config or standard ERC20 when dataType=6.
//
//easyjson:json
type CallSmartContractABIReq struct {
	common.BaseReq
	AccountID       string   `json:"accountID"`
	Symbol          string   `json:"symbol"`
	ContractAddress string   `json:"contractAddress"`
	ABIParam        []string `json:"abiParam"`
	Raw             string   `json:"raw"`
	RawType         uint64   `json:"rawType"`
	Value           string   `json:"value"`
	TxFrom          string   `json:"txFrom"`
	DataType        int64    `json:"dataType"`
	Decimals        int64    `json:"decimals"`
}

// CallSmartContractABIRes is the read-only call result.
//
//easyjson:json
type CallSmartContractABIRes struct {
	Method         string `json:"method"`
	Value          string `json:"value"`
	RawHex         string `json:"rawHex"`
	Status         uint64 `json:"status"`
	Exception      string `json:"exception"`
	Uint256Wei     string `json:"uint256Wei,omitempty"`
	Uint256Human   string `json:"uint256Human,omitempty"`
}
