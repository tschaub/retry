package retry_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tschaub/retry"
)

var fail = errors.New("fail")

func failForever() error {
	return fail
}

func failTwice(ctx context.Context, call int) error {
	if call < 2 {
		return failForever()
	}
	return nil
}

func TestLimitWithPermanentFailure(t *testing.T) {
	retries := 3
	ctx := context.Background()

	calls := 0
	err := retry.Limit(ctx, retries, func(context.Context, int) error {
		calls++
		return failForever()
	})
	assert.Equal(t, fail, err)
	assert.Equal(t, retries, calls)
}

func TestLimitWithACoupleFailures(t *testing.T) {
	retries := 5
	ctx := context.Background()

	calls := 0
	err := retry.Limit(ctx, retries, func(ctx context.Context, call int) error {
		calls++
		return failTwice(ctx, call)
	})
	assert.Nil(t, err)
	assert.Equal(t, 3, calls)
}

func TestLimitFinishEarly(t *testing.T) {
	retries := 100
	ctx := context.Background()

	early := 10

	calls := 0
	err := retry.Limit(ctx, retries, func(ctx context.Context, call int) error {
		calls++
		if call < early {
			return failForever()
		}
		return nil
	})

	assert.Nil(t, err)
	assert.Equal(t, early+1, calls)
}

func TestLimitStopped(t *testing.T) {
	retries := 100
	ctx := context.Background()

	early := 10

	tired := errors.New("tired")
	calls := 0
	err := retry.Limit(ctx, retries, func(ctx context.Context, call int) error {
		calls++
		if call < early {
			return failForever()
		}
		return retry.Stop(tired)
	})

	assert.Equal(t, err, tired)
	assert.Equal(t, early+1, calls)
}

func TestLimitWithContextCancel(t *testing.T) {
	retries := 100
	ctx, cancel := context.WithCancel(context.Background())

	early := 10

	calls := 0
	err := retry.Limit(ctx, retries, func(ctx context.Context, call int) error {
		calls++
		if call == early {
			cancel()
		}
		return failForever()
	})

	assert.Equal(t, err, context.Canceled)
	assert.Equal(t, early+1, calls)
}
