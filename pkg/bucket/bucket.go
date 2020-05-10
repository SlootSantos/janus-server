//go:generate mockgen -source=$GOFILE -destination=mock.$GOFILE -package=$GOPACKAGE

package bucket

import (
	"log"

	"github.com/SlootSantos/janus-server/pkg/queue"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type bucketHandler interface {
	ListBuckets(*s3.ListBucketsInput) (*s3.ListBucketsOutput, error)
	CreateBucket(*s3.CreateBucketInput) (*s3.CreateBucketOutput, error)
	DeleteBucket(*s3.DeleteBucketInput) (*s3.DeleteBucketOutput, error)
	ListObjects(*s3.ListObjectsInput) (*s3.ListObjectsOutput, error)
	DeleteObjects(*s3.DeleteObjectsInput) (*s3.DeleteObjectsOutput, error)
	PutBucketPolicy(*s3.PutBucketPolicyInput) (*s3.PutBucketPolicyOutput, error)
}

// Bucket contains all data to interact w/ AWS S3
type Bucket struct {
	s3    bucketHandler
	queue *queue.Q
}

// New creates a new Bucket creator
func New(s *session.Session, q *queue.Q) *Bucket {
	log.Print("DONE: setting up Bucket-Creator")

	bucket := &Bucket{
		s3:    s3.New(s),
		queue: q,
	}

	q.AccessID.SetListener(bucket.setBucketAccessID)

	return bucket
}
