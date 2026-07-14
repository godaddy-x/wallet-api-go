package client

import (
	"encoding/hex"
	"testing"

	"github.com/godaddy-x/freego/utils"
	adapter "github.com/godaddy-x/wallet-adapter"
)

func TestValidateCreateTradeRequestHook(t *testing.T) {
	raw := &adapter.RawTransaction{
		Sid: "sid-1",
		Coin: adapter.Coin{
			Symbol: "ETH",
		},
		Account: &adapter.AssetsAccount{AccountID: "acc-1"},
		TxTo:    []string{"0xabc:0.01"},
	}
	data, err := utils.JsonMarshal(raw)
	if err != nil {
		t.Fatal(err)
	}
	pending := &adapter.PendingSignTx{Data: utils.Bytes2Str(data)}

	req := &CreateTradeReq{
		Sid:       "sid-1",
		AccountID: "acc-1",
		Coin:      CoinInfo{Symbol: "ETH"},
		To:        map[string]string{"0xabc": "0.01"},
	}
	ctx := &TradeCreatedContext{
		Kind:    TradeKindCreate,
		Request: req,
		Pending: []*adapter.PendingSignTx{pending},
	}
	if err := ValidateCreateTradeRequestHook()(ctx); err != nil {
		t.Fatalf("expected match: %v", err)
	}

	req.To["0xabc"] = "0.02"
	if err := ValidateCreateTradeRequestHook()(ctx); err == nil {
		t.Fatal("expected amount mismatch error")
	}
}

func TestValidateCreateTradeRequestHook_btcChangeOutput(t *testing.T) {
	raw := &adapter.RawTransaction{
		Sid: "sid-btc",
		Coin: adapter.Coin{
			Symbol: "BTC",
		},
		Account: &adapter.AssetsAccount{AccountID: "acc-1"},
		TxFrom:  []string{"bcrt1qay6v8dmyqu6lu6z448fx9re0c5nzy2ye22shua:50"},
		TxTo: []string{
			"bcrt1qzqwuj487l2weae4vqpdqfku5gk7ssj8h5ry6ec:0.0001",
			"bcrt1qay6v8dmyqu6lu6z448fx9re0c5nzy2ye22shua:41.79989000",
		},
	}
	data, err := utils.JsonMarshal(raw)
	if err != nil {
		t.Fatal(err)
	}
	pending := &adapter.PendingSignTx{Data: utils.Bytes2Str(data)}

	req := &CreateTradeReq{
		Sid:       "sid-btc",
		AccountID: "acc-1",
		Coin:      CoinInfo{Symbol: "BTC"},
		To:        map[string]string{"bcrt1qzqwuj487l2weae4vqpdqfku5gk7ssj8h5ry6ec": "0.0001"},
	}
	ctx := &TradeCreatedContext{
		Kind:    TradeKindCreate,
		Request: req,
		Pending: []*adapter.PendingSignTx{pending},
	}
	if err := ValidateCreateTradeRequestHook()(ctx); err != nil {
		t.Fatalf("expected BTC change output allowed: %v", err)
	}
}

func TestValidateCreateTradeRequestHook_rejectsExtraRecipient(t *testing.T) {
	raw := &adapter.RawTransaction{
		Sid: "sid-btc",
		Coin: adapter.Coin{
			Symbol: "BTC",
		},
		Account: &adapter.AssetsAccount{AccountID: "acc-1"},
		TxFrom:  []string{"bcrt1qay6v8dmyqu6lu6z448fx9re0c5nzy2ye22shua:50"},
		TxTo: []string{
			"bcrt1qzqwuj487l2weae4vqpdqfku5gk7ssj8h5ry6ec:0.0001",
			"bcrt1qattackerxxxxxxxxxxxxxxxxxxxxxxxxxx:0.5",
			"bcrt1qay6v8dmyqu6lu6z448fx9re0c5nzy2ye22shua:41.79989000",
		},
	}
	data, err := utils.JsonMarshal(raw)
	if err != nil {
		t.Fatal(err)
	}
	pending := &adapter.PendingSignTx{Data: utils.Bytes2Str(data)}
	req := &CreateTradeReq{
		Sid:       "sid-btc",
		AccountID: "acc-1",
		Coin:      CoinInfo{Symbol: "BTC"},
		To:        map[string]string{"bcrt1qzqwuj487l2weae4vqpdqfku5gk7ssj8h5ry6ec": "0.0001"},
	}
	ctx := &TradeCreatedContext{
		Kind:    TradeKindCreate,
		Request: req,
		Pending: []*adapter.PendingSignTx{pending},
	}
	if err := ValidateCreateTradeRequestHook()(ctx); err == nil {
		t.Fatal("expected extra recipient error")
	}
}

func TestValidatePendingDataSignHook(t *testing.T) {
	appKey := hex.EncodeToString([]byte("0123456789abcdef"))
	key, _ := hex.DecodeString(appKey)
	data := `{"sid":"s1"}`
	checkSign := utils.Base64Encode(utils.HMAC_SHA256_BASE(utils.Str2Bytes(data), key))
	pending := &adapter.PendingSignTx{
		Data:     data,
		DataSign: checkSign,
	}
	ctx := &TradeCreatedContext{Pending: []*adapter.PendingSignTx{pending}}
	if err := ValidatePendingDataSignHook(appKey)(ctx); err != nil {
		t.Fatalf("valid sign: %v", err)
	}
	pending.DataSign = "bad"
	if err := ValidatePendingDataSignHook(appKey)(ctx); err == nil {
		t.Fatal("expected invalid sign error")
	}
}

func TestRunTradeCreatedHooks_emptyPending(t *testing.T) {
	wc := &WalletClient{}
	err := wc.runTradeCreatedHooks(TradeKindCreate, nil, nil)
	if err == nil {
		t.Fatal("expected empty pending error")
	}
}
