package cdn

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudfront"
)

const lambdaEventOriginReq = "origin-request"
const lambdaFuncARN = "arn:aws:lambda:us-east-1:976589619057:function:janus-exmaple-redirect:34" // => env?

func blueGreenLambdaFuncConfig() *cloudfront.LambdaFunctionAssociations {
	return &cloudfront.LambdaFunctionAssociations{
		Items: []*cloudfront.LambdaFunctionAssociation{
			{
				LambdaFunctionARN: aws.String(lambdaFuncARN),
				IncludeBody:       aws.Bool(false),
				EventType:         aws.String(lambdaEventOriginReq),
			},
		},
		Quantity: aws.Int64(1),
	}
}
