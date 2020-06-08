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
	subdomain := param.CDN.Subdomain

	accessID, err := c.createOrginAccess(bucketID)
	if err != nil {
		return "", err
	}

	distroConfigInput := &constructDistroConfigInput{
		bucketID:       bucketID,
		originAccessID: *accessID,
		stackID:        param.ID,
		subdomain:      subdomain,
	}

	config := c.constructStandardDistroConfig(distroConfigInput)
	createDistroOuput, err := c.cdn.CreateDistribution(config)
	if err != nil {
		log.Println("could not create Cloudfront distro", err)
	}

	c.createDNSRecord(*createDistroOuput.Distribution.DomainName, subdomain)
	c.issueCertificate(subdomain, *createDistroOuput.Distribution.Id)

	out.CDN = &jam.StackCDN{
		Subdomain: subdomain,
		AccessID:  *accessID,
		Domain:    *createDistroOuput.Distribution.DomainName,
		ID:        *createDistroOuput.Distribution.Id,
	}

	log.Println("DONE: creating up CDN ID:", out.CDN.ID)
	return "", nil
}

// Destroy deletes a Cloudfront-Distro at AWS
func (c *CDN) Destroy(ctx context.Context, param *jam.DeletionParam) error {
	log.Println("START: destroying CDN")

	getDistroInput := &cloudfront.GetDistributionInput{
		Id: aws.String(param.CDN.ID),
	}

	output, err := c.cdn.GetDistribution(getDistroInput)
	if err != nil {
		return err
	}

	conf := *output.Distribution.DistributionConfig
	conf.Enabled = aws.Bool(false)

	input := &cloudfront.UpdateDistributionInput{
		DistributionConfig: &conf,
		IfMatch:            output.ETag,
		Id:                 output.Distribution.Id,
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
				StringValue: output.Distribution.Id,
				DataType:    aws.String("String"),
			},
			queue.MessageAccessEtag: &sqs.MessageAttributeValue{
				StringValue: res.ETag,
				DataType:    aws.String("String"),
			},
		},
	)

	c.destroyDNSRecord(param.CDN.Domain, param.CDN.Subdomain)

	log.Println("DONE: destroying CDN ID:", param.CDN.ID)
	return nil
}

// List returns a list of all Cloudfront-Distro at AWS
func (c *CDN) List(ctx context.Context) string {
	return "CDN_1"
}
