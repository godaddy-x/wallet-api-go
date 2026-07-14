package types

import (
	"github.com/godaddy-x/freego/node/common"
	adapter "github.com/godaddy-x/wallet-adapter"
)

//easyjson:json
type AssetsAccount struct {
	WalletID        string   `json:"walletID"`
	AccountID       string   `json:"accountID"`
	Index           uint64   `json:"index"`
	HdPath          string   `json:"hdPath"`
	PublicKey       string   `json:"publicKey"`
	OwnerKeys       []string `json:"ownerKeys"`
	ContractAddress string   `json:"contractAddress"`
	Symbol          string   `json:"symbol"`
	AddressIndex    int64    `json:"addressIndex"`
	MainSymbol      string   `json:"mainSymbol"`
	Alias           string   `json:"alias"`
	ReqSigs         int64    `json:"reqSigs"`
	IsTrust         bool     `json:"isTrust"`
}

//easyjson:json
type CoinInfo struct {
	Symbol          string        `json:"symbol"`
	IsContract      bool          `json:"isContract"`
	ContractAddress string        `json:"contractAddress"`
	ContractABI     string        `json:"contractABI"`
	Contract        SmartCoinInfo `json:"contract"`
}

//easyjson:json
type FeesSupportAccount struct {
	AccountID        string `json:"accountID"`
	FixSupportAmount string `json:"fixSupportAmount"`
	FeesSupportScale string `json:"feesSupportScale"`
}

//easyjson:json
type SmartCoinInfo struct {
	Symbol   string `json:"symbol"`
	Address  string `json:"address"`
	Token    string `json:"token"`
	Protocol string `json:"protocol"`
	Name     string `json:"name"`
	Decimals uint64 `json:"decimals"`
}

//easyjson:json
type KeySig struct {
	Address     string `json:"address"`
	DerivedPath string `json:"derivedPath"`
	EccType     uint32 `json:"eccType"`
	InputIndex  int64  `json:"inputIndex"`
	Msg         string `json:"msg"`
	Nonce       string `json:"nonce"`
	Signed      string `json:"signed"`
	WalletID    string `json:"walletID"`
	IsImport    int64  `json:"isImport"`
	PublicKey   string `json:"publicKey"`
	Rsv         bool   `json:"rsv"`
}

//easyjson:json
type SigParts struct {
	Key string   `json:"key"`
	Sig []KeySig `json:"sig"`
}

//easyjson:json
type SubmitRawTransaction struct {
	AppID     string            `json:"appID"`
	WalletID  string            `json:"walletID"`
	AccountID string            `json:"accountID"`
	Coin      CoinInfo          `json:"coin"`
	RawHex    string            `json:"rawHex"`
	RawHexSig string            `json:"rawHexSig"`
	ReqSigs   uint64            `json:"reqSigs"`
	Sid       string            `json:"sid"`
	SigCount  int64             `json:"sigCount"`
	ExtParam  string            `json:"extParam"`
	Fees      string            `json:"fees"`
	To        map[string]string `json:"to"`
	FeeRate   string            `json:"feeRate"`
	SigParts  []SigParts        `json:"sigParts"`
	ErrorMsg  map[string]string `json:"errorMsg"`
	UserID    int64             `json:"userID"`
	Account   AssetsAccount     `json:"account"`
	DataType  int64             `json:"dataType"` // 1=transfer 2=summary
}

//easyjson:json
type SubmitRawTransactionReq struct {
	common.BaseReq
	PendingSignTx *adapter.PendingSignTx `json:"pendingSignTx"`
}

//easyjson:json
type SubmitRawTransactionRes struct {
	TxID        string   `json:"txID"`
	AccountID   string   `json:"accountID"`
	Coin        CoinInfo `json:"coin"`
	From        []string `json:"from"`
	To          []string `json:"to"`
	Amount      string   `json:"amount"`
	Decimal     int32    `json:"decimal"`
	TxType      uint64   `json:"txType"`
	TxAction    string   `json:"txAction"`
	Confirm     int64    `json:"confirm"`
	BlockHash   string   `json:"blockHash"`
	BlockHeight uint64   `json:"blockHeight"`
	IsMemo      bool     `json:"isMemo"`
	Memo        string   `json:"memo"`
	Fees        string   `json:"fees"`
	Received    bool     `json:"received"`
	SubmitTime  int64    `json:"submitTime"`
	ConfirmTime int64    `json:"confirmTime"`
	Status      string   `json:"status"`
	Reason      string   `json:"reason"`
	ExtParam    string   `json:"extParam"`
}

