package stacker

import (
	"context"
	"log"
	"os"
	"strconv"
	"time"

	aws "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	awsSession "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/lambda"

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
func New(defaultSession *awsSession.Session) *Stacker {
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

func setupDefaultCreator(defaultSession *awsSession.Session, s *Stacker) {
	lambdaARN := setupLambdaFunc(defaultSession)
	bucket := bucket.New(defaultSession, &s.queue)
	cloudfront := cdn.New(&cdn.CreateCDNParams{
		HostedZoneID: os.Getenv("DOMAIN_ZONE_ID"),
		CertARN:      os.Getenv("DOMAIN_CERT_ARN"),
		LambdaARN:    lambdaARN,
		Session:      defaultSession,
		Domain:       os.Getenv("DOMAIN_HOST"),
		Queue:        &s.queue,
	})

	s._defaultCreator = jam.New(bucket, cloudfront, s.repo)
}

func setupLambdaFunc(sess *awsSession.Session) string {
	i := iam.New(sess)
	roleName := "stackers-cdn-routing-lambda-policy"
	functionName := "stackers-cdn-routing-origin-lambda"
	var roleARN string
	role, err := i.CreateRole(&iam.CreateRoleInput{
		Path:                     aws.String("/service-role/"),
		AssumeRolePolicyDocument: aws.String("{\"Version\": \"2012-10-17\",\"Statement\": [{\"Effect\": \"Allow\",\"Principal\": {\"Service\": [\"lambda.amazonaws.com\",\"edgelambda.amazonaws.com\"]},\"Action\": \"sts:AssumeRole\"}]}"),
		RoleName:                 aws.String(roleName),
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case iam.ErrCodeEntityAlreadyExistsException:
				log.Printf("User: already has role named: %s", roleName)
				r, _ := i.GetRole(&iam.GetRoleInput{
					RoleName: aws.String(roleName),
				})
				roleARN = *r.Role.Arn
			default:
				panic("Could not handle AWS err" + err.Error())
			}
		}
	} else {
		roleARN = *role.Role.Arn
		i.PutRolePolicy(&iam.PutRolePolicyInput{
			PolicyName:     aws.String("stackers-lambda-exec-policy"),
			PolicyDocument: aws.String("{\"Version\": \"2012-10-17\", \"Statement\": [ { \"Effect\": \"Allow\", \"Action\": [ \"logs:CreateLogGroup\", \"logs:CreateLogStream\", \"logs:PutLogEvents\" ], \"Resource\": [ \"arn:aws:logs:*:*:*\" ] } ] }"),
			RoleName:       aws.String(roleName),
		})

		time.Sleep(time.Second * 60)
	}

	var functionARN string
	lam := lambda.New(sess)
	res, err := lam.CreateFunction(&lambda.CreateFunctionInput{
		FunctionName: aws.String(functionName),
		Handler:      aws.String("index.handler"),
		Runtime:      aws.String("nodejs12.x"),
		Role:         aws.String(roleARN),
		Code: &lambda.FunctionCode{
			S3Bucket: aws.String("stackers.io-lambda-public-functions"),
			S3Key:    aws.String("exx.zip"),
		},
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case lambda.ErrCodeResourceConflictException:
				log.Printf("User: already has function named: %s", functionName)
				existingFunction, _ := lam.GetFunction(&lambda.GetFunctionInput{
					FunctionName: aws.String(functionName),
				})

				versions, _ := lam.ListVersionsByFunction(&lambda.ListVersionsByFunctionInput{
					FunctionName: &functionName,
				})
				latestVersion := strconv.Itoa(len(versions.Versions) - 1)

				functionARN = *existingFunction.Configuration.FunctionArn + ":" + latestVersion
			default:
				panic("Could not handle AWS err" + err.Error())
			}
		}

	} else {
		_, err := lam.PublishVersion(&lambda.PublishVersionInput{
			FunctionName: aws.String(*res.FunctionArn),
			Description:  aws.String("Stackers.io Routing Lambda Function"),
		})
		if err != nil {
			panic("Could not publish function" + err.Error())
		}

		versions, _ := lam.ListVersionsByFunction(&lambda.ListVersionsByFunctionInput{
			FunctionName: &functionName,
		})
		latestVersion := strconv.Itoa(len(versions.Versions) - 1)

		functionARN = *res.FunctionArn + ":" + latestVersion
	}

	return functionARN
}
