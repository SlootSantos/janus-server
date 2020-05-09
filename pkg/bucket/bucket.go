package bucket

import (
	"log"

	"github.com/SlootSantos/janus-server/pkg/queue"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// Bucket contains all data to interact w/ AWS S3
type Bucket struct {
	s3    *s3.S3
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
