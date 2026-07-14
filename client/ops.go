package client

import (
	"errors"
	"fmt"
	"strings"

	adapter "github.com/godaddy-x/wallet-adapter"
)

func (c *WalletClient) requireOPS() error {
	if c == nil {
		return errors.New("wallet client is nil")
	}
	if c.opsSDK == nil {
		return errors.New("ops not connected")
	}
	return nil
}

func (c *WalletClient) sendOPS(path string, req, res interface{}) error {
	if isNilReq(req) {
		return ErrNilRequest
	}
	if err := c.requireOPS(); err != nil {
		return err
	}
	clampPageLimit(path, req)
	if err := sendWS(c.opsSDK, path, req, res, c.wsTimeout); err != nil {
		name := strings.TrimPrefix(path, "/api/")
		return fmt.Errorf("%s: %w", name, err)
	}
	return nil
}

// --- account ---

func (c *WalletClient) CreateAccount(req *CreateAccountReq) (CreateAccountRes, error) {
	var res CreateAccountRes
	return res, c.sendOPS("/api/CreateAccount", req, &res)
}

func (c *WalletClient) FindAccountByAccountID(req *FindAccountByAccountIDReq) (FindAccountByAccountIDRes, error) {
	var res FindAccountByAccountIDRes
	return res, c.sendOPS("/api/FindAccountByAccountID", req, &res)
}

func (c *WalletClient) FindAccountByWalletID(req *FindAccountByWalletIDReq) (FindAccountByWalletIDRes, error) {
	var res FindAccountByWalletIDRes
	return res, c.sendOPS("/api/FindAccountByWalletID", req, &res)
}

func (c *WalletClient) GetBalanceByAccount(req *GetBalanceByAccountReq) (GetBalanceByAccountRes, error) {
	var res GetBalanceByAccountRes
	return res, c.sendOPS("/api/GetBalanceByAccount", req, &res)
}

func (c *WalletClient) GetAccountBalanceList(req *GetAccountBalanceListReq) (GetAccountBalanceListRes, error) {
	var res GetAccountBalanceListRes
	return res, c.sendOPS("/api/GetAccountBalanceList", req, &res)
}

// --- address ---

func (c *WalletClient) ImportAddress(req *ImportAddressReq) (ImportAddressRes, error) {
	var res ImportAddressRes
	return res, c.sendOPS("/api/ImportAddress", req, &res)
}

func (c *WalletClient) FindAddressByAddress(req *FindAddressByAddressReq) (FindAddressByAddressRes, error) {
	var res FindAddressByAddressRes
	return res, c.sendOPS("/api/FindAddressByAddress", req, &res)
}

func (c *WalletClient) FindAddressByAccountID(req *FindAddressByAccountIDReq) (FindAddressByAccountIDRes, error) {
	var res FindAddressByAccountIDRes
	return res, c.sendOPS("/api/FindAddressByAccountID", req, &res)
}

func (c *WalletClient) VerifyAddress(req *VerifyAddressReq) (VerifyAddressRes, error) {
	var res VerifyAddressRes
	return res, c.sendOPS("/api/VerifyAddress", req, &res)
}

func (c *WalletClient) GetBalanceByAddress(req *GetBalanceByAddressReq) (GetBalanceByAddressRes, error) {
	var res GetBalanceByAddressRes
	return res, c.sendOPS("/api/GetBalanceByAddress", req, &res)
}

func (c *WalletClient) GetAddressBalanceList(req *GetAddressBalanceListReq) (GetAddressBalanceListRes, error) {
	var res GetAddressBalanceListRes
	return res, c.sendOPS("/api/GetAddressBalanceList", req, &res)
}

// --- contract ---

func (c *WalletClient) GetContracts(req *GetContractsReq) (GetContractsRes, error) {
	var res GetContractsRes
	return res, c.sendOPS("/api/GetContracts", req, &res)
}

func (c *WalletClient) GetContractTemplates(req *GetContractTemplatesReq) (GetContractTemplatesRes, error) {
	var res GetContractTemplatesRes
	return res, c.sendOPS("/api/GetContractTemplates", req, &res)
}

