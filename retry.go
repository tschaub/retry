package retry

import (
	"context"
)

type errStop struct {
	err error
}

func (e errStop) Error() string {
	return "retry stopped"
}

// Stop returns an error that will stop retries.
func Stop(err error) error {
	return errStop{err: err}
}

// Func is a function that can be retried.  Called with the number of previous attempts.
type Func func(context.Context, int) error

// Limit will retry a function until it does not error, the limit is reached, or the context is cancelled.
// If the limit is reached, the last error will be returned.  To stop retries early, return Stop(err).
func Limit(ctx context.Context, limit int, fn Func) error {
	var err error
	for attempt := 0; attempt < limit; attempt++ {
		if err := ctx.Err(); err != nil {
			return err
		}

		err = fn(ctx, attempt)
		if err == nil {
			return nil
		}
		if stopped, ok := err.(errStop); ok {
			return stopped.err
		}
	}

	return err
}
