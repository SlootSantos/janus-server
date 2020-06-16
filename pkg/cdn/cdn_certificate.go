package cdn

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/SlootSantos/janus-server/pkg/api/auth"
	"github.com/SlootSantos/janus-server/pkg/queue"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/acm"
	"github.com/aws/aws-sdk-go/service/cloudfront"
	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/aws/aws-sdk-go/service/sqs"
)

type issueCertificateConfig struct {
	subdomain string
	distroID  string
	owner     string
}

type createCertificateDNSConfig struct {
	*issueCertificateConfig
	certifcateARN string
}

func (c *CDN) issueCertificate(ctx context.Context, config *issueCertificateConfig) string {
	fullyQualifiedDomain := config.subdomain + "." + c.config.domain

	res, err := c.acm.RequestCertificate(&acm.RequestCertificateInput{
		DomainName:       aws.String(fullyQualifiedDomain),
		ValidationMethod: aws.String("DNS"),
		SubjectAlternativeNames: []*string{
			aws.String("*." + fullyQualifiedDomain),
			aws.String("*.pr." + fullyQualifiedDomain),
		},
	})
	if err != nil {
		panic("Error when creating cert" + err.Error())
	}

	if *res.CertificateArn == "" {
		panic("No certifcate ARN")
	}

	go func() {
		for {
			time.Sleep(time.Second * 5)
			certDNSConfig := &createCertificateDNSConfig{
				config,
				*res.CertificateArn,
			}
			done := c.createCertificateDNSRecords(ctx, certDNSConfig)
			if done {
				break
			}
		}
	}()

	return *res.CertificateArn
}

func (c *CDN) createCertificateDNSRecords(ctx context.Context, config *createCertificateDNSConfig) bool {
	rr, err := c.acm.DescribeCertificate(&acm.DescribeCertificateInput{
		CertificateArn: &config.certifcateARN,
	})
	if err != nil {
		log.Println("Error describing certificate")
		return false
	}

	if len(rr.Certificate.DomainValidationOptions) == 0 {
		return false
	}

	dnsEntryChanges := []*route53.Change{}
	createdChanges := make(map[string]bool)
	for _, validationOption := range rr.Certificate.DomainValidationOptions {
		if validationOption.ResourceRecord == nil || *validationOption.ResourceRecord.Name == "" {
			continue

		}

		if _, ok := createdChanges[*validationOption.ResourceRecord.Name]; ok {
			continue
		}

		dnsEntryChanges = append(dnsEntryChanges, &route53.Change{
			Action: aws.String(dnsActionUpsert),
			ResourceRecordSet: &route53.ResourceRecordSet{
				Name: validationOption.ResourceRecord.Name,
				Type: validationOption.ResourceRecord.Type,
				ResourceRecords: []*route53.ResourceRecord{
					{
						Value: validationOption.ResourceRecord.Value,
					},
				},
				TTL:           aws.Int64(60),
				Weight:        aws.Int64(1),
				SetIdentifier: aws.String("Custom PR preview domain: " + *validationOption.ResourceRecord.Value),
			},
		})

		createdChanges[*validationOption.ResourceRecord.Name] = true
	}

	recordParams := &route53.ChangeResourceRecordSetsInput{
		ChangeBatch: &route53.ChangeBatch{
			Changes: dnsEntryChanges,
			Comment: aws.String("DNS Validation for PR preview."),
		},
		// HostedZoneId: aws.String("/hostedzone/" + os.Getenv("DOMAIN_ZONE_ID")),
		HostedZoneId: aws.String(c.config.hostedZoneID),
	}

	_, err = c.dns.ChangeResourceRecordSets(recordParams)
	if err != nil {
		fmt.Println(err.Error())
	}

	isThirdPartyStr := strconv.FormatBool(ctx.Value(auth.ContextKeyIsThirdParty).(bool))

	c.queue.Certificate.Push(queue.QueueMessage{
		queue.MessageCertificateARN: &sqs.MessageAttributeValue{
			DataType:    aws.String("String"),
			StringValue: &config.certifcateARN,
		},
		queue.MessageCertificateDistroID: &sqs.MessageAttributeValue{
			DataType:    aws.String("String"),
			StringValue: aws.String(config.distroID),
		},
		queue.MessageCertificateSubDomain: &sqs.MessageAttributeValue{
			DataType:    aws.String("String"),
			StringValue: aws.String(config.subdomain),
		},
		queue.MessageCommonUser: &sqs.MessageAttributeValue{
			DataType:    aws.String("String"),
			StringValue: aws.String(config.owner),
		},
		queue.MessageCommonIsThirdParty: &sqs.MessageAttributeValue{
			DataType:    aws.String("String"),
			StringValue: aws.String(isThirdPartyStr),
		},
	})

	return true
}

func (c *CDN) HandleQueueMessageCertificate(distroID string, certificateARN string, subdomain string) (ack bool) {
	getDistroInput := &cloudfront.GetDistributionInput{
		Id: aws.String(distroID),
	}

	output, err := c.cdn.GetDistribution(getDistroInput)
	if err != nil {
		log.Println("Error! getting distro", err)
		return false
	}

	conf := *output.Distribution.DistributionConfig
	conf.ViewerCertificate = constructCertificate(certificateARN)
	conf.Aliases = constructAliases(subdomain, c.config.domain)

	input := &cloudfront.UpdateDistributionInput{
		DistributionConfig: &conf,
		IfMatch:            output.ETag,
		Id:                 output.Distribution.Id,
	}

	_, err = c.cdn.UpdateDistribution(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case cloudfront.ErrCodeInvalidViewerCertificate:
				return ack
			default:
				log.Println("Could not handle AWS err", err.Error())
				return ack
			}
		}
	}

	ack = true
	return ack
}

func (c *CDN) destroyCertificate(certARN string) error {
	_, err := c.acm.DeleteCertificate(&acm.DeleteCertificateInput{
		CertificateArn: aws.String(certARN),
	})

	return err
}
