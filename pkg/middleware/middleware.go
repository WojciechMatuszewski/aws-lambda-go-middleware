package middleware

import (
	"context"

	"github.com/aws/aws-lambda-go/lambda"
)

type Wrapper struct {
	handler    lambda.Handler
	middleware middlewareFunc
}

func Wrap(root interface{}) Wrapper {
	return Wrapper{
		handler:    lambda.NewHandler(root),
		middleware: nil,
	}
}

func (w Wrapper) Invoke(ctx context.Context, payload []byte) ([]byte, error) {
	if w.middleware == nil {
		return w.handler.Invoke(ctx, payload)
	}

	newCtx, out, err := w.middleware.Invoke(ctx, payload)
	if err != nil {
		panic(err.Error())
	}
	return w.handler.Invoke(newCtx, out)
}

func (w Wrapper) Use(next middlewareFunc) Wrapper {
	return Wrapper{handler: w.handler, middleware: func(ctx context.Context, payload []byte) (context.Context, []byte, error) {
		if w.middleware == nil {
			return next.Invoke(ctx, payload)
		}

		newCtx, out, err := w.middleware.Invoke(ctx, payload)
		if err != nil {
			panic(err.Error())
		}

		return next.Invoke(newCtx, out)
	}}
}

type middlewareFunc func(ctx context.Context, payload []byte) (context.Context, []byte, error)

func (m middlewareFunc) Invoke(ctx context.Context, payload []byte) (context.Context, []byte, error) {
	return m(ctx, payload)
}
