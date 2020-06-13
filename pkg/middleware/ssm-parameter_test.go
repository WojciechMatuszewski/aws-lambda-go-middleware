package middleware_test

import (
	"context"
	"errors"
	"testing"

	"lambda-middleware/pkg/middleware"
	"lambda-middleware/pkg/middleware/mock"

	"github.com/stretchr/testify/assert"

	"github.com/aws/aws-lambda-go/lambda"

	"github.com/golang/mock/gomock"
)

func TestWithSSMParameter(t *testing.T) {
	ctx := context.Background()
	t.Run("success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		getter := mock.NewMockParameterGetter(ctrl)
		getter.EXPECT().Get(ctx, "parameterName").Return("parameterValue", nil)

		rootCalled := false
		rootHandler := lambda.NewHandler(func(ctx context.Context) error {
			rootCalled = true
			v, found := ctx.Value("contextKey").(string)

			assert.True(t, found)
			assert.Equal(t, "parameterValue", v)
			return nil
		})

		handler := middleware.WithSSMParameter(getter, "parameterName", "contextKey")(rootHandler)
		_, err := handler.Invoke(ctx, nil)
		assert.NoError(t, err)
		assert.True(t, rootCalled)
	})

	t.Run("getter failure", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		getter := mock.NewMockParameterGetter(ctrl)
		getter.EXPECT().Get(ctx, "parameterName").Return("", errors.New("boom"))

		rootCalled := false
		rootHandler := lambda.NewHandler(func() {
			rootCalled = true
		})

		handler := middleware.WithSSMParameter(getter, "parameterName", "contextKey")(rootHandler)
		_, err := handler.Invoke(ctx, nil)
		assert.Error(t, err)
		assert.False(t, rootCalled)
		assert.Equal(t, errors.New("boom"), err)
	})

}
