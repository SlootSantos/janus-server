//go:generate mockgen -source=$GOFILE -destination=mock.$GOFILE -package=$GOPACKAGE

package cdn

import (
	"log"

	"github.com/SlootSantos/janus-server/pkg/queue"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudfront"
	"github.com/aws/aws-sdk-go/service/route53"
)

type cdnandler interface {
	CreateDistribution(*cloudfront.CreateDistributionInput) (*cloudfront.CreateDistributionOutput, error)
	DeleteDistribution(*cloudfront.DeleteDistributionInput) (*cloudfront.DeleteDistributionOutput, error)
	GetDistribution(*cloudfront.GetDistributionInput) (*cloudfront.GetDistributionOutput, error)
	UpdateDistribution(*cloudfront.UpdateDistributionInput) (*cloudfront.UpdateDistributionOutput, error)
	GetCloudFrontOriginAccessIdentity(*cloudfront.GetCloudFrontOriginAccessIdentityInput) (*cloudfront.GetCloudFrontOriginAccessIdentityOutput, error)
	CreateCloudFrontOriginAccessIdentity(*cloudfront.CreateCloudFrontOriginAccessIdentityInput) (*cloudfront.CreateCloudFrontOriginAccessIdentityOutput, error)
}

type dnshandler interface {
	ChangeResourceRecordSets(*route53.ChangeResourceRecordSetsInput) (*route53.ChangeResourceRecordSetsOutput, error)
}

// CDN contains all data to interact w/ AWS Cloudfront
type CDN struct {
	cdn   cdnandler
	dns   dnshandler
	queue *queue.Q
}

// New creates a new CDN creator
func New(s *session.Session, q *queue.Q) *CDN {
	log.Print("DONE: setting up CDN-Creator")

	cdn := &CDN{
		cdn:   cloudfront.New(s),
		dns:   route53.New(s),
		queue: q,
	}

	q.DestroyCDN.SetListener(cdn.deleteDisabledDistro)

	return cdn
}
