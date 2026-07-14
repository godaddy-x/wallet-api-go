# Wallet API SDK (Go)

Go SDK for **OPS** and **MPC** WebSocket integration: signed transfers, wallet lifecycle, and MPC signing.

Architecture mirrors [wallet-api-java](https://github.com/godaddy-x/wallet-api-java) and [wallet-api-rust](https://github.com/godaddy-x/wallet-api-rust):

```
types/         — OPS/MPC request and response DTOs (source of truth for other SDKs)
client/        — WalletClient, TradeHook, LogWatch, WebSocket helpers
integration/   — live and offline integration tests
config/        — sample ops.json / cli.json
```

WebSocket transport and Plan2 crypto are provided by [freego](https://github.com/godaddy-x/freego) (`github.com/godaddy-x/freego/utils/sdk`).

## Requirements

- Go **1.26+**
- Reachable **OPS** and/or **MPC broker** endpoints

---

## Install

### From module proxy (recommended)

Pin a release version in production:

```bash
go get github.com/godaddy-x/wallet-api-go@v0.1.0
```

Prefer an exact module version in `go.mod` when integrating against a live OPS/MPC stack.

### From GitHub

Same module path; pin a tag or commit instead of a release version:

```bash
# By release tag
go get github.com/godaddy-x/wallet-api-go@v0.1.0

# By commit
go get github.com/godaddy-x/wallet-api-go@abcdef1234567890
```

### Local development

```bash
git clone https://github.com/godaddy-x/wallet-api-go.git
cd wallet-api-go
go test ./...
```

---

## Third-party integration

### Credentials from the platform operator

The SDK does not mint credentials. Request the following from whoever operates OPS/MPC:

| Item | Used for |
|------|----------|
| OPS host (`domain`) | OPS WebSocket endpoint |
| MPC host (`domain`) | MPC WebSocket endpoint |
| `appID` | OPS app identity |
| `appKey` (hex) | OPS login HMAC + transfer `dataSign` validation |
| `clientNo` | ML-DSA user id on OPS |
| `clientPrk` / `serverPub` | ML-DSA key pair for Plan2 login |
| MPC `clientNo`, `clientPrk`, `serverPub` | MPC broker session |
| IP whitelist entry | Required for non-loopback OPS clients |

Store secrets in environment variables or a secret manager. Never commit `appKey`, `clientPrk`, or private keys.

Broker binding: MPC `source` should match your `appID`; broker `tradeKey` must align with the OPS tenant.

### Configuration files

Sample shapes ship in `config/`. Replace values with credentials issued for your app.

**OPS — `config/ops.json`**

```json
{
  "domain": "ops.example.com",
  "keyPath": "/api/PublicKey",
  "loginPath": "/api/Login",
  "appID": "your-app-id",
  "appKey": "hex-encoded-app-key",
  "clientNo": "202606271558577948",
  "ssl": false,
  "clientPrk": "base64-mldsa-seed-32-bytes",
  "serverPub": "base64-mldsa-server-public-key"
}
```

**MPC — `config/cli.json`**

```json
{
  "domain": "mpc.example.com",
  "keyPath": "/api/PublicKey",
  "loginPath": "/api/Login",
  "clientNo": 3,
  "source": "your-app-id",
  "clientPrk": "base64-mldsa-seed-32-bytes",
  "serverPub": "base64-mldsa-server-public-key",
  "ssl": false
}
```

Notes:

- `domain` is the host (and port, if required) issued by the platform operator.
- OPS login uses `appID` + `appKey` (HMAC: **data** = app key bytes, **key** = `nonce + time + source`).
- MPC Plan2 login uses `source` when no app key is configured on the MPC side.
- `clientPrk` is a **32-byte ML-DSA seed** (base64), not an expanded secret key.
- Large `clientNo` values **should** be JSON strings in `ops.json` to avoid precision issues in other tooling; Go parses string or number safely.

### Quick start — connect and transfer

```go
package main

import (
    "log"
    "os"

    "github.com/godaddy-x/wallet-api-go/client"
)

func main() {
    appKey := os.Getenv("APP_KEY")
    wc, err := client.NewFromFiles("config/ops.json", "config/cli.json", appKey, 300)
    if err != nil {
        log.Fatal(err)
    }
    defer wc.Close()

    opsOK, mpcOK := wc.Connected()
    if !opsOK || !mpcOK {
        log.Fatal("OPS and MPC must both be connected")
    }

    res, createMs, signMs, err := wc.Transfer(
        "account-id",
        map[string]string{"0xabc...": "0.01"},
        "ETH",
        "", // token contract; empty for native coin
    )
    if err != nil {
        log.Fatal(err)
    }
    log.Printf("txID=%s createMs=%v signMs=%v", res.TxID, createMs, signMs)
}
```

`Transfer` runs **OPS CreateTrade → MPC SignTransaction → OPS SubmitTrade**.

Connect only one side with an empty path:

```go
opsOnly, _ := client.NewFromFiles("config/ops.json", "", appKey, 300)
mpcOnly, _ := client.NewFromFiles("", "config/cli.json", "", 300)
```

`Transfer` requires **both** OPS and MPC sessions.

### Programmatic configuration

```go
wc, err := client.New(client.Config{
    AppKey:       os.Getenv("APP_KEY"),
    WSTimeoutSec: 300,
    Ops: client.SdkConfig{
        Domain:    os.Getenv("OPS_DOMAIN"),
        AppID:     os.Getenv("APP_ID"),
        AppKey:    os.Getenv("APP_KEY"),
        ClientNo:  clientNo,
        ClientPrk: os.Getenv("CLIENT_PRK"),
        ServerPub: os.Getenv("SERVER_PUB"),
        SSL:       false,
    },
    MPC: client.SdkConfig{
        Domain:    os.Getenv("MPC_DOMAIN"),
        ClientNo:  3,
        Source:    os.Getenv("APP_ID"),
        ClientPrk: os.Getenv("MPC_CLIENT_PRK"),
        ServerPub: os.Getenv("MPC_SERVER_PUB"),
        SSL:       false,
    },
})
```

Runtime overrides via config hook:

```go
wc, err := client.NewFromFiles(opsPath, mpcPath, appKey, 300, func(cfg *client.Config) error {
    if v := os.Getenv("OPS_DOMAIN"); v != "" {
        cfg.Ops.Domain = v
    }
    if v := os.Getenv("MPC_DOMAIN"); v != "" {
        cfg.MPC.Domain = v
    }
    if v := os.Getenv("APP_KEY"); v != "" {
        cfg.AppKey = v
        cfg.Ops.AppKey = v
    }
    return nil
})
```

### Wallet onboarding

Typical flow to register a wallet, account, and on-chain address:

```go
wallets, err := wc.FindWalletList(&client.CliFindWalletListReq{})
if err != nil || len(wallets.Result) == 0 {
    // handle
}
w := wallets.Result[0]

_, err = wc.CreateWallet(&client.CreateWalletReq{
    WalletID:  w.WalletID,
    RootPath:  w.RootPath,
    Alias:     "my-wallet",
    Algorithm: w.Algorithm,
})

acc, err := wc.CliCreateAccount(&client.CliCreateAccountReq{
    WalletID:  w.WalletID,
    LastIndex: -1,
})

_, err = wc.CreateAccount(&client.CreateAccountReq{
    WalletID:     w.WalletID,
    AccountID:    acc.AccountID,
    Alias:        "my-account",
    Symbol:       "BTC",
    PublicKey:    acc.PublicKey,
    AccountIndex: acc.AccountIndex,
    HdPath:       acc.HdPath,
    ReqSigs:      acc.ReqSigs,
})

addrs, err := wc.CliCreateAddress(&client.CliCreateAddressReq{
    WalletID:     w.WalletID,
    AccountID:    acc.AccountID,
    AccountIndex: acc.AccountIndex,
    MainSymbol:   "BTC",
    LastIndex:    -1,
    Count:        1,
    Change:       0,
})

item := addrs.AddressList[0]
_, err = wc.ImportAddress(&client.ImportAddressReq{
    AccountID: acc.AccountID,
    Addresses: []client.ImportAddressItem{
        {
            AddrIndex: item.AddressIndex,
            PublicKey: item.AddressPubHex,
            HdPath:    item.HdPath,
        },
    },
})
```

### Types and imports

All OPS/MPC request and response bodies are **strongly typed** DTOs in `github.com/godaddy-x/wallet-api-go/types`. The `client` package re-exports them as type aliases for convenience.

```go
import (
    "github.com/godaddy-x/wallet-api-go/client"
    "github.com/godaddy-x/wallet-api-go/types"
)

// DTOs: client.CreateTradeReq == types.CreateTradeReq
var _ types.CreateTradeReq
```

Low-level WebSocket transport: `github.com/godaddy-x/freego/utils/sdk.SocketSDK` (used internally by `client`).

### Security

- Do not log or persist `appKey`, `clientPrk`, or JWT secrets.
- OPS enforces client IP whitelist for non-loopback callers.
- Session: Plan2 login → JWT → encrypted business messages (`p=1`).
- With `AppKey` set, default trade hooks verify `pendingSignTx.dataSign` before MPC signing.
- Call `Close()` on shutdown.

### Versioning

1. Pin an exact module version in `go.mod`.
2. Read release notes before upgrading.
3. Run integration tests against your OPS/MPC stack after any bump.

---

## Protocol

| Phase | Plan | Transport |
|-------|------|-----------|
| Plan2 bootstrap + login | `p=2` | WebSocket `/ws` + ML-KEM + ML-DSA + AES-GCM |
| Post-login business | `p=1` | JWT `Authorization` + AES-GCM body |
| Heartbeat | `p=0` | `/ws/ping` |

Additional details:

- `SHARED_INFO = "freego-ecdh-aes-gcm"` for HKDF key derivation
- Empty-key HMAC uses zero-length key (matches Java BouncyCastle behavior)

---

## API surface

| Area | Methods |
|------|---------|
| Lifecycle | `New`, `NewFromFiles`, `Connected`, `Close`, `Transfer` |
| OPS | `CreateWallet`, `CreateAccount`, `CreateTrade`, `SubmitTrade`, `ImportAddress`, … |
| MPC | `FindWalletList`, `CreateMPCWallet`, `CliCreateAccount`, `CliCreateAddress`, `SignTransaction` |
| Log watch | `WatchTradeLog`, `WatchBalanceLog`, `WatchMonitorAlert` |
| Trade hooks | `AddTradeCreatedHook`, `ValidatePendingDataSignHook`, … |

All methods require a non-nil request pointer (`&client.XxxReq{}` at minimum). `nil` returns `client.ErrNilRequest`.

### OPS detail

| Category | Methods |
|----------|---------|
| Account | `CreateAccount`, `FindAccountByAccountID`, `FindAccountByWalletID`, `GetBalanceByAccount`, `GetAccountBalanceList` |
| Address | `ImportAddress`, `FindAddressByAddress`, `FindAddressByAccountID`, `VerifyAddress`, `GetBalanceByAddress`, `GetAddressBalanceList` |
| Contract | `GetContracts`, `GetContractTemplates`, `DeployContract`, `SubmitDeployContract`, `SubmitSmartContractTrade` |
| Chain | `SymbolBlockList`, `GetBlockStatus` |
| Trade | `CreateTrade`, `SubmitTrade`, `SpeedUpTransferTrade`, `CancelTransferTrade`, `CreateSummaryTx`, `EvaluateSummaryFeeDeficit`, `FindTradeLog`, `FindBalanceLog`, `FindMonitorAlert` |
| Batch | `CreateBatchTransferTrade`, `CreateBatchTransferApproveTrade`, `GetBatchTransferAllowance`, `SpeedUpBatchTransferTrade`, `CancelBatchTransferTrade` |
| Stake | `CreateStakeTrade`, `CreateUnstakeTrade`, `CreateWithdrawUnfreezeTrade`, `GetAccountResourceDetail` |
| Wallet | `CreateWallet`, `FindWalletByWalletID` |
| Subscribe | `CreateSubscribe` |

### MPC detail

| Method | Path |
|--------|------|
| `FindWalletList` | `/api/FindWalletList` |
| `CreateMPCWallet` | `/api/CreateMPCWallet` |
| `CliCreateAccount` | `/api/CreateAccount` |
| `CliCreateAddress` | `/api/CreateAddress` |
| `SignTransaction` | `/api/SignTransaction` |

When MPC runs in signing-only mode, only `SignTransaction` may be enabled server-side.

---

## Trade hooks

When `AppKey` is configured, default hooks validate OPS `pendingSignTx` before MPC signing:

```go
wc, err := client.New(client.Config{
    AppKey:                   appKey,
    Ops:                      opsCfg,
    MPC:                      mpcCfg,
    DisableDefaultTradeHooks: false,
})

wc.AddTradeCreatedHook(func(ctx *client.TradeCreatedContext) error {
    // ctx.Kind, ctx.Pending, ctx.Request, ctx.Client
    return nil
})
```

Disable defaults with `DisableDefaultTradeHooks: true` and supply hooks via `TradeCreatedHooks` in config.

Built-in validators:

- `ValidatePendingDataSignHook`
- `ValidateCreateTradeRequestHook`
- `ValidateCreateSummaryTxRequestHook`
- `ValidateRBFSidRequestHook`

Trade-creating OPS methods (`CreateTrade`, `DeployContract`, `CreateSummaryTx`, batch, stake, …) run hooks automatically after the response is received.

---

## Log watch

Incremental `lastID` polling; fixed page size and backoff inside the SDK.

```go
ctx, cancel := context.WithCancel(context.Background())
defer cancel()

go func() {
    _ = wc.WatchTradeLog(ctx, startLastID, func(row client.TradeLogResult) error {
        fmt.Println(row.ID, row.TxID, row.Success)
        return nil
    })
}()

// cancel() to stop the watch loop
```

`WatchBalanceLog` and `WatchMonitorAlert` follow the same pattern with `BalanceLogResult` and `MonitorAlertResult`.

---

## Layout

| Path | Purpose |
|------|---------|
| `types/` | OPS/MPC request/response DTOs (source definitions for Java/Rust generators) |
| `client/` | `WalletClient` SDK entry, TradeHook, LogWatch, TradeValidate |
| `integration/` | Live and offline integration tests |
| `config/` | Sample `ops.json` / `cli.json` |
| `docs/` | API reference and integration guides |
| `scripts/` | Doc and regression helpers |

---

## Tests

```bash
# Unit / offline
go test ./client/... -v

# Trade hook validators
go test ./client/... -run 'TestValidate|TestRunTradeCreatedHooks' -v

# Offline integration guards
go test ./integration -run 'TestNewSkipsEmptyConfig|TestCLIRequiresMPC|TestOPSRequiresOPS|TestTransferRequiresBothSessions|TestNilRequestRejected' -v

# Live OPS + MPC connect
go test ./integration -run TestConnectOPSAndCLI -v

# BTC wallet → account → address flow (writes to OPS MongoDB)
RUN_BTC_WALLET_ADDRESS=1 go test ./integration -run TestBTCWalletToAddress -v

# Live API suites
go test ./integration -run TestOPSAPI -v
go test ./integration -run TestCLIAPI -v

# End-to-end transfer (on-chain)
RUN_OPS_TRANSFER=1 go test ./integration -run TestTransfer -v

# Skip all live integration tests
go test ./integration -short
```

Copy `integration/integration.env.example` to `integration/integration.env`. Set `TEST_OPS_DOMAIN`, `TEST_MPC_DOMAIN`, and `TEST_APP_KEY` to override config during integration runs.

Optional env vars for the BTC address flow:

| Variable | Default | Purpose |
|----------|---------|---------|
| `TEST_OPS_SYMBOL` | `BTC` | Chain symbol |
| `TEST_BTC_ADDRESS_PREFIX` | `bcrt1` | Expected address prefix (regtest) |

Write gates (default off): `RUN_OPS_WRITE`, `RUN_OPS_SUBMIT`, `RUN_OPS_TRANSFER`.

---

## License

MIT
