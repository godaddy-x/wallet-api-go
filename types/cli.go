package types

import "github.com/godaddy-x/freego/node/common"

// CliFindWalletListReq lists wallets on MPC.
//
//easyjson:json
type CliFindWalletListReq struct {
	common.BaseReq
}

// CliWalletResult is a wallet entry returned by FindWalletList.
//
//easyjson:json
type CliWalletResult struct {
	ID           int64  `json:"id"`
	AppID        string `json:"appID"`
	WalletID     string `json:"walletID"`
	RootPath     string `json:"rootPath"`
	Alias        string `json:"alias"`
	Algorithm    string `json:"algorithm,omitempty"`
	AccountIndex int64  `json:"accountIndex"`
	CreateAt     int64  `json:"createAt"`
}

// CliFindWalletListRes is the FindWalletList response.
//
//easyjson:json
type CliFindWalletListRes struct {
	Result []CliWalletResult `json:"result"`
}

// CliCreateMPCWalletReq starts MPC wallet keygen.
//
//easyjson:json
type CliCreateMPCWalletReq struct {
	common.BaseReq
	Alias     string `json:"alias"`
	Algorithm string `json:"algorithm"` // ecdsa | ed25519
}

// CliCreateMPCWalletRes is the CreateMPCWallet response.
//
//easyjson:json
type CliCreateMPCWalletRes struct {
	WalletID  string `json:"walletID"`
	Algorithm string `json:"algorithm"`
}

// CliCreateAccountReq derives the next account for a wallet.
//
//easyjson:json
type CliCreateAccountReq struct {
	common.BaseReq
	WalletID  string `json:"walletID"`
	LastIndex int64  `json:"lastIndex"`
}

// CliCreateAccountRes is the CreateAccount response.
//
//easyjson:json
type CliCreateAccountRes struct {
	WalletID       string   `json:"walletID"`
	AccountID      string   `json:"accountID"`
	OtherOwnerKeys []string `json:"otherOwnerKeys"`
	ReqSigs        int64    `json:"reqSigs"`
	PublicKey      string   `json:"publicKey"`
	HdPath         string   `json:"hdPath"`
	AccountIndex   int64    `json:"accountIndex"`
	AddressIndex   int64    `json:"addressIndex"`
}

// CliCreateAddressReq derives addresses for an account.
//
//easyjson:json
type CliCreateAddressReq struct {
	common.BaseReq
	WalletID     string `json:"walletID"`
	AccountID    string `json:"accountID"`
	AccountIndex int64  `json:"accountIndex"`
	MainSymbol   string `json:"symbol"`
	LastIndex    int64  `json:"lastIndex"`
	Count        int64  `json:"count"`
	Change       int64  `json:"change"` // 0=external, 1=change
}

// CliAddressData is one derived address.
//
//easyjson:json
type CliAddressData struct {
	AddressIndex  int64  `json:"addressIndex"`
	AddressPubHex string `json:"addressPubHex"`
	HdPath        string `json:"hdPath"`
}

// CliCreateAddressRes is the CreateAddress response.
//
//easyjson:json
type CliCreateAddressRes struct {
	AddressList []CliAddressData `json:"addressList"`
}

// CliSignTransactionReq signs pending transaction data on MPC.
//
//easyjson:json
type CliSignTransactionReq struct {
	common.BaseReq
	Type      int64  `json:"type"` // 0=transfer, 1=summary, 2=smart contract
	Data      string `json:"data"`
	TradeSign string `json:"tradeSign"`
}

// CliSignTransactionRes is the SignTransaction response.
//
//easyjson:json
type CliSignTransactionRes struct {
	SignerList map[string]string `json:"signerList"`
}
