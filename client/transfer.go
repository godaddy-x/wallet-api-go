package client

import (
	"fmt"
	"time"

	"github.com/godaddy-x/freego/utils"
	adapter "github.com/godaddy-x/wallet-adapter"
)

// SubmitResult is an alias for SubmitRawTransactionRes returned by Transfer.
type SubmitResult = SubmitRawTransactionRes

func (c *WalletClient) createAndSignRawTrade(
	fromAccountID string,
	to map[string]string,
	symbol, contractAddress string,
) (*adapter.PendingSignTx, time.Duration, time.Duration, error) {
	sid := utils.GetUUID(true)
	req := CreateTradeReq{
		Sid:       sid,
		AccountID: fromAccountID,
		Coin: CoinInfo{
			Symbol:          symbol,
			ContractAddress: contractAddress,
		},
		To: to,
	}
	createStart := time.Now()
	res, err := c.CreateTrade(&req)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("CreateTrade: %w", err)
	}
	createElapsed := time.Since(createStart)

	txData := res.PendingSignTx[0]

	signReq := CliSignTransactionReq{Type: 0, Data: txData.Data, TradeSign: txData.TradeSign}
	signRes := CliSignTransactionRes{}
	signStart := time.Now()
	if err := signTransaction(c.mpcSDK, &signReq, &signRes, c.wsTimeout); err != nil {
		return nil, createElapsed, time.Since(signStart), fmt.Errorf("SignTransaction sid=%s: %w", sid, err)
	}
	signElapsed := time.Since(signStart)
	if len(signRes.SignerList) == 0 {
		return nil, createElapsed, signElapsed, fmt.Errorf("cli signer is nil sid=%s", sid)
	}
	txData.SignerList = signRes.SignerList
	return txData, createElapsed, signElapsed, nil
}

func (c *WalletClient) submitSignedRawTrade(signed *adapter.PendingSignTx) (SubmitRawTransactionRes, error) {
	req := SubmitRawTransactionReq{PendingSignTx: signed}
	return c.SubmitTrade(&req)
}
