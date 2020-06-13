package middleware

import (
	"context"

	"github.com/aws/aws-lambda-go/lambda"
)

type ParameterGetter interface {
	Get(ctx context.Context, name string) (string, error)
}

func WithSSMParameter(getter ParameterGetter, name string, ctxKey string) func(next lambda.Handler) lambda.Handler {
	return func(next lambda.Handler) lambda.Handler {
		return middlewareFunc(func(ctx context.Context, payload []byte) ([]byte, error) {
			p, err := getter.Get(ctx, name)
			if err != nil {
				return nil, err
			}

			ctx = context.WithValue(ctx, ctxKey, p)
			return next.Invoke(ctx, payload)
		})
	}
}
