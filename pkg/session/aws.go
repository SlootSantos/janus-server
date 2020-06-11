package session

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
)

// AWSSession authenticates against AWS and returns a session object
func AWSSession() (*session.Session, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"),
	})

	return sess, err
}

// AWSSession authenticates against AWS and returns a session object
func AWSSessionThirdParty(access string, secret string) (*session.Session, error) {
	creds := credentials.NewStaticCredentials(access, secret, "")
	sess, err := session.NewSession(&aws.Config{
		Credentials: creds,
		Region:      aws.String("us-east-1"),
	})

	return sess, err
}
