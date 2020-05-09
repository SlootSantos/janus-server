package bucket

import (
	"context"
	"errors"
	"log"

	"github.com/SlootSantos/janus-server/pkg/jam"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

const constantSuccessResponse = "Success"
const constantPublicReadACL = "public-read"
const constantPrivate = "private"

// Create generates a S3-Bucket at AWS
func (b *Bucket) Create(ctx context.Context, param *jam.CreationParam, out *jam.OutputParam) (string, error) {
	createBucketInput := &s3.CreateBucketInput{
		Bucket: aws.String(param.Bucket.ID),
		ACL:    aws.String(constantPrivate),
	}

	_, err := b.s3.CreateBucket(createBucketInput)
	if err != nil {
		log.Print("FAILED: creating ID:", param.Bucket.ID)
		return "", errors.New("bucket could not be created: " + err.Error())
	}

	out.BucketID = param.Bucket.ID
	log.Print("DONE: creating up Bucket ID:", param.Bucket.ID)

	return param.Bucket.ID, nil
}

// Destroy deletes a S3-Bucket at AWS
func (b *Bucket) Destroy(ctx context.Context, param *jam.DeletionParam) error {
	log.Println("START: destroying bucket", param.BucketID)

	err := b.emptyBucket(param.BucketID)
	if err != nil {
		log.Println(err)
		return err
	}

	deleteBucketInput := &s3.DeleteBucketInput{
		Bucket: aws.String(param.BucketID),
	}

	_, err = b.s3.DeleteBucket(deleteBucketInput)
	if err != nil {
		log.Println(err)
	}

	log.Println("DONE: destroying Bucket ID:", param.BucketID)
	return nil
}

// List returns a list of all S3-Buckets at AWS
func (b *Bucket) List(ctx context.Context) string {
	res, err := b.s3.ListBuckets(&s3.ListBucketsInput{})
	if err != nil {
		log.Println("FUCK", err)
		return ""
	}

	return res.String()
}
