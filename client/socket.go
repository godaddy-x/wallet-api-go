package client

import (
	"encoding/hex"
	"fmt"
	"strings"
	"sync"

	"github.com/godaddy-x/freego/utils"
	"github.com/godaddy-x/freego/utils/sdk"
	"github.com/godaddy-x/freego/zlog"
	"github.com/godaddy-x/wallet-api-go/types"
)

var cliSignMu sync.Mutex

func signTransaction(
	cliSDK *sdk.SocketSDK,
	req *CliSignTransactionReq,
	res *CliSignTransactionRes,
	wsTimeoutSec int64,
) error {
	cliSignMu.Lock()
	defer cliSignMu.Unlock()
	return cliSDK.SendWebSocketMessage("/api/SignTransaction", req, res, true, true, wsTimeoutSec)
}

func newLongLivedSocket(cfg SdkConfig) (*sdk.SocketSDK, error) {
	c := sdk.NewSocketSDK(cfg.Domain)
	c.SetClientNo(cfg.ClientNo)
	_ = c.SetMLDSA87Object(cfg.ClientNo, cfg.ClientPrk, cfg.ServerPub)
	c.SetSSL(cfg.SSL)

	resp, err := loginSocket(cfg)
	if err != nil {
		return nil, fmt.Errorf("auth appID=%s: %w", cfg.AppID, err)
	}
	c.AuthToken(resp)

	c.EnableReconnect()
	c.SetHealthPing(10)
	c.SetTokenExpiredCallback(func() {
		if token, err := loginSocket(cfg); err != nil {
			zlog.Error("sdk token refresh error", 0, zlog.String("errMsg", err.Error()))
			return
		} else {
			c.AuthToken(token)
		}
	})
	if err := c.ConnectWebSocket(); err != nil {
		return nil, err
	}
	return c, nil
}

func loginSocket(cfg SdkConfig) (sdk.AuthToken, error) {
	cliSignMu.Lock()
	defer cliSignMu.Unlock()

	loginSDK := sdk.NewSocketSDK(cfg.Domain)
	loginSDK.SetClientNo(cfg.ClientNo)
	_ = loginSDK.SetMLDSA87Object(cfg.ClientNo, cfg.ClientPrk, cfg.ServerPub)
	loginSDK.SetSSL(cfg.SSL)
	defer loginSDK.DisconnectWebSocket()

	req, err := loginRequestForConfig(cfg)
	if err != nil {
		return sdk.AuthToken{}, err
	}
	resp := sdk.AuthToken{}
	if err := loginSDK.LoginByWebSocketPlan2Auto(cfg.KeyPath, cfg.LoginPath, req, &resp, 10); err != nil {
		return sdk.AuthToken{}, err
	}
	return resp, nil
}

func loginRequestForConfig(cfg SdkConfig) (interface{}, error) {
	// OPS: AppID + AppKey HMAC login.
	// MPC: Plan2 source login — AppKey on cfg is for Transfer dataSign only.
	if strings.TrimSpace(cfg.AppID) != "" && strings.TrimSpace(cfg.AppKey) != "" {
		return newAppLoginRequest(cfg)
	}
	source := strings.TrimSpace(cfg.Source)
	if source == "" {
		source = "API"
	}
	return &types.CliPlan2LoginReq{Source: source}, nil
}

func newAppLoginRequest(cfg SdkConfig) (*types.AppLoginReq, error) {
	req := &types.AppLoginReq{
		AppID:  cfg.AppID,
		Nonce:  utils.Base64Encode(utils.GetRandomSecure(32)),
		Time:   utils.UnixSecond(),
		Source: "API",
	}
	h, err := hex.DecodeString(cfg.AppKey)
	if err != nil {
		return nil, fmt.Errorf("invalid appKey hex: %w", err)
	}
	req.Sign = utils.Base64Encode(utils.HMAC_SHA256_BASE(h, utils.Str2Bytes(utils.AddStr(req.Nonce, req.Time, req.Source))))
	return req, nil
}
