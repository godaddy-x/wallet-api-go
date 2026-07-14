package client

import (
	"testing"

	adapter "github.com/godaddy-x/wallet-adapter"
)

func TestValidateBTCPendingOutputsAcceptsCancelSelfPay(t *testing.T) {
	tx := &adapter.RawTransaction{
		Coin:   adapter.Coin{Symbol: "BTC"},
		To:     map[string]string{"bcrt1qchange": "0"},
		TxTo:   []string{"bcrt1qchange:0.00013058"},
		TxFrom: []string{"bcrt1qchange:50"},
	}
	if err := validateBTCPendingOutputs(tx, map[string]string{"bcrt1qchange": "0"}); err != nil {
		t.Fatalf("expected cancel self-pay allowed: %v", err)
	}
}

func TestValidateBTCPendingOutputsAcceptsDecimalFormat(t *testing.T) {
	tx := &adapter.RawTransaction{
		Coin: adapter.Coin{Symbol: "BTC"},
		To:   map[string]string{"bcrt1qrecipient": "0.00100000"},
		TxTo: []string{"bcrt1qrecipient:0.001"},
	}
	if err := validateBTCPendingOutputs(tx, map[string]string{"bcrt1qrecipient": "0.00100000"}); err != nil {
		t.Fatalf("expected decimal-equivalent match: %v", err)
	}
}

func TestValidateBTCPendingOutputsRejectsForeignChange(t *testing.T) {
	tx := &adapter.RawTransaction{
		Coin:   adapter.Coin{Symbol: "BTC"},
		To:     map[string]string{"bcrt1qrecipient": "0.0001"},
		TxTo:   []string{"bcrt1qrecipient:0.0001", "bcrt1qattacker:1.0"},
		TxFrom: []string{"bcrt1qchange:50"},
	}
	if err := validateBTCPendingOutputs(tx, map[string]string{"bcrt1qrecipient": "0.0001"}); err == nil {
		t.Fatal("expected foreign change rejection")
	}
}
