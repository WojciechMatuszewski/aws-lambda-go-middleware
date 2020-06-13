package middleware

import (
	"context"
	"time"

	"lambda-middleware/internal/logger"

	"github.com/aws/aws-lambda-go/lambda"
)

func WithTimeoutLogger(next lambda.Handler) lambda.Handler {
	return middlewareFunc(func(ctx context.Context, payload []byte) ([]byte, error) {
		go func() {
			deadline, ok := ctx.Deadline()
			if !ok {
				panic("not ok")
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

		return next.Invoke(ctx, payload)
	})
}
