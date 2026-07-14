package client

import (
	"fmt"
	"strings"

	"github.com/godaddy-x/freego/utils/decimal"
	adapter "github.com/godaddy-x/wallet-adapter"
)

func validateRawTradePending(
	tx *adapter.RawTransaction,
	sid, fromAccountID, symbol, contractAddress string,
	to map[string]string,
) error {
	if tx.Sid != sid {
		return fmt.Errorf("sid invalid")
	}
	if tx.Coin.Symbol != symbol {
		return fmt.Errorf("symbol invalid")
	}
	if tx.Coin.Contract.Address != contractAddress {
		return fmt.Errorf("contractAddress invalid")
	}
	if tx.Account.AccountID != fromAccountID {
		return fmt.Errorf("accountID invalid, got=%s want=%s", tx.Account.AccountID, fromAccountID)
	}
	if strings.EqualFold(strings.TrimSpace(symbol), "BTC") {
		return validateBTCPendingOutputs(tx, to)
	}
	return validateExactPendingOutputs(tx, to, maxAllowedExtraOutputs(symbol))
}

func validateExactPendingOutputs(tx *adapter.RawTransaction, to map[string]string, maxExtra int) error {
	wantLines := make(map[string]struct{}, len(to))
	for addr, amt := range to {
		wantLines[strings.TrimSpace(addr)+":"+strings.TrimSpace(amt)] = struct{}{}
	}
	found := make(map[string]struct{}, len(wantLines))
	for _, line := range tx.TxTo {
		if _, ok := wantLines[line]; ok {
			found[line] = struct{}{}
		}
	}
	if len(found) != len(wantLines) {
		for line := range wantLines {
			if _, ok := found[line]; !ok {
				return fmt.Errorf("missing tx.To entry %q in %v", line, tx.TxTo)
			}
		}
	}
	extra := len(tx.TxTo) - len(found)
	if extra > maxExtra {
		return fmt.Errorf("unexpected extra outputs in tx.TxTo: got %d entries, want %d (+ at most %d change)",
			len(tx.TxTo), len(wantLines), maxExtra)
	}
	return nil
}

func validateBTCPendingOutputs(tx *adapter.RawTransaction, to map[string]string) error {
	want := make(map[string]string, len(to))
	for addr, amt := range to {
		want[strings.TrimSpace(addr)] = strings.TrimSpace(amt)
	}
	found := make(map[string]struct{}, len(want))
	var extraLines []string
	for _, line := range tx.TxTo {
		addr, amt, ok := parseAddrAmountLine(line)
		if !ok {
			return fmt.Errorf("invalid tx.To line %q", line)
		}
		if wantAmt, ok := want[addr]; ok {
			if isZeroAmount(wantAmt) {
				inputAddrs := txFromAddresses(tx.TxFrom)
				if !inputAddrs[addr] {
					return fmt.Errorf("cancel output not owned by input addresses: %s", addr)
				}
				found[addr] = struct{}{}
				continue
			}
			if !amountStringsEqual(wantAmt, amt) {
				return fmt.Errorf("tx.To amount mismatch for %s: got %q want %q", addr, amt, wantAmt)
			}
			found[addr] = struct{}{}
			continue
		}
		extraLines = append(extraLines, line)
	}
	for addr := range want {
		if _, ok := found[addr]; !ok {
			return fmt.Errorf("missing tx.To entry for %q in %v", addr, tx.TxTo)
		}
	}
	if len(extraLines) > 1 {
		return fmt.Errorf("unexpected extra outputs in tx.TxTo: got %d entries, want %d (+ at most 1 change)",
			len(tx.TxTo), len(want))
	}
	if len(extraLines) == 1 {
		inputAddrs := txFromAddresses(tx.TxFrom)
		addr, _, ok := parseAddrAmountLine(extraLines[0])
		if !ok {
			return fmt.Errorf("invalid change output line %q", extraLines[0])
		}
		if !inputAddrs[addr] {
			return fmt.Errorf("change output not owned by input addresses: %s", addr)
		}
	}
	return nil
}

func parseAddrAmountLine(line string) (addr, amt string, ok bool) {
	line = strings.TrimSpace(line)
	idx := strings.LastIndex(line, ":")
	if idx <= 0 {
		return "", "", false
	}
	return strings.TrimSpace(line[:idx]), strings.TrimSpace(line[idx+1:]), true
}

func amountStringsEqual(a, b string) bool {
	da, errA := decimal.NewFromString(strings.TrimSpace(a))
	db, errB := decimal.NewFromString(strings.TrimSpace(b))
	if errA == nil && errB == nil {
		return da.Equal(db)
	}
	return strings.TrimSpace(a) == strings.TrimSpace(b)
}

func isZeroAmount(v string) bool {
	d, err := decimal.NewFromString(strings.TrimSpace(v))
	return err == nil && d.Equal(decimal.Zero)
}

func txFromAddresses(txFrom []string) map[string]bool {
	out := make(map[string]bool, len(txFrom))
	for _, line := range txFrom {
		addr, _, ok := strings.Cut(line, ":")
		if !ok {
			continue
		}
		addr = strings.TrimSpace(addr)
		if addr != "" {
			out[addr] = true
		}
	}
	return out
}

func maxAllowedExtraOutputs(symbol string) int {
	if strings.EqualFold(strings.TrimSpace(symbol), "BTC") {
		return 1
	}
	return 0
}
