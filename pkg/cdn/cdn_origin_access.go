package cdn

import (
	"log"
	"time"

	"github.com/SlootSantos/janus-server/pkg/queue"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudfront"
	"github.com/aws/aws-sdk-go/service/sqs"
)

func (c *CDN) createOrginAccess(bucketID string) (*string, error) {
	accessID, err := c.cdn.CreateCloudFrontOriginAccessIdentity(&cloudfront.CreateCloudFrontOriginAccessIdentityInput{
		CloudFrontOriginAccessIdentityConfig: &cloudfront.OriginAccessIdentityConfig{
			Comment:         aws.String("source-cdn-" + bucketID),
			CallerReference: aws.String(time.Now().String() + bucketID),
		},
	})
	if err != nil {
		log.Println(err)
		return aws.String(""), err
	}

	originAccess, err := c.cdn.GetCloudFrontOriginAccessIdentity(&cloudfront.GetCloudFrontOriginAccessIdentityInput{
		Id: accessID.CloudFrontOriginAccessIdentity.Id,
	})
	if err != nil {
		log.Println("IDENTIITY ERROR", err)
		return aws.String(""), err
	}

	c.queue.AccessID.Push(queue.QueueMessage{
		queue.MessageDestroyBucketID: &sqs.MessageAttributeValue{
			DataType:    aws.String("String"),
			StringValue: &bucketID,
		},
		queue.MessageDestroyAccessID: &sqs.MessageAttributeValue{
			DataType:    aws.String("String"),
			StringValue: accessID.CloudFrontOriginAccessIdentity.Id,
		},
	})

	return originAccess.CloudFrontOriginAccessIdentity.Id, nil
}
