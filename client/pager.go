package client

import (
	"reflect"

	"github.com/godaddy-x/freego/node/common"
)

const (
	// DefaultPageLimit is the SDK page size when limit is omitted on list APIs.
	DefaultPageLimit int64 = 200
	// MaxPageLimit is the SDK-side cap for list API page size (server may allow more).
	MaxPageLimit int64 = 200
)

var paginatedPaths = map[string]struct{}{
	"/api/FindAccountByWalletID": {},
	"/api/GetAccountBalanceList": {},
	"/api/FindAddressByAccountID": {},
	"/api/GetAddressBalanceList":  {},
	"/api/GetContracts":           {},
	"/api/GetContractTemplates":   {},
	"/api/FindTradeLog":           {},
	"/api/FindBalanceLog":         {},
	"/api/FindMonitorAlert":       {},
	"/api/FindWalletList":         {},
}

func normalizePageLimit(limit int64) int64 {
	if limit <= 0 {
		return DefaultPageLimit
	}
	if limit > MaxPageLimit {
		return MaxPageLimit
	}
	return limit
}

func clampPageLimit(path string, req any) {
	if _, ok := paginatedPaths[path]; !ok {
		return
	}
	base := baseReqPtr(req)
	if base == nil {
		return
	}
	base.Limit = normalizePageLimit(base.Limit)
}

func baseReqPtr(req any) *common.BaseReq {
	if req == nil {
		return nil
	}
	v := reflect.ValueOf(req)
	if v.Kind() != reflect.Ptr || v.IsNil() {
		return nil
	}
	f := v.Elem().FieldByName("BaseReq")
	if !f.IsValid() || !f.CanAddr() {
		return nil
	}
	base, ok := f.Addr().Interface().(*common.BaseReq)
	if !ok {
		return nil
	}
	return base
}
