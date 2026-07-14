package client

import (
	"errors"
	"fmt"
	"strings"
)

func (c *WalletClient) requireMPC() error {
	if c == nil {
		return errors.New("wallet client is nil")
	}
	if c.mpcSDK == nil {
		return errors.New("mpc not connected")
	}
	return nil
}

func (c *WalletClient) sendCLI(path string, req, res interface{}) error {
	if isNilReq(req) {
		return ErrNilRequest
	}
	if err := c.requireMPC(); err != nil {
		return err
	}
	clampPageLimit(path, req)
	if err := sendWS(c.mpcSDK, path, req, res, c.wsTimeout); err != nil {
		name := strings.TrimPrefix(path, "/api/")
		return fmt.Errorf("%s: %w", name, err)
	}
	return nil
}

// FindWalletList calls MPC /api/FindWalletList.
func (c *WalletClient) FindWalletList(req *CliFindWalletListReq) (CliFindWalletListRes, error) {
	var res CliFindWalletListRes
	return res, c.sendCLI("/api/FindWalletList", req, &res)
}

// CreateMPCWallet calls MPC /api/CreateMPCWallet.
func (c *WalletClient) CreateMPCWallet(req *CliCreateMPCWalletReq) (CliCreateMPCWalletRes, error) {
	var res CliCreateMPCWalletRes
	return res, c.sendCLI("/api/CreateMPCWallet", req, &res)
}

// CliCreateAccount calls MPC /api/CreateAccount.
func (c *WalletClient) CliCreateAccount(req *CliCreateAccountReq) (CliCreateAccountRes, error) {
	var res CliCreateAccountRes
	return res, c.sendCLI("/api/CreateAccount", req, &res)
}

// CliCreateAddress calls MPC /api/CreateAddress.
func (c *WalletClient) CliCreateAddress(req *CliCreateAddressReq) (CliCreateAddressRes, error) {
	var res CliCreateAddressRes
	return res, c.sendCLI("/api/CreateAddress", req, &res)
}

// SignTransaction calls MPC /api/SignTransaction under the global sign lock.
func (c *WalletClient) SignTransaction(req *CliSignTransactionReq) (CliSignTransactionRes, error) {
	var res CliSignTransactionRes
	if isNilReq(req) {
		return res, ErrNilRequest
	}
	if err := c.requireMPC(); err != nil {
		return res, err
	}
	if err := signTransaction(c.mpcSDK, req, &res, c.wsTimeout); err != nil {
		return res, fmt.Errorf("SignTransaction: %w", err)
	}
	return res, nil
}