func (c *WalletClient) DeployContract(req *DeployContractReq) (DeployContractRes, error) {
	var res DeployContractRes
	if err := c.sendOPS("/api/DeployContract", req, &res); err != nil {
		return res, err
	}
	return res, c.runTradeCreatedHooks(TradeKindDeploy, req, res.PendingSignTx)
}

func (c *WalletClient) SubmitDeployContract(req *SubmitDeployContractReq) (SubmitDeployContractRes, error) {
	var res SubmitDeployContractRes
	return res, c.sendOPS("/api/SubmitDeployContract", req, &res)
}

func (c *WalletClient) SubmitSmartContractTrade(req *SubmitSmartContractTradeReq) (SubmitSmartContractTradeRes, error) {
	var res SubmitSmartContractTradeRes
	return res, c.sendOPS("/api/SubmitSmartContractTrade", req, &res)
}

// --- chain ---

func (c *WalletClient) SymbolBlockList(req *SymbolBlockListReq) (SymbolBlockListRes, error) {
	var res SymbolBlockListRes
	return res, c.sendOPS("/api/SymbolBlockList", req, &res)
}

func (c *WalletClient) GetBlockStatus(req *GetBlockStatusReq) (GetBlockStatusRes, error) {
	var res GetBlockStatusRes
	return res, c.sendOPS("/api/GetBlockStatus", req, &res)
}

// --- trade ---

func (c *WalletClient) CreateTrade(req *CreateTradeReq) (CreateTradeRes, error) {
	var res CreateTradeRes
	if err := c.sendOPS("/api/CreateTrade", req, &res); err != nil {
		return res, err
	}
	return res, c.runTradeCreatedHooks(TradeKindCreate, req, res.PendingSignTx)
}

func (c *WalletClient) SubmitTrade(req *SubmitRawTransactionReq) (SubmitRawTransactionRes, error) {
	var res SubmitRawTransactionRes
	return res, c.sendOPS("/api/SubmitTrade", req, &res)
}

func (c *WalletClient) SpeedUpTransferTrade(req *SpeedUpTransferTradeReq) (CreateTradeRes, error) {
	var res CreateTradeRes
	if err := c.sendOPS("/api/SpeedUpTransferTrade", req, &res); err != nil {
		return res, err
	}
	return res, c.runTradeCreatedHooks(TradeKindSpeedUp, req, res.PendingSignTx)
}

func (c *WalletClient) CancelTransferTrade(req *CancelTransferTradeReq) (CreateTradeRes, error) {
	var res CreateTradeRes
	if err := c.sendOPS("/api/CancelTransferTrade", req, &res); err != nil {
		return res, err
	}
	return res, c.runTradeCreatedHooks(TradeKindCancel, req, res.PendingSignTx)
}

func (c *WalletClient) CreateSummaryTx(req *CreateSummaryTxReq) (CreateSummaryTxRes, error) {
	var res CreateSummaryTxRes
	if err := c.sendOPS("/api/CreateSummaryTx", req, &res); err != nil {
		return res, err
	}
	return res, c.runTradeCreatedHooks(TradeKindSummary, req, res.SummaryPendingSignTx)
}

func (c *WalletClient) EvaluateSummaryFeeDeficit(req *EvaluateSummaryFeeDeficitReq) (adapter.SummaryFeeDeficitEvalResult, error) {
	var res adapter.SummaryFeeDeficitEvalResult
	if err := c.sendOPS("/api/EvaluateSummaryFeeDeficit", req, &res); err != nil {
		return res, err
	}
	return res, nil
}

func (c *WalletClient) FindTradeLog(req *FindTradeLogReq) (FindTradeLogRes, error) {
	var res FindTradeLogRes
	return res, c.sendOPS("/api/FindTradeLog", req, &res)
}

func (c *WalletClient) FindBalanceLog(req *FindBalanceLogReq) (FindBalanceLogRes, error) {
	var res FindBalanceLogRes
	return res, c.sendOPS("/api/FindBalanceLog", req, &res)
}

func (c *WalletClient) FindMonitorAlert(req *FindMonitorAlertReq) (FindMonitorAlertRes, error) {
	var res FindMonitorAlertRes
	return res, c.sendOPS("/api/FindMonitorAlert", req, &res)
}

// --- batch transfer ---

