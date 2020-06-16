package cdn

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudfront"
)

const lambdaEventOriginReq = "origin-request"

func blueGreenLambdaFuncConfig(lambdaARN string) *cloudfront.LambdaFunctionAssociations {
	return &cloudfront.LambdaFunctionAssociations{
		Items: []*cloudfront.LambdaFunctionAssociation{
			{
				LambdaFunctionARN: aws.String(lambdaARN),
				IncludeBody:       aws.Bool(false),
				EventType:         aws.String(lambdaEventOriginReq),
			},
		},
		Quantity: aws.Int64(1),
	}
}
