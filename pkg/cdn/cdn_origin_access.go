package cdn

import (
	"context"
	"log"
	"strconv"

	"github.com/SlootSantos/janus-server/pkg/api/auth"
	"github.com/SlootSantos/janus-server/pkg/queue"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudfront"
	"github.com/aws/aws-sdk-go/service/sqs"
)

type createOriginConfig struct {
	bucketID string
	owner    string
}

func (c *CDN) createOrginAccess(ctx context.Context, config *createOriginConfig) (*string, error) {
	accessID, err := c.cdn.CreateCloudFrontOriginAccessIdentity(&cloudfront.CreateCloudFrontOriginAccessIdentityInput{
		CloudFrontOriginAccessIdentityConfig: &cloudfront.OriginAccessIdentityConfig{
			Comment:         aws.String("source-cdn-" + config.bucketID),
			CallerReference: aws.String(config.bucketID),
		},
	})
	if err != nil {
		log.Println(err)
		return aws.String(""), err
	}

	isThirdPartyStr := strconv.FormatBool(ctx.Value(auth.ContextKeyIsThirdParty).(bool))

	c.queue.AccessID.Push(queue.QueueMessage{
		queue.MessageDestroyBucketID: &sqs.MessageAttributeValue{
			DataType:    aws.String("String"),
			StringValue: &config.bucketID,
		},
		queue.MessageDestroyAccessID: &sqs.MessageAttributeValue{
			DataType:    aws.String("String"),
			StringValue: accessID.CloudFrontOriginAccessIdentity.Id,
		},
		queue.MessageCommonUser: &sqs.MessageAttributeValue{
			DataType:    aws.String("String"),
			StringValue: aws.String(config.owner),
		},
		queue.MessageCommonIsThirdParty: &sqs.MessageAttributeValue{
			DataType:    aws.String("String"),
			StringValue: aws.String(isThirdPartyStr),
		},
	})

	return accessID.CloudFrontOriginAccessIdentity.Id, nil
}
