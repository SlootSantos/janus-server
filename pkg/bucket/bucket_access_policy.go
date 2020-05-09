package bucket

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/SlootSantos/janus-server/pkg/queue"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

type bucketPolicy map[string]interface{}

func (b *Bucket) setBucketAccessID(message queue.QueueMessage) (ack bool) {
	bucketID, ok := message[queue.MessageDestroyBucketID]
	if !ok {
		log.Println("Handling queue message for Bucket-Policy failed. attribute:", queue.MessageDestroyBucketID, " does not exist on message")
		return ack
	}

	accessID, ok := message[queue.MessageDestroyAccessID]
	if !ok {
		log.Println("Handling queue message for Bucket-Policy failed. attribute:", queue.MessageDestroyAccessID, " does not exist on message")
		return ack
	}

	policyString := createPolicy(*accessID.StringValue, *bucketID.StringValue)
	policy, err := json.Marshal(policyString)
	if err != nil {
		log.Println(err)
		return ack
	}

	_, err = b.s3.PutBucketPolicy(&s3.PutBucketPolicyInput{
		Bucket: bucketID.StringValue,
		Policy: aws.String(string(policy)),
	})
	if err != nil {
		return ack
	}

	ack = true
	return ack
}

func createPolicy(accessID string, bucketID string) bucketPolicy {
	principalString := "arn:aws:iam::cloudfront:user/CloudFront Origin Access Identity " + accessID
	readOnlyAnonUserPolicy := bucketPolicy{
		"Version": "2012-10-17",
		"Statement": []map[string]interface{}{
			{
				"Sid":    "",
				"Effect": "Allow",
				"Principal": map[string]string{
					"AWS": principalString,
				},
				"Action": []string{
					"s3:GetObject",
				},
				"Resource": []string{
					fmt.Sprintf("arn:aws:s3:::%s/*", bucketID),
				},
			},
			{
				"Sid":    "",
				"Effect": "Allow",
				"Principal": map[string]string{
					"AWS": principalString,
				},
				"Action": []string{
					"s3:ListBucket",
				},
				"Resource": []string{
					fmt.Sprintf("arn:aws:s3:::%s", bucketID),
				},
			},
		},
	}

	return readOnlyAnonUserPolicy
}
