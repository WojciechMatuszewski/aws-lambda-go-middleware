package parameter

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go-v2/service/ssm/ssmiface"
)

type Service struct {
	ssm ssmiface.ClientAPI
}

func NewService(ssm ssmiface.ClientAPI) Service {
	return Service{ssm}
}

func (s Service) Get(ctx context.Context, name string) (string, error) {
	req := s.ssm.GetParameterRequest(&ssm.GetParameterInput{
		Name: aws.String(name),
	})

	resp, err := req.Send(ctx)
	if err != nil {
		return "", err
	}

	return *resp.Parameter.Value, nil
}
