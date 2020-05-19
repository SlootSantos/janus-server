package cdn

import (
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudfront"
)

const originAccessIDPrefix = "origin-access-identity/cloudfront/"

func (c *CDN) constructStandardDistroConfig(bucketID string, originAccessID string, stackID string) *cloudfront.CreateDistributionInput {
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

	aliases := &cloudfront.Aliases{
		Items: []*string{
			aws.String(stackID + "." + os.Getenv("DOMAIN_HOST")),
		},
		Quantity: aws.Int64(1),
	}

	certificate := &cloudfront.ViewerCertificate{
		ACMCertificateArn:      aws.String(os.Getenv("DOMAIN_CERT_ARN")),
		MinimumProtocolVersion: aws.String("TLSv1.2_2018"),
		SSLSupportMethod:       aws.String("sni-only"),
	}

	customErrorBehaviour := &cloudfront.CustomErrorResponses{
		Items: []*cloudfront.CustomErrorResponse{
			{
				ErrorCode:        aws.Int64(404),
				ResponseCode:     aws.String("200"),
				ResponsePagePath: aws.String("/index.html"),
			},
		},
		Quantity: aws.Int64(1),
	}

	config := &cloudfront.CreateDistributionInput{
		DistributionConfig: &cloudfront.DistributionConfig{
			CustomErrorResponses: customErrorBehaviour,
			Aliases:              aliases,
			ViewerCertificate:    certificate,
			DefaultRootObject:    aws.String("index.html"),
			CallerReference:      aws.String(bucketID),
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
