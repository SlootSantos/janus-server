//go:generate mockgen -source=$GOFILE -destination=mock.$GOFILE -package=$GOPACKAGE

package cdn

import (
	"log"

	"github.com/SlootSantos/janus-server/pkg/queue"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/acm"
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

type certificateHandler interface {
	RequestCertificate(*acm.RequestCertificateInput) (*acm.RequestCertificateOutput, error)
	GetCertificate(*acm.GetCertificateInput) (*acm.GetCertificateOutput, error)
	DescribeCertificate(*acm.DescribeCertificateInput) (*acm.DescribeCertificateOutput, error)
	DeleteCertificate(*acm.DeleteCertificateInput) (*acm.DeleteCertificateOutput, error)
}

// CDN contains all data to interact w/ AWS Cloudfront
type CDN struct {
	cdn    cdnandler
	dns    dnshandler
	queue  *queue.Q
	acm    certificateHandler
	config *cdnConfig
}

type CreateCDNParams struct {
	Domain       string
	HostedZoneID string
	CertARN      string
	LambdaARN    string
	Session      *session.Session
	Queue        *queue.Q
}

type cdnConfig struct {
	domain       string
	hostedZoneID string
	certARN      string
	lambdaARN    string
}

// New creates a new CDN creator
func New(params *CreateCDNParams) *CDN {
	log.Print("DONE: setting up CDN-Creator")

	cdn := &CDN{
		cdn:   cloudfront.New(params.Session),
		dns:   route53.New(params.Session),
		queue: params.Queue,
		acm:   acm.New(params.Session),
		config: &cdnConfig{
			domain:       params.Domain,
			hostedZoneID: params.HostedZoneID,
			certARN:      params.CertARN,
			lambdaARN:    params.LambdaARN,
		},
	}

	return cdn
}

func (c *CDN) HandleQueueMessaeDestroyCDN(distroID string, etag string, certARN string) (ack bool) {
	deleteDistroInput := &cloudfront.DeleteDistributionInput{
		Id:      &distroID,
		IfMatch: &etag,
	}

	_, err := c.cdn.DeleteDistribution(deleteDistroInput)
	if err != nil {
		log.Println("could not delete distro", err.Error())
		return ack
	}

	err = c.destroyCertificate(certARN)
	if err != nil {
		log.Println("could not delete distro Certificate", err.Error())
		return ack
	}

	ack = true
	return ack
}
