// Package client is a minimal wallet SDK: connect OPS/MPC and run signed transfers.
package client

import (
	"errors"
	"reflect"
	"strings"
	"time"

	"github.com/godaddy-x/freego/utils"
	"github.com/godaddy-x/freego/utils/sdk"
)

// ErrNilRequest is returned when an API method is called with a nil request pointer.
var ErrNilRequest = errors.New("request is nil")

// SdkConfig holds OPS/MPC WebSocket connection and auth parameters.
type SdkConfig struct {
	Domain    string `json:"domain"`
	KeyPath   string `json:"keyPath"`
	LoginPath string `json:"loginPath"`
	Source    string `json:"source"`
	AppID     string `json:"appID"`
	AppKey    string `json:"appKey"`
	ClientPrk string `json:"clientPrk"`
	ServerPub string `json:"serverPub"`
	ClientNo  int64  `json:"clientNo"`
	SSL       bool   `json:"ssl"`
}

// Config holds OPS/MPC connection settings and the tenant appKey for dataSign validation.
type Config struct {
	Ops    SdkConfig
	MPC    SdkConfig
	AppKey string
	// WSTimeoutSec defaults to 300 when zero.
	WSTimeoutSec int64
	// TradeCreatedHooks run after OPS Create* trade APIs succeed, before MPC signing.
	// When AppKey is set, ValidatePendingDataSignHook and ValidateCreateTradeRequestHook
	// are prepended automatically unless DisableDefaultTradeHooks is true.
	TradeCreatedHooks []TradeCreatedHook
	// DisableDefaultTradeHooks skips built-in pending validation hooks.
	DisableDefaultTradeHooks bool
}

// ConfigHook mutates Config after load and before WebSocket connect.
type ConfigHook func(*Config) error

// WalletClient manages long-lived OPS and MPC WebSocket sessions.
type WalletClient struct {
	appKey            string
	wsTimeout         int64
	opsSDK            *sdk.SocketSDK
	mpcSDK            *sdk.SocketSDK
	tradeCreatedHooks []TradeCreatedHook
}

// New opens long-lived WebSocket sessions for OPS and/or MPC.
// A side is skipped when its SdkConfig is empty (Domain is blank).
func New(cfg Config, hooks ...ConfigHook) (*WalletClient, error) {
	for _, hook := range hooks {
		if hook == nil {
			continue
		}
		if err := hook(&cfg); err != nil {
			return nil, err
		}
	}
	normalizeConfig(&cfg)
	tmo := wsTimeout(cfg.WSTimeoutSec)
	wc := &WalletClient{
		appKey:    cfg.AppKey,
		wsTimeout: tmo,
	}
	if !cfg.DisableDefaultTradeHooks {
		wc.tradeCreatedHooks = append(wc.tradeCreatedHooks, defaultTradeCreatedHooks(cfg)...)
	}
	wc.tradeCreatedHooks = append(wc.tradeCreatedHooks, cfg.TradeCreatedHooks...)

	if sdkConfigEnabled(cfg.Ops) {
		opsSDK, err := newLongLivedSocket(cfg.Ops)
		if err != nil {
			return nil, err
		}
		wc.opsSDK = opsSDK
	}
	if sdkConfigEnabled(cfg.MPC) {
		mpcSDK, err := newLongLivedSocket(cfg.MPC)
		if err != nil {
			wc.Close()
			return nil, err
		}
		wc.mpcSDK = mpcSDK
	}
	return wc, nil
}

// NewFromFiles opens a WalletClient from ops.json / cli.json paths and appKey.
// Empty opsPath or mpcPath skips loading and connecting that side.
func NewFromFiles(opsPath, mpcPath, appKey string, wsTimeoutSec int64, hooks ...ConfigHook) (*WalletClient, error) {
	cfg := Config{
		AppKey:       appKey,
		WSTimeoutSec: wsTimeoutSec,
	}
	if strings.TrimSpace(opsPath) != "" {
		opsCfg, err := readConfigFile(opsPath)
		if err != nil {
			return nil, err
		}
		cfg.Ops = opsCfg
	}
	if strings.TrimSpace(mpcPath) != "" {
		mpcCfg, err := readConfigFile(mpcPath)
		if err != nil {
			return nil, err
		}
		cfg.MPC = mpcCfg
	}
	return New(cfg, hooks...)
}

