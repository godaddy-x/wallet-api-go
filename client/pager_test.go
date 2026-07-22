package client

import "testing"

func TestNormalizePageLimit(t *testing.T) {
	tests := []struct {
		in   int64
		want int64
	}{
		{0, DefaultPageLimit},
		{-1, DefaultPageLimit},
		{50, 50},
		{200, 200},
		{500, MaxPageLimit},
		{5000, MaxPageLimit},
	}
	for _, tc := range tests {
		if got := normalizePageLimit(tc.in); got != tc.want {
			t.Fatalf("normalizePageLimit(%d): got %d want %d", tc.in, got, tc.want)
		}
	}
}

func TestClampPageLimitPaginatedPath(t *testing.T) {
	req := &FindTradeLogReq{}
	req.Limit = 5000
	clampPageLimit("/api/FindTradeLog", req)
	if req.Limit != MaxPageLimit {
		t.Fatalf("FindTradeLog limit: got %d want %d", req.Limit, MaxPageLimit)
	}

	req.Limit = 0
	clampPageLimit("/api/FindTradeLog", req)
	if req.Limit != DefaultPageLimit {
		t.Fatalf("FindTradeLog default limit: got %d want %d", req.Limit, DefaultPageLimit)
	}

	newlyReq := &FindTradeNewlyReq{}
	newlyReq.Limit = 5000
	clampPageLimit("/api/FindTradeNewly", newlyReq)
	if newlyReq.Limit != MaxPageLimit {
		t.Fatalf("FindTradeNewly limit: got %d want %d", newlyReq.Limit, MaxPageLimit)
	}
}

func TestClampPageLimitNonPaginatedPath(t *testing.T) {
	req := &CreateWalletReq{}
	req.Limit = 5000
	clampPageLimit("/api/CreateWallet", req)
	if req.Limit != 5000 {
		t.Fatalf("CreateWallet limit should be unchanged: got %d", req.Limit)
	}
}

func TestClampPageLimitMPCPath(t *testing.T) {
	req := &CliFindWalletListReq{}
	req.Limit = 999
	clampPageLimit("/api/FindWalletList", req)
	if req.Limit != MaxPageLimit {
		t.Fatalf("FindWalletList limit: got %d want %d", req.Limit, MaxPageLimit)
	}
}
