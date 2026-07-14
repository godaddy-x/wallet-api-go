package client

import (
	"encoding/hex"
	"errors"
	"fmt"
	"strings"

	"github.com/godaddy-x/freego/utils"
	adapter "github.com/godaddy-x/wallet-adapter"
)

// TradeKind identifies which OPS Create* API produced the pending transaction.
type TradeKind string

const (
	TradeKindCreate        TradeKind = "create_trade"
	TradeKindSpeedUp       TradeKind = "speed_up"
	TradeKindCancel        TradeKind = "cancel"
	TradeKindSummary       TradeKind = "summary"
	TradeKindBatchTransfer TradeKind = "batch_transfer"
	TradeKindBatchApprove  TradeKind = "batch_approve"
	TradeKindBatchSpeedUp  TradeKind = "batch_speed_up"
	TradeKindBatchCancel   TradeKind = "batch_cancel"
	TradeKindDeploy        TradeKind = "deploy"
)

// TradeCreatedHook runs after OPS successfully builds pending sign transaction(s)
// and before MPC signing. Return a non-nil error to abort the flow.
type TradeCreatedHook func(ctx *TradeCreatedContext) error

// TradeCreatedContext carries the original request and OPS-built pending txs.
type TradeCreatedContext struct {
	Kind    TradeKind
	Request any
	Pending []*adapter.PendingSignTx
	Client  *WalletClient
}

// FirstPending returns the first pending sign tx or an error when the list is empty.
func (ctx *TradeCreatedContext) FirstPending() (*adapter.PendingSignTx, error) {
	if ctx == nil || len(ctx.Pending) == 0 {
		return nil, errors.New("empty pendingSignTx")
	}
	return ctx.Pending[0], nil
}

// DecodeRaw unmarshals pending data as RawTransaction (MPC sign type=0).
func (ctx *TradeCreatedContext) DecodeRaw() (*adapter.RawTransaction, error) {
	pending, err := ctx.FirstPending()
	if err != nil {
		return nil, err
	}
	tx := &adapter.RawTransaction{}
	if err := utils.JsonUnmarshal(utils.Str2Bytes(pending.Data), tx); err != nil {
		return nil, fmt.Errorf("unmarshal raw tx: %w", err)
	}
	return tx, nil
}

// DecodeSmart unmarshals pending data as SmartContractRawTransaction (MPC sign type=2).
func (ctx *TradeCreatedContext) DecodeSmart() (*adapter.SmartContractRawTransaction, error) {
	pending, err := ctx.FirstPending()
	if err != nil {
		return nil, err
	}
	tx := &adapter.SmartContractRawTransaction{}
	if err := utils.JsonUnmarshal(utils.Str2Bytes(pending.Data), tx); err != nil {
		return nil, fmt.Errorf("unmarshal smart contract tx: %w", err)
	}
	return tx, nil
}

// AddTradeCreatedHook appends a hook invoked on subsequent Create* trade calls.
func (c *WalletClient) AddTradeCreatedHook(h TradeCreatedHook) {
	if c == nil || h == nil {
		return
	}
	c.tradeCreatedHooks = append(c.tradeCreatedHooks, h)
}

func (c *WalletClient) runTradeCreatedHooks(kind TradeKind, req any, pending []*adapter.PendingSignTx) error {
	if c == nil {
		return errors.New("wallet client is nil")
	}
	if len(pending) == 0 {
		return errors.New("empty pendingSignTx")
	}
	ctx := &TradeCreatedContext{
		Kind:    kind,
		Request: req,
		Pending: pending,
		Client:  c,
	}
	for _, h := range c.tradeCreatedHooks {
		if h == nil {
			continue
		}
		if err := h(ctx); err != nil {
			return err
		}
	}
	return nil
}

// ValidatePendingDataSignHook verifies OPS dataSign (HMAC with appKey) on each pending tx.
func ValidatePendingDataSignHook(appKey string) TradeCreatedHook {
	key, err := hex.DecodeString(appKey)
	if err != nil {
		return func(*TradeCreatedContext) error {
			return fmt.Errorf("invalid appKey: %w", err)
		}
	}
	return func(ctx *TradeCreatedContext) error {
		for i, tx := range ctx.Pending {
			if err := validatePendingDataSign(key, tx); err != nil {
				return fmt.Errorf("pending[%d]: %w", i, err)
			}
		}
		return nil
	}
}

// ValidateCreateTradeRequestHook matches CreateTradeReq fields (to, amount, account, symbol, contract) to pending raw tx.
// No-op for other TradeKind values or create-style APIs that use a different request type (e.g. stake).
func ValidateCreateTradeRequestHook() TradeCreatedHook {
	return func(ctx *TradeCreatedContext) error {
		if ctx.Kind != TradeKindCreate {
			return nil
		}
		req, ok := ctx.Request.(*CreateTradeReq)
		if !ok || req == nil {
			return nil
		}
		tx, err := ctx.DecodeRaw()
		if err != nil {
			return err
		}
		return validateRawTradePending(
			tx,
			req.Sid,
			req.AccountID,
			req.Coin.Symbol,
			req.Coin.ContractAddress,
			req.To,
		)
	}
}

