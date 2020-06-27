package middleware

import (
	"context"
)

type ParameterGetter interface {
	Get(ctx context.Context, name string) (string, error)
}

func WithSSMParameter(getter ParameterGetter, name string, ctxKey string) middlewareFunc {
	return func(ctx context.Context, payload []byte) (context.Context, []byte, error) {
		p, err := getter.Get(ctx, name)
		if err != nil {
			return ctx, nil, err
		}

		ctx = context.WithValue(ctx, ctxKey, p)
		return ctx, payload, nil
	}
}
