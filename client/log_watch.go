package client

import (
	"context"
	"errors"
	"time"
)

// Fixed incremental pull tuning — not exposed to callers to avoid OPS overload.
const incrementalPageSize = DefaultPageLimit

const incrementalCatchUpWait = 5 * time.Second

var incrementalEmptyBackoff = []time.Duration{
	time.Second,
	2 * time.Second,
	3 * time.Second,
	5 * time.Second,
}

// TradeLogHandler is invoked for each new trade log row. Return a non-nil error to stop watching.
type TradeLogHandler func(entry TradeLogResult) error

// BalanceLogHandler is invoked for each new balance log row. Return a non-nil error to stop watching.
type BalanceLogHandler func(entry BalanceLogResult) error

// MonitorAlertHandler is invoked for each new monitor alert row. Return a non-nil error to stop watching.
type MonitorAlertHandler func(entry MonitorAlertResult) error

// WatchTradeLog incrementally pulls trade logs using lastID pagination and invokes fn for each row.
// Page size (200) and poll/backoff intervals are fixed by the SDK.
// startLastID is the persisted watermark; 0 starts from the oldest available row.
// It blocks until ctx is cancelled, fn returns an error, or an OPS request fails.
func (c *WalletClient) WatchTradeLog(ctx context.Context, startLastID int64, fn TradeLogHandler) error {
	if fn == nil {
		return errors.New("trade log handler is nil")
	}
	return watchIncremental(ctx, startLastID, incrementalPageSize, incrementalCatchUpWait, incrementalEmptyBackoff, func(cursor int64) (int, int64, error) {
		req := &FindTradeLogReq{}
		req.LastID = cursor
		req.Limit = incrementalPageSize
		res, err := c.FindTradeLog(req)
		if err != nil {
			return 0, cursor, err
		}
		next := incrementalNextCursor(cursor, res.Limit.LastID, len(res.Result), func(i int) int64 {
			return res.Result[i].ID
		})
		for _, row := range res.Result {
			if err := fn(row); err != nil {
				return 0, next, err
			}
		}
		return len(res.Result), next, nil
	})
}

// WatchBalanceLog incrementally pulls balance logs using lastID pagination and invokes fn for each row.
func (c *WalletClient) WatchBalanceLog(ctx context.Context, startLastID int64, fn BalanceLogHandler) error {
	if fn == nil {
		return errors.New("balance log handler is nil")
	}
	return watchIncremental(ctx, startLastID, incrementalPageSize, incrementalCatchUpWait, incrementalEmptyBackoff, func(cursor int64) (int, int64, error) {
		req := &FindBalanceLogReq{}
		req.LastID = cursor
		req.Limit = incrementalPageSize
		res, err := c.FindBalanceLog(req)
		if err != nil {
			return 0, cursor, err
		}
		next := incrementalNextCursor(cursor, res.Limit.LastID, len(res.Result), func(i int) int64 {
			return res.Result[i].ID
		})
		for _, row := range res.Result {
			if err := fn(row); err != nil {
				return 0, next, err
			}
		}
		return len(res.Result), next, nil
	})
}

// WatchMonitorAlert incrementally pulls monitor alerts using lastID pagination and invokes fn for each row.
func (c *WalletClient) WatchMonitorAlert(ctx context.Context, startLastID int64, fn MonitorAlertHandler) error {
	if fn == nil {
		return errors.New("monitor alert handler is nil")
	}
	return watchIncremental(ctx, startLastID, incrementalPageSize, incrementalCatchUpWait, incrementalEmptyBackoff, func(cursor int64) (int, int64, error) {
		req := &FindMonitorAlertReq{}
		req.LastID = cursor
		req.Limit = incrementalPageSize
		res, err := c.FindMonitorAlert(req)
		if err != nil {
			return 0, cursor, err
		}
		next := incrementalNextCursor(cursor, res.Limit.LastID, len(res.Result), func(i int) int64 {
			return res.Result[i].ID
		})
		for _, row := range res.Result {
			if err := fn(row); err != nil {
				return 0, next, err
			}
		}
		return len(res.Result), next, nil
	})
}

type incrementalFetch func(cursor int64) (count int, nextCursor int64, err error)

func watchIncremental(
	ctx context.Context,
	startLastID int64,
	pageSize int64,
	catchUpWait time.Duration,
	emptyBackoff []time.Duration,
	fetch incrementalFetch,
) error {
	cursor := startLastID
	emptyStreak := 0

	for {
		if err := ctx.Err(); err != nil {
			return err
		}

		count, nextCursor, err := fetch(cursor)
		if err != nil {
			return err
		}
		cursor = nextCursor

		if count == 0 {
			if err := sleepWithContext(ctx, emptyBackoffWait(emptyBackoff, emptyStreak)); err != nil {
				return err
			}
			emptyStreak++
			continue
		}

		emptyStreak = 0
		if int64(count) >= pageSize {
			continue
		}
		if err := sleepWithContext(ctx, catchUpWait); err != nil {
			return err
		}
	}
}

func emptyBackoffWait(steps []time.Duration, streak int) time.Duration {
	if len(steps) == 0 {
		return incrementalEmptyBackoff[0]
	}
	if streak >= len(steps) {
		return steps[len(steps)-1]
	}
	return steps[streak]
}

func incrementalNextCursor(cursor, limitLastID int64, count int, idAt func(int) int64) int64 {
	if count == 0 {
		return cursor
	}
	if limitLastID > cursor {
		return limitLastID
	}
	last := idAt(count - 1)
	if last > cursor {
		return last
	}
	return cursor
}

func sleepWithContext(ctx context.Context, d time.Duration) error {
	if d <= 0 {
		return nil
	}
	timer := time.NewTimer(d)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}
