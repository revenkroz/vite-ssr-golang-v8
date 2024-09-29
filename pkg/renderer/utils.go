package renderer

import (
	"errors"
	"fmt"
	"rogchap.com/v8go"
)

func resolvePromise(ctx *v8go.Context, val *v8go.Value, err error) (*v8go.Value, error) {
	if err != nil || !val.IsPromise() {
		return val, err
	}
	for {
		switch p, _ := val.AsPromise(); p.State() {
		case v8go.Fulfilled:
			return p.Result(), nil
		case v8go.Rejected:
			return nil, errors.New(p.Result().DetailString())
		case v8go.Pending:
			ctx.PerformMicrotaskCheckpoint() // run VM to make progress on the promise
			// go round the loop again...
		default:
			return nil, fmt.Errorf("illegal v8go.Promise state %d", p) // unreachable
		}
	}
}

func formatError(err error) error {
	var jsErr *v8go.JSError
	if errors.As(err, &jsErr) {
		err = fmt.Errorf("%v", jsErr.StackTrace)
	}

	return err
}
