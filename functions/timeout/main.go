package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"lambda-middleware/internal/logger"
	"lambda-middleware/internal/parameter"
	"lambda-middleware/pkg/middleware"

	"github.com/aws/aws-sdk-go-v2/service/ssm"

	"github.com/aws/aws-sdk-go-v2/aws/external"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type Payload struct {
	Timeout bool `json:"timeout"`
}

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	rBody := request.Body
	ssmP, found := ctx.Value("ssmparam").(string)
	if !found {
		return events.APIGatewayProxyResponse{StatusCode: http.StatusBadRequest, Body: "ssm parameter not found"}, nil
	}

	logger.Log.Info().Str("parameter", ssmP).Msg("within ssm")

	var payload Payload
	err := json.Unmarshal([]byte(rBody), &payload)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: http.StatusInternalServerError, Body: "marshal problem"}, nil
	}

	if payload.Timeout {
		fmt.Println("sleeping")
		time.Sleep(2 * time.Second)
	}

	fmt.Println("returning")
	return events.APIGatewayProxyResponse{StatusCode: http.StatusOK, Body: http.StatusText(http.StatusOK)}, nil
}

func main() {
	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		panic(err.Error())
	}

	ssmAPI := ssm.New(cfg)
	paramService := parameter.NewService(ssmAPI)

	root := lambda.NewHandler(handler)
	withTimeout := middleware.WithTimeoutLogger(root)
	h := middleware.WithSSMParameter(paramService, os.Getenv("PARAMETER_NAME"), "ssmparam")(withTimeout)

	lambda.StartHandler(h)
}
