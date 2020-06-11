package cdn

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudfront"
)

type constructDistroConfigInput struct {
	originAccessID string
	subdomain      string
	certificateARN string
	bucketID       string
	stackID        string
	isThirdParty   bool
}

const (
	originAccessIDPrefix      = "origin-access-identity/cloudfront/"
	defaultRootObject         = "index.html"
	cacheProtocolPolicy       = "redirect-to-https"
	cacheForwardHeaderHost    = "Host"
	certificateSSLMethods     = "sni-only"
	certifcateProtocolVersion = "TLSv1.2_2018"
	errorRequestCode          = 404
	errorResponseCode         = "200"
)

func (c *CDN) constructStandardDistroConfig(input *constructDistroConfigInput) *cloudfront.CreateDistributionInput {
	config := &cloudfront.CreateDistributionInput{
		DistributionConfig: &cloudfront.DistributionConfig{
			DefaultRootObject: aws.String(defaultRootObject),
			CallerReference:   aws.String(input.bucketID),
			Comment:           aws.String(input.bucketID),
			Enabled:           aws.Bool(true),
			// Aliases:           constructAliases(input.subdomain),
			Origins: constructS3Origins(input.bucketID, input.originAccessID),
			// ViewerCertificate:    constructCertificate(input.certificateARN),
			DefaultCacheBehavior: constructDefaultCacheBehavior(input.bucketID, input.isThirdParty),
			CustomErrorResponses: constructErrorBehavior(),
		},
	}

	return config
}

func constructCertificate(certARN string) *cloudfront.ViewerCertificate {
	return &cloudfront.ViewerCertificate{
		SSLSupportMethod:       aws.String(certificateSSLMethods),
		ACMCertificateArn:      aws.String(certARN),
		MinimumProtocolVersion: aws.String(certifcateProtocolVersion),
	}
}

func constructErrorBehavior() *cloudfront.CustomErrorResponses {
	return &cloudfront.CustomErrorResponses{
		Items: []*cloudfront.CustomErrorResponse{
			{
				ErrorCode:        aws.Int64(errorRequestCode),
				ResponseCode:     aws.String(errorResponseCode),
				ResponsePagePath: aws.String("/" + defaultRootObject),
			},
		},
		Quantity: aws.Int64(1),
	}
}

func constructAliases(subdomain string, domainHost string) *cloudfront.Aliases {
	alias := subdomain + "." + domainHost

	return &cloudfront.Aliases{
		Quantity: aws.Int64(5),
		Items: []*string{
			aws.String(alias),
			aws.String(aliasPrefixGreenDeployment + alias),
			aws.String(aliasPrefixDevelopmentEnv + alias),
			aws.String(aliasPrefixStageEnv + alias),
			aws.String(aliasPrefixPRPreview + alias),
		},
	}
}

func constructS3Origins(bucketID string, originAccessID string) *cloudfront.Origins {
	return &cloudfront.Origins{
		Quantity: aws.Int64(1),
		Items: []*cloudfront.Origin{
			{
				Id:         aws.String(cdnPrefix + bucketID),
				DomainName: aws.String(bucketID + cdnS3OriginSuffix),
				S3OriginConfig: &cloudfront.S3OriginConfig{
					OriginAccessIdentity: aws.String(originAccessIDPrefix + originAccessID),
				},
			},
		},
	}
}

func constructDefaultCacheBehavior(bucketID string, isThirdParty bool) *cloudfront.DefaultCacheBehavior {
	return &cloudfront.DefaultCacheBehavior{
		MinTTL:                     aws.Int64(10),
		Compress:                   aws.Bool(true),
		TargetOriginId:             aws.String(cdnPrefix + bucketID),
		ViewerProtocolPolicy:       aws.String(cacheProtocolPolicy),
		LambdaFunctionAssociations: blueGreenLambdaFuncConfig(isThirdParty),
		TrustedSigners: &cloudfront.TrustedSigners{
			Quantity: aws.Int64(0),
			Enabled:  aws.Bool(false),
		},
		ForwardedValues: &cloudfront.ForwardedValues{
			QueryString: aws.Bool(false),
			Cookies: &cloudfront.CookiePreference{
				Forward: aws.String("none"),
			},
			Headers: &cloudfront.Headers{
				Items: []*string{
					aws.String(cacheForwardHeaderHost),
				},
				Quantity: aws.Int64(1),
			},
		},
	}
}
