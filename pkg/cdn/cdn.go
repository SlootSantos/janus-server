package cdn

import (
	"log"

	"github.com/SlootSantos/janus-server/pkg/queue"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudfront"
)

// CDN contains all data to interact w/ AWS Cloudfront
type CDN struct {
	cdn   *cloudfront.CloudFront
	queue *queue.Q
}

// New creates a new CDN creator
func New(s *session.Session, q *queue.Q) *CDN {
	log.Print("DONE: setting up CDN-Creator")

	cdn := &CDN{
		cdn:   cloudfront.New(s),
		queue: q,
	}

	q.DestroyCDN.SetListener(cdn.deleteDisabledDistro)

	return cdn
}
