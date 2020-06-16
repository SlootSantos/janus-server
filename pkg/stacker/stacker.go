package stacker

import (
	"context"
	"os"

	aws "github.com/aws/aws-sdk-go/aws/session"

	"github.com/SlootSantos/janus-server/pkg/bucket"
	"github.com/SlootSantos/janus-server/pkg/cdn"
	"github.com/SlootSantos/janus-server/pkg/jam"
	"github.com/SlootSantos/janus-server/pkg/queue"
	"github.com/SlootSantos/janus-server/pkg/repo"
	"github.com/SlootSantos/janus-server/pkg/session"
)

type Stacker struct {
	queue           queue.Q
	repo            *repo.Repo
	_defaultCreator *jam.Creator
}

type launchParam struct {
	hasOwnCreds  bool
	creds        awsCredentials
	domain       string
	lambdaARN    string
	hostedZoneID string
}

type awsCredentials struct {
	accessKey string
	secretKey string
}

// New does what ever
func New(defaultSession *aws.Session) *Stacker {
	stacker := &Stacker{
		repo:  repo.New(),
		queue: queue.New(defaultSession),
	}

	setupQueueHandlers(&stacker.queue)
	setupDefaultCreator(defaultSession, stacker)

	return stacker
}

func (s *Stacker) launchCreator(ctx context.Context, launchParamOps *launchParamOptions) (*jam.Creator, error) {
	param, err := setLaunchParam(ctx, launchParamOps)

	if err != nil {
		return nil, err
	}

	if !param.hasOwnCreds || !launchParamOps.isThirdParty {
		return s._defaultCreator, nil
	}

	awsSess, _ := session.AWSSessionThirdParty(param.creds.accessKey, param.creds.secretKey)
	bucket := bucket.New(awsSess, &s.queue)
	cloudfront := cdn.New(&cdn.CreateCDNParams{
		HostedZoneID: param.hostedZoneID,
		CertARN:      os.Getenv("DOMAIN_CERT_ARN"),
		LambdaARN:    param.lambdaARN,
		Session:      awsSess,
		Domain:       param.domain,
		Queue:        &s.queue,
	})

	return jam.New(bucket, cloudfront, s.repo), nil
}

func setupQueueHandlers(q *queue.Q) {
	q.Certificate.SetListener(updateCDNCertificate)
	q.DestroyCDN.SetListener(deleteDisabledDistro)
	q.AccessID.SetListener(setAccessID)
}

func setupDefaultCreator(defaultSession *aws.Session, s *Stacker) {
	bucket := bucket.New(defaultSession, &s.queue)
	cloudfront := cdn.New(&cdn.CreateCDNParams{
		HostedZoneID: os.Getenv("DOMAIN_ZONE_ID"),
		CertARN:      os.Getenv("DOMAIN_CERT_ARN"),
		LambdaARN:    "arn:aws:lambda:us-east-1:976589619057:function:janus-exmaple-redirect:34",
		Session:      defaultSession,
		Domain:       os.Getenv("DOMAIN_HOST"),
		Queue:        &s.queue,
	})

	s._defaultCreator = jam.New(bucket, cloudfront, s.repo)
}