// Connected reports whether OPS and MPC WebSocket sessions are up.
func (c *WalletClient) Connected() (ops bool, mpc bool) {
	if c == nil {
		return false, false
	}
	if c.opsSDK != nil {
		ops = c.opsSDK.IsWebSocketConnected()
	}
	if c.mpcSDK != nil {
		mpc = c.mpcSDK.IsWebSocketConnected()
	}
	return ops, mpc
}

// Close disconnects OPS and MPC sessions.
func (c *WalletClient) Close() {
	if c == nil {
		return
	}
	if c.opsSDK != nil {
		c.opsSDK.DisconnectWebSocket()
	}
	if c.mpcSDK != nil {
		c.mpcSDK.DisconnectWebSocket()
	}
}

// Transfer performs CreateTrade -> validate -> Sign -> SubmitTrade.
// Requires both OPS and MPC sessions; use New with both sides enabled.
func (c *WalletClient) Transfer(
	fromAccountID string,
	to map[string]string,
	symbol, contractAddress string,
) (SubmitResult, time.Duration, time.Duration, error) {
	if err := c.requireTransferSessions(); err != nil {
		return SubmitResult{}, 0, 0, err
	}
	signed, createMs, signMs, err := c.createAndSignRawTrade(fromAccountID, to, symbol, contractAddress)
	if err != nil {
		return SubmitResult{}, createMs, signMs, err
	}
	res, err := c.submitSignedRawTrade(signed)
	return res, createMs, signMs, err
}

func (c *WalletClient) requireTransferSessions() error {
	if c == nil {
		return errors.New("wallet client is nil")
	}
	if c.opsSDK == nil {
		return errors.New("ops not connected")
	}
	if c.mpcSDK == nil {
		return errors.New("mpc not connected")
	}
	return nil
}

func sdkConfigEnabled(cfg SdkConfig) bool {
	return strings.TrimSpace(cfg.Domain) != ""
}

// normalizeConfig keeps Config.AppKey and SdkConfig.AppKey aligned for login and Transfer validation.
func normalizeConfig(cfg *Config) {
	if cfg == nil {
		return
	}
	if strings.TrimSpace(cfg.AppKey) == "" && strings.TrimSpace(cfg.Ops.AppKey) != "" {
		cfg.AppKey = cfg.Ops.AppKey
	}
	if strings.TrimSpace(cfg.AppKey) == "" {
		return
	}
	if strings.TrimSpace(cfg.Ops.AppKey) == "" {
		cfg.Ops.AppKey = cfg.AppKey
	}
	if strings.TrimSpace(cfg.MPC.AppKey) == "" {
		cfg.MPC.AppKey = cfg.AppKey
	}
}

func readConfigFile(path string) (SdkConfig, error) {
	data, err := utils.ReadFile(path)
	if err != nil {
		return SdkConfig{}, err
	}
	cfg := SdkConfig{}
	if err := utils.JsonUnmarshal(data, &cfg); err != nil {
		return SdkConfig{}, err
	}
	return cfg, nil
}

func wsTimeout(sec int64) int64 {
	if sec > 0 {
		return sec
	}
	return 300
}

func sendWS(sdk *sdk.SocketSDK, path string, req, res interface{}, timeoutSec int64) error {
	if isNilReq(req) {
		return ErrNilRequest
	}
	return sdk.SendWebSocketMessage(path, req, res, true, true, timeoutSec)
}

func isNilReq(v interface{}) bool {
	if v == nil {
		return true
	}
	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Ptr, reflect.Map, reflect.Slice, reflect.Interface, reflect.Chan, reflect.Func:
		return rv.IsNil()
	default:
		return false
	}
}
