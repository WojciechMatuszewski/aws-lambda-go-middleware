package middleware

import (
	"context"
	"time"

	"lambda-middleware/internal/logger"
)

func WithTimeoutLogger() middlewareFunc {
	return func(ctx context.Context, payload []byte) (context.Context, []byte, error) {
		go func() {
			deadline, ok := ctx.Deadline()
			if !ok {
				panic("deadline not found")
			}

			deadline = deadline.Add(-50 * time.Millisecond)
			timeoutChan := time.After(time.Until(deadline))

			select {
			case <-timeoutChan:
				logger.Log.Error().Msg("timeout")
			case <-ctx.Done():
				return
			}
		}()

		return ctx, payload, nil
	}
}
