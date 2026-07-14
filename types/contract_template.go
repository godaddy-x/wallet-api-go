package types

import (
	"github.com/godaddy-x/freego/node/common"
	"github.com/godaddy-x/freego/ormx/sqlc"
	adapter "github.com/godaddy-x/wallet-adapter"
)

//easyjson:json
type ContractTemplateResult struct {
	ID          int64  `json:"id"`
	Code        string `json:"code"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Category    string `json:"category"`
	MainSymbol  string `json:"mainSymbol"`
	Symbol      string `json:"symbol"`
	Protocol    string `json:"protocol"`
	ABI         string `json:"abi"`
	Bytecode    string `json:"bytecode"`
	SourceCode  string `json:"sourceCode"`
	Version     string `json:"version"`
	OrderNo     int64  `json:"orderNo"`
	CreateAt    int64  `json:"createAt"`
	UpdateAt    int64  `json:"updateAt"`
}

//easyjson:json
type GetContractTemplatesReq struct {
	common.BaseReq
	Category   string `json:"category"`
	Code       string `json:"code"`
	MainSymbol string `json:"mainSymbol"`
	Protocol   string `json:"protocol"`
}

//easyjson:json
type GetContractTemplatesRes struct {
	Result []ContractTemplateResult `json:"result"`
	Limit  sqlc.Limit               `json:"limit"`
}

// DeployContractReq creates a pending contract deploy trade from a template (dataType=5).
//
//easyjson:json
type DeployContractReq struct {
	common.BaseReq
	TemplateID   int64  `json:"templateID"`
	AccountID    string `json:"accountID"`
	Sid          string `json:"sid"`
	TxFrom       string `json:"txFrom"`
	Symbol       string `json:"symbol"`
	FeeRate      string `json:"feeRate"`
	ExtParam     string `json:"extParam"`
	Memo         string `json:"memo"`
	Value        string `json:"value"`
	AwaitResult  bool   `json:"awaitResult"`
	AwaitTimeout uint64 `json:"awaitTimeout"`
}

//easyjson:json
type DeployContractRes struct {
	Sid           string                   `json:"sid"`
	PendingSignTx []*adapter.PendingSignTx `json:"pendingSignTx"`
}

// SubmitDeployContractReq broadcasts a signed deploy tx and persists the deployment record on success.
// templateID is deprecated; the trade record linked by pendingSignTx.data.sid is authoritative.
//
//easyjson:json
type SubmitDeployContractReq struct {
	common.BaseReq
	TemplateID    int64                  `json:"templateID"` // deprecated, ignored
	PendingSignTx *adapter.PendingSignTx `json:"pendingSignTx"`
}

//easyjson:json
type SubmitDeployContractRes struct {
	Receipt         *adapter.SmartContractReceipt `json:"receipt,omitempty"`
	DeployID        int64                         `json:"deployID"`
	TemplateID      int64                         `json:"templateID"`
	ContractAddress string                        `json:"contractAddress"`
	Status          int64                         `json:"status"`
	SubmitAt        int64                         `json:"submitAt"`
	ConfirmAt       int64                         `json:"confirmAt"`
	Reason          string                        `json:"reason"`
}