func validatePendingDataSign(appKey []byte, tx *adapter.PendingSignTx) error {
	if tx == nil {
		return errors.New("pendingSignTx is nil")
	}
	checkSign := utils.HMAC_SHA256_BASE(utils.Str2Bytes(tx.Data), appKey)
	if !utils.CompareBase64Sign(checkSign, tx.DataSign) {
		return errors.New("tx data check sign invalid")
	}
	return nil
}

// ValidateCreateSummaryTxRequestHook ensures summary pending txs sweep to the requested target address.
func ValidateCreateSummaryTxRequestHook() TradeCreatedHook {
	return func(ctx *TradeCreatedContext) error {
		if ctx.Kind != TradeKindSummary {
			return nil
		}
		req, ok := ctx.Request.(*CreateSummaryTxReq)
		if !ok || req == nil {
			return fmt.Errorf("CreateSummaryTx hook: invalid request type")
		}
		target := strings.TrimSpace(req.Address)
		if target == "" {
			return fmt.Errorf("summary target address is empty")
		}
		for i, pending := range ctx.Pending {
			tx, err := decodePendingRaw(pending)
			if err != nil {
				return fmt.Errorf("pending[%d]: %w", i, err)
			}
			if err := validateSummaryRawPending(tx, req.Sid, req.AccountID, req.Coin.Symbol, target); err != nil {
				return fmt.Errorf("pending[%d]: %w", i, err)
			}
		}
		return nil
	}
}

// ValidateRBFSidRequestHook binds speed-up/cancel pending txs to the caller request sid/account.
func ValidateRBFSidRequestHook() TradeCreatedHook {
	return func(ctx *TradeCreatedContext) error {
		var sid, accountID string
		switch ctx.Kind {
		case TradeKindSpeedUp:
			req, ok := ctx.Request.(*SpeedUpTransferTradeReq)
			if !ok || req == nil {
				return fmt.Errorf("SpeedUp hook: invalid request type")
			}
			sid, accountID = req.Sid, req.AccountID
		case TradeKindCancel:
			req, ok := ctx.Request.(*CancelTransferTradeReq)
			if !ok || req == nil {
				return fmt.Errorf("Cancel hook: invalid request type")
			}
			sid, accountID = req.Sid, req.AccountID
		default:
			return nil
		}
		for i, pending := range ctx.Pending {
			tx, err := decodePendingRaw(pending)
			if err != nil {
				return fmt.Errorf("pending[%d]: %w", i, err)
			}
			if tx.Sid != sid {
				return fmt.Errorf("pending[%d]: sid invalid, got=%s want=%s", i, tx.Sid, sid)
			}
			if tx.Account == nil || tx.Account.AccountID != accountID {
				got := ""
				if tx.Account != nil {
					got = tx.Account.AccountID
				}
				return fmt.Errorf("pending[%d]: accountID invalid, got=%s want=%s", i, got, accountID)
			}
		}
		return nil
	}
}

func decodePendingRaw(pending *adapter.PendingSignTx) (*adapter.RawTransaction, error) {
	if pending == nil {
		return nil, errors.New("pendingSignTx is nil")
	}
	tx := &adapter.RawTransaction{}
	if err := utils.JsonUnmarshal(utils.Str2Bytes(pending.Data), tx); err != nil {
		return nil, fmt.Errorf("unmarshal raw tx: %w", err)
	}
	return tx, nil
}

func validateSummaryRawPending(
	tx *adapter.RawTransaction,
	sid, accountID, symbol, target string,
) error {
	if tx == nil {
		return errors.New("raw tx is nil")
	}
	if !strings.HasPrefix(tx.Sid, sid) {
		return fmt.Errorf("sid invalid, got=%s want prefix=%s", tx.Sid, sid)
	}
	if tx.Coin.Symbol != symbol {
		return fmt.Errorf("symbol invalid")
	}
	if tx.Account == nil || tx.Account.AccountID != accountID {
		got := ""
		if tx.Account != nil {
			got = tx.Account.AccountID
		}
		return fmt.Errorf("accountID invalid, got=%s want=%s", got, accountID)
	}
	if len(tx.To) != 1 {
		return fmt.Errorf("summary tx target count=%d want 1", len(tx.To))
	}
	for addr := range tx.To {
		if !strings.EqualFold(strings.TrimSpace(addr), target) {
			return fmt.Errorf("summary target mismatch, got=%s want=%s", addr, target)
		}
	}
	return nil
}

func defaultTradeCreatedHooks(cfg Config) []TradeCreatedHook {
	if cfg.AppKey == "" {
		return nil
	}
	return []TradeCreatedHook{
		ValidatePendingDataSignHook(cfg.AppKey),
		ValidateCreateTradeRequestHook(),
		ValidateCreateSummaryTxRequestHook(),
		ValidateRBFSidRequestHook(),
	}
}