//easyjson:json
type GetTransactionFeeEstimatedReq struct {
	common.BaseReq
	Symbol        string `json:"symbol"`
	From          string `json:"from"`
	To            string `json:"to"`
	Amount        string `json:"amount"`
	Data          string `json:"data"`
	Decimals      int32  `json:"decimals"`
	TokenAddress  string `json:"tokenAddress"`
	TokenDecimals int32  `json:"tokenDecimals"`
}

//easyjson:json
type GetTransactionFeeEstimatedRes struct {
	GasLimit    string `json:"gasLimit"`
	GasPrice    string `json:"gasPrice"`
	Fee         string `json:"fee"`
	MainBalance string `json:"mainBalance"`
}

//easyjson:json
type GetTransactionCountOnChainReq struct {
	common.BaseReq
	Symbol  string `json:"symbol"`
	Address string `json:"address"`
}

//easyjson:json
type GetTransactionCountOnChainRes struct {
	Nonce uint64 `json:"nonce"`
}

//easyjson:json
type CreateTradeReq struct {
	common.BaseReq
	AccountID string            `json:"accountID"`
	Sid       string            `json:"sid"`
	Coin      CoinInfo          `json:"coin"`
	FeeRate   string            `json:"feeRate"`
	ExtParam  string            `json:"extParam"`
	Memo      string            `json:"memo"`
	To        map[string]string `json:"to"`
}

//easyjson:json
type CreateTradeRes struct {
	PendingSignTx []*adapter.PendingSignTx `json:"pendingSignTx"`
}

// CreateSummaryTxRes response for CreateSummaryTx (summary pending legs + optional fee deficits).
//
//easyjson:json
type CreateSummaryTxRes struct {
	SummaryPendingSignTx []*adapter.PendingSignTx     `json:"summaryPendingSignTx"`
	FeeDeficits          []*adapter.SummaryFeeDeficit `json:"feeDeficits"`
}

// EvaluateSummaryFeeDeficitReq re-estimates native fee gap for an existing summary leg.
//
//easyjson:json
type EvaluateSummaryFeeDeficitReq struct {
	common.BaseReq
	AccountID         string `json:"accountID"`
	SummarySid        string `json:"summarySid"`
	PreviousShortfall string `json:"previousShortfall"`
	PreviousReason    string `json:"previousReason"`
}

// CancelTransferTradeReq cancels a pending transfer using originSid (same nonce, zero-value self-send); requires a new sid.
//
//easyjson:json
type CancelTransferTradeReq struct {
	common.BaseReq
	AccountID      string `json:"accountID"`
	OriginSid      string `json:"originSid"`
	Sid            string `json:"sid"`
	FeeRate        string `json:"feeRate"`
	BaseFeeRate    string `json:"baseFeeRate"`
	FeeBumpPercent uint64 `json:"feeBumpPercent"`
	FeeBumpWei     string `json:"feeBumpWei"`
}

// SpeedUpTransferTradeReq replaces a pending transfer with a higher fee using originSid; requires a new sid.
//
//easyjson:json
type SpeedUpTransferTradeReq struct {
	common.BaseReq
	AccountID      string `json:"accountID"`
	OriginSid      string `json:"originSid"`
	Sid            string `json:"sid"`
	FeeRate        string `json:"feeRate"`
	BaseFeeRate    string `json:"baseFeeRate"`
	FeeBumpPercent uint64 `json:"feeBumpPercent"`
	FeeBumpWei     string `json:"feeBumpWei"`
}

// CreateSummaryTxReq creates a summary (sweep) transaction.
//
//easyjson:json
type CreateSummaryTxReq struct {
	common.BaseReq
	AccountID            string             `json:"accountID"`
	MinTransfer          string             `json:"minTransfer"`
	RetainedBalance      string             `json:"retainedBalance"`
	Address              string             `json:"address"`
	Coin                 CoinInfo           `json:"coin"`
	FeeRate              string             `json:"feeRate"`
	AddressStartIndex    int64              `json:"addressStartIndex"`
	AddressLimit         int64              `json:"addressLimit"`
	Confirms             int64              `json:"confirms"`
	FeesSupportAccount   FeesSupportAccount `json:"feesSupportAccount"`
	Memo                 string             `json:"memo"`
	Sid                  string             `json:"sid"`
}
