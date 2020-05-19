package cdn

import (
	"context"
	"fmt"
	"log"

	"github.com/SlootSantos/janus-server/pkg/jam"
	"github.com/SlootSantos/janus-server/pkg/queue"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/cloudfront"
	"github.com/aws/aws-sdk-go/service/sqs"
)

const cdnPrefix = "janus_cdn_distro"
const cdnS3OriginSuffix = ".s3.us-east-1.amazonaws.com"

// Create generates a Cloudfront-Distro at AWS
func (c *CDN) Create(ctx context.Context, param *jam.CreationParam, out *jam.OutputParam) (string, error) {
	fmt.Println("STARTING: creating bucket ID:", param.Bucket.ID)
	bucketID := param.Bucket.ID

	accessID, err := c.createOrginAccess(bucketID)
	if err != nil {
		return "", err
	}

	config := c.constructStandardDistroConfig(bucketID, *accessID, param.ID)
	createDistroOuput, err := c.cdn.CreateDistribution(config)
	if err != nil {
		fmt.Println(err)
	}

	out.CDN = &jam.StackCDN{
		Domain:   *createDistroOuput.Distribution.DomainName,
		ID:       *createDistroOuput.Distribution.Id,
		AccessID: *accessID,
	}

	c.createDNSRecord(*createDistroOuput.Distribution.DomainName, param.ID)

	log.Println("DONE: creating up CDN ID:", out.CDN.ID)
	return "", nil
}

// Destroy deletes a Cloudfront-Distro at AWS
func (c *CDN) Destroy(ctx context.Context, param *jam.DeletionParam) error {
	log.Println("START: destroying CDN")

	in := &cloudfront.GetDistributionInput{
		Id: aws.String(param.CDN.ID),
	}

	output, err := c.cdn.GetDistribution(in)
	if err != nil {
		return err
	}

	conf := *output.Distribution.DistributionConfig
	conf.Enabled = aws.Bool(false)

	input := &cloudfront.UpdateDistributionInput{
		Id:                 output.Distribution.Id,
		DistributionConfig: &conf,
		IfMatch:            output.ETag,
	}

	res, err := c.cdn.UpdateDistribution(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			fmt.Println(aerr.Error())
		}
	}

	c.queue.DestroyCDN.Push(
		queue.QueueMessage{
			queue.MessageAccessDistroID: &sqs.MessageAttributeValue{
				DataType:    aws.String("String"),
				StringValue: output.Distribution.Id,
			},
			queue.MessageAccessEtag: &sqs.MessageAttributeValue{
				DataType:    aws.String("String"),
				StringValue: res.ETag,
			},
		},
	)

	c.destroyDNSRecord(param.CDN.Domain, param.ID)

	log.Println("DONE: destroying CDN ID:", param.CDN.ID)
	return nil
}

// List returns a list of all Cloudfront-Distro at AWS
func (c *CDN) List(ctx context.Context) string {
	return "CDN_1"
}
