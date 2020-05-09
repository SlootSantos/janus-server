package cdn

import (
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudfront"
)

const originAccessIDPrefix = "origin-access-identity/cloudfront/"

func (c *CDN) constructStandardDistroConfig(bucketID string, originAccessID string) *cloudfront.CreateDistributionInput {
	cacheBehavior := &cloudfront.DefaultCacheBehavior{
		TargetOriginId:       aws.String(cdnPrefix + bucketID),
		ViewerProtocolPolicy: aws.String("allow-all"),
		MinTTL:               aws.Int64(10),
		TrustedSigners: &cloudfront.TrustedSigners{
			Quantity: aws.Int64(0),
			Enabled:  aws.Bool(false),
		},
		ForwardedValues: &cloudfront.ForwardedValues{
			QueryString: aws.Bool(false),
			Cookies: &cloudfront.CookiePreference{
				Forward: aws.String("none"),
			},
		},
	}

	origins := []*cloudfront.Origin{
		{
			DomainName: aws.String(bucketID + cdnS3OriginSuffix),
			Id:         aws.String(cdnPrefix + bucketID),
			S3OriginConfig: &cloudfront.S3OriginConfig{
				OriginAccessIdentity: aws.String("origin-access-identity/cloudfront/" + originAccessID),
			},
		},
	}

	config := &cloudfront.CreateDistributionInput{
		DistributionConfig: &cloudfront.DistributionConfig{
			CallerReference:      aws.String(time.Now().String() + bucketID),
			Comment:              aws.String(bucketID),
			Enabled:              aws.Bool(true),
			DefaultCacheBehavior: cacheBehavior,
			Origins: &cloudfront.Origins{
				Quantity: aws.Int64(1),
				Items:    origins,
			},
		},
	}

	return config
}