func (c *WalletClient) CreateBatchTransferTrade(req *CreateBatchTransferTradeReq) (CreateTradeRes, error) {
	var res CreateTradeRes
	if err := c.sendOPS("/api/CreateBatchTransferTrade", req, &res); err != nil {
		return res, err
	}
	return res, c.runTradeCreatedHooks(TradeKindBatchTransfer, req, res.PendingSignTx)
}

func (c *WalletClient) CreateBatchTransferApproveTrade(req *CreateBatchTransferApproveTradeReq) (CreateTradeRes, error) {
	var res CreateTradeRes
	if err := c.sendOPS("/api/CreateBatchTransferApproveTrade", req, &res); err != nil {
		return res, err
	}
	return res, c.runTradeCreatedHooks(TradeKindBatchApprove, req, res.PendingSignTx)
}

func (c *WalletClient) GetBatchTransferAllowance(req *GetBatchTransferAllowanceReq) (GetBatchTransferAllowanceRes, error) {
	var res GetBatchTransferAllowanceRes
	return res, c.sendOPS("/api/GetBatchTransferAllowance", req, &res)
}

// --- TRX Stake 2.0 ---

func (c *WalletClient) CreateStakeTrade(req *CreateStakeTradeReq) (CreateTradeRes, error) {
	var res CreateTradeRes
	if err := c.sendOPS("/api/CreateStakeTrade", req, &res); err != nil {
		return res, err
	}
	return res, c.runTradeCreatedHooks(TradeKindCreate, req, res.PendingSignTx)
}

func (c *WalletClient) CreateUnstakeTrade(req *CreateUnstakeTradeReq) (CreateTradeRes, error) {
	var res CreateTradeRes
	if err := c.sendOPS("/api/CreateUnstakeTrade", req, &res); err != nil {
		return res, err
	}
	return res, c.runTradeCreatedHooks(TradeKindCreate, req, res.PendingSignTx)
}

func (c *WalletClient) CreateWithdrawUnfreezeTrade(req *CreateWithdrawUnfreezeTradeReq) (CreateTradeRes, error) {
	var res CreateTradeRes
	if err := c.sendOPS("/api/CreateWithdrawUnfreezeTrade", req, &res); err != nil {
		return res, err
	}
	return res, c.runTradeCreatedHooks(TradeKindCreate, req, res.PendingSignTx)
}

func (c *WalletClient) GetAccountResourceDetail(req *GetAccountResourceDetailReq) (GetAccountResourceDetailRes, error) {
	var res GetAccountResourceDetailRes
	return res, c.sendOPS("/api/GetAccountResourceDetail", req, &res)
}

func (c *WalletClient) SpeedUpBatchTransferTrade(req *SpeedUpBatchTransferTradeReq) (CreateTradeRes, error) {
	var res CreateTradeRes
	if err := c.sendOPS("/api/SpeedUpBatchTransferTrade", req, &res); err != nil {
		return res, err
	}
	return res, c.runTradeCreatedHooks(TradeKindBatchSpeedUp, req, res.PendingSignTx)
}

func (c *WalletClient) CancelBatchTransferTrade(req *CancelBatchTransferTradeReq) (CreateTradeRes, error) {
	var res CreateTradeRes
	if err := c.sendOPS("/api/CancelBatchTransferTrade", req, &res); err != nil {
		return res, err
	}
	return res, c.runTradeCreatedHooks(TradeKindBatchCancel, req, res.PendingSignTx)
}

// --- wallet ---

func (c *WalletClient) CreateWallet(req *CreateWalletReq) (CreateWalletRes, error) {
	var res CreateWalletRes
	return res, c.sendOPS("/api/CreateWallet", req, &res)
}

func (c *WalletClient) FindWalletByWalletID(req *FindWalletByWalletIDReq) (FindWalletByWalletIDRes, error) {
	var res FindWalletByWalletIDRes
	return res, c.sendOPS("/api/FindWalletByWalletID", req, &res)
}

// --- subscribe ---

func (c *WalletClient) CreateSubscribe(req *CreateSubscribeReq) (CreateSubscribeRes, error) {
	var res CreateSubscribeRes
	return res, c.sendOPS("/api/CreateSubscribe", req, &res)
}
