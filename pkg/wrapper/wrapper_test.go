package wrapper_test

import (
	"context"
	"fmt"
	"testing"

	"lambda-middleware/pkg/wrapper"

	"github.com/aws/aws-lambda-go/lambda"

	"github.com/aws/aws-lambda-go/events"

	"github.com/stretchr/testify/assert"
)

func TestWrapper(t *testing.T) {
	t.Run("works", func(t *testing.T) {
		firstCalled := false
		first := wrapper.MiddlewareFunc(func(ctx context.Context, payload []byte) ([]byte, error) {
			firstCalled = true
			return nil, nil
		})

		secondCalled := false
		second := wrapper.MiddlewareFunc(func(ctx context.Context, payload []byte) ([]byte, error) {
			secondCalled = true
			return nil, nil
		})

		handlerCalled := false
		handler := func() {
			handlerCalled = true
		}

		wrapper.Wrap(handler).Next(first).Next(second).Invoke(context.TODO(), nil)
		assert.True(t, firstCalled)
		assert.True(t, secondCalled)
		assert.True(t, handlerCalled)

	})

	t.Run("works2", func(t *testing.T) {
		handler := func(context context.Context, request events.APIGatewayV2HTTPRequest) {
			fmt.Println(request.Body)
		}

		lambda.StartHandler(wrapper.Wrap(handler))

		lambda.Start()
	})
}
