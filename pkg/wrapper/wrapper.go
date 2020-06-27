package wrapper

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/lambda"
)

// Wrap(handler).Next(middleware).Next(middleware)

type Wrapper struct {
	current MiddlewareFunc
	handler lambda.Handler
}

func Wrap(handlerFunc interface{}) *Wrapper {
	lambdaHandler := lambda.NewHandler(handlerFunc)
	return &Wrapper{
		handler: lambdaHandler,
		current: nil,
	}
}

func (w *Wrapper) Invoke(ctx context.Context, payload []byte) ([]byte, error) {
	if w.current == nil {
		return w.handler.Invoke(ctx, payload)
	}

	out, err := w.current.Invoke(ctx, payload)
	if err != nil {
		return nil, err
	}

	return w.handler.Invoke(ctx, out)
}

func (w *Wrapper) Next(next MiddlewareFunc) *Wrapper {
	return &Wrapper{current: func(ctx context.Context, payload []byte) ([]byte, error) {
		if w.current == nil {
			return next.Invoke(ctx, payload)
		}

		out, err := w.current.Invoke(ctx, payload)
		if err != nil {
			panic(err.Error())
		}

		return next.Invoke(ctx, out)
	}}
}

type MiddlewareFunc func(ctx context.Context, payload []byte) ([]byte, error)

func (m MiddlewareFunc) Invoke(ctx context.Context, payload []byte) ([]byte, error) {
	return m(ctx, payload)
}

var Test1 = MiddlewareFunc(func(ctx context.Context, payload []byte) ([]byte, error) {
	fmt.Println("test1")
	return payload, nil
})

var Test2 = MiddlewareFunc(func(ctx context.Context, payload []byte) ([]byte, error) {
	fmt.Println("test2")
	return payload, nil
})
