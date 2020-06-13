package middleware_test

import (
	"bytes"
	"context"
	"log"
	"os"
	"testing"
	"time"

	"lambda-middleware/internal/logger"
	"lambda-middleware/pkg/middleware"

	"github.com/aws/aws-lambda-go/lambda"
)

func TestWithTimeoutLogger(t *testing.T) {

	captureLog := func(fn func()) string {
		defer log.SetOutput(os.Stdout)
		var buf bytes.Buffer
		logger.Out(&buf)
		fn()
		return buf.String()
	}

	t.Run("no timeout", func(t *testing.T) {
		buf := captureLog(func() {
			timestamp := time.Now().Add(70 * time.Millisecond)
			ctx, cancel := context.WithDeadline(context.Background(), timestamp)
			defer cancel()

			h := lambda.NewHandler(func() {
				time.Sleep(10 * time.Millisecond)
			})
			_, _ = middleware.WithTimeoutLogger(h).Invoke(ctx, nil)

		})

		if buf != "" {
			t.Fatalf("expected %v, got %v", "", buf)
		}
	})

	t.Run("timeout logged", func(t *testing.T) {
		buf := captureLog(func() {
			timestamp := time.Now().Add(70 * time.Millisecond)
			ctx, cancel := context.WithDeadline(context.Background(), timestamp)
			defer cancel()

			h := lambda.NewHandler(func() {
				time.Sleep(100 * time.Millisecond)
			})
			_, _ = middleware.WithTimeoutLogger(h).Invoke(ctx, nil)
		})

		if buf == "" {
			t.Fatalf("expected log output, got none")
		}
	})
}
