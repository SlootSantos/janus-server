package cdn

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudfront"
)

const lambdaEventOriginReq = "origin-request"

const lambdaFuncARN = "arn:aws:lambda:us-east-1:976589619057:function:janus-exmaple-redirect:34"              // => env?
const lambdaFuncARNThirdParty = "arn:aws:lambda:us-east-1:108151951856:function:stackers-handler-routing-2:1" // => env?

func blueGreenLambdaFuncConfig(isThirdParty bool) *cloudfront.LambdaFunctionAssociations {
	lambdaARN := lambdaFuncARN

	if isThirdParty {
		lambdaARN = lambdaFuncARNThirdParty
	}

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
