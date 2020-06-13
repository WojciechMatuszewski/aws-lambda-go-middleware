package middleware

import (
	"context"
)

type middlewareFunc func(ctx context.Context, payload []byte) ([]byte, error)

func (m middlewareFunc) Invoke(ctx context.Context, payload []byte) ([]byte, error) {
	return m(ctx, payload)
}
