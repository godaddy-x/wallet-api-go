package client

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestWatchIncrementalDrainsBacklogThenWaits(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(context.Background(), 25*time.Millisecond)
	defer cancel()

	pages := [][]int{
		{1, 2, 3},
		{4, 5},
	}
	page := 0
	attempts := 0

	err := watchIncremental(ctx, 0, 3, 5*time.Millisecond, []time.Duration{2 * time.Millisecond, 4 * time.Millisecond}, func(cursor int64) (int, int64, error) {
		attempts++
		if page >= len(pages) {
			return 0, cursor, nil
		}
		rows := pages[page]
		page++
		next := int64(rows[len(rows)-1])
		return len(rows), next, nil
	})
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("expected deadline error, got %v attempts=%d page=%d", err, attempts, page)
	}
	if page != 2 {
		t.Fatalf("expected both backlog pages drained before wait, page=%d", page)
	}
	if attempts < 3 {
		t.Fatalf("expected immediate second page then idle polls, attempts=%d", attempts)
	}
}

func TestWatchIncrementalEmptyBackoffResetsAfterData(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(context.Background(), 40*time.Millisecond)
	defer cancel()

	attempts := 0
	err := watchIncremental(ctx, 0, 10, 2*time.Millisecond, []time.Duration{2 * time.Millisecond, 4 * time.Millisecond}, func(cursor int64) (int, int64, error) {
		attempts++
		switch attempts {
		case 1, 2:
			return 0, cursor, nil
		case 3:
			return 1, 100, nil
		default:
			return 0, 100, nil
		}
	})
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("expected deadline error, got %v attempts=%d", err, attempts)
	}
	if attempts < 4 {
		t.Fatalf("expected at least 4 fetch attempts, got %d", attempts)
	}
}

func TestWatchIncrementalHandlerErrorStops(t *testing.T) {
	t.Parallel()

	want := errors.New("stop")
	err := watchIncremental(context.Background(), 0, 10, time.Millisecond, []time.Duration{time.Millisecond}, func(cursor int64) (int, int64, error) {
		return 0, cursor, want
	})
	if !errors.Is(err, want) {
		t.Fatalf("expected %v, got %v", want, err)
	}
}

func TestIncrementalNextCursor(t *testing.T) {
	t.Parallel()

	got := incrementalNextCursor(10, 25, 2, func(i int) int64 {
		if i == 0 {
			return 20
		}
		return 25
	})
	if got != 25 {
		t.Fatalf("limit lastID: got %d want 25", got)
	}

	got = incrementalNextCursor(10, 0, 1, func(int) int64 { return 15 })
	if got != 15 {
		t.Fatalf("row id fallback: got %d want 15", got)
	}

	got = incrementalNextCursor(10, 0, 0, func(int) int64 { return 0 })
	if got != 10 {
		t.Fatalf("empty page keeps cursor: got %d want 10", got)
	}
}

func TestEmptyBackoffWait(t *testing.T) {
	t.Parallel()

	steps := incrementalEmptyBackoff
	cases := []struct {
		streak int
		want   time.Duration
	}{
		{0, time.Second},
		{1, 2 * time.Second},
		{2, 3 * time.Second},
		{3, 5 * time.Second},
		{99, 5 * time.Second},
	}
	for _, tc := range cases {
		got := emptyBackoffWait(steps, tc.streak)
		if got != tc.want {
			t.Fatalf("streak %d: got %v want %v", tc.streak, got, tc.want)
		}
	}
}

func TestIncrementalPullConstants(t *testing.T) {
	t.Parallel()

	if incrementalPageSize != DefaultPageLimit {
		t.Fatalf("page size: got %d want %d", incrementalPageSize, DefaultPageLimit)
	}
	if incrementalCatchUpWait != 5*time.Second {
		t.Fatalf("catch-up wait: got %v want 5s", incrementalCatchUpWait)
	}
	if len(incrementalEmptyBackoff) != 4 {
		t.Fatalf("backoff steps: got %d want 4", len(incrementalEmptyBackoff))
	}
}
