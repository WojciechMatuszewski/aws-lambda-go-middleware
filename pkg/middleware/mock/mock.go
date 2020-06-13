// //go:generate mockgen -package=mock -destination=./ssm.go github.com/aws/aws-sdk-go-v2/service/ssm/ssmiface ClientAPI
//go:generate mockgen -package=mock -destination=./ssm-parameter.go -source=../ssm-parameter.go
package mock
