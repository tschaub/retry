package retry_test

import (
	"context"
	"errors"
	"fmt"

	"github.com/tschaub/retry"
)

func Example_basic() {

	// assume you have a function that fails twice and then succeeds
	calls := 0
	failTwice := func() error {
		calls += 1
		if calls < 2 {
			return errors.New("failed")
		}
		return nil
	}

	// try calling the function at most 5 times
	retries := 5
	ctx := context.Background()
	err := retry.Limit(ctx, retries, func(context.Context, int) error {
		return failTwice()
	})

	fmt.Printf("called: %d\n", calls)
	fmt.Printf("error: %v\n", err)

	// Output:
	// called: 2
	// error: <nil>
}

func Example_stop() {

	// In some cases (like http requests), you may want
	// to stop retrying early under certain conditions
	// (like a 4xx response).  The retry.Stop(err) function
	// can be used for this purpose.

	retries := 500
	ctx := context.Background()
	err := retry.Limit(ctx, retries, func(ctx context.Context, call int) error {
		// here we might be checking for a response code
		if call > 3 {
			return retry.Stop(errors.New("stopped"))
		}
		return errors.New("failed")
	})

	fmt.Printf("error: %v\n", err)

	// Output:
	// error: stopped
}
