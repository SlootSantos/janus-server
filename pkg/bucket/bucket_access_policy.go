package bucket

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

type bucketPolicy map[string]interface{}

func (b *Bucket) HandleQueueMessageAccessID(bucketID string, accessID string) (ack bool) {
	policyString := createPolicy(accessID, bucketID)
	policy, err := json.Marshal(policyString)
	if err != nil {
		log.Println(err)
		return ack
	}

	_, err = b.s3.PutBucketPolicy(&s3.PutBucketPolicyInput{
		Bucket: &bucketID,
		Policy: aws.String(string(policy)),
	})
	if err != nil {
		log.Println("Error setting policy", err.Error())
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
