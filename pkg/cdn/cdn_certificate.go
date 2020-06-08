package cdn

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/SlootSantos/janus-server/pkg/queue"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/acm"
	"github.com/aws/aws-sdk-go/service/cloudfront"
	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/aws/aws-sdk-go/service/sqs"
)

func (c *CDN) issueCertificate(subdomain string, distroID string) {
	fullyQualifiedDomain := subdomain + "." + os.Getenv("DOMAIN_HOST")

	res, _ := c.acm.RequestCertificate(&acm.RequestCertificateInput{
		DomainName:       aws.String(fullyQualifiedDomain),
		ValidationMethod: aws.String("DNS"),
		SubjectAlternativeNames: []*string{
			aws.String("*." + fullyQualifiedDomain),
			aws.String("*.pr." + fullyQualifiedDomain),
		},
	})

	for {
		time.Sleep(time.Second * 5)
		done := c.createCertificateDNSRecords(*res.CertificateArn, distroID, subdomain)
		if done {
			break
		}
	}
}

func (c *CDN) createCertificateDNSRecords(certARN string, distroID string, subdomain string) bool {
	rr, err := c.acm.DescribeCertificate(&acm.DescribeCertificateInput{
		CertificateArn: &certARN,
	})

	if len(rr.Certificate.DomainValidationOptions) == 0 {
		return false
	}

	dnsEntryChanges := []*route53.Change{}
	createdChanges := make(map[string]bool)
	for _, validationOption := range rr.Certificate.DomainValidationOptions {
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
		HostedZoneId: aws.String("/hostedzone/" + os.Getenv("DOMAIN_ZONE_ID")),
	}

	_, err = c.dns.ChangeResourceRecordSets(recordParams)
	if err != nil {
		fmt.Println(err.Error())
	}

	c.queue.Certificate.Push(queue.QueueMessage{
		queue.MessageCertificateARN: &sqs.MessageAttributeValue{
			DataType:    aws.String("String"),
			StringValue: &certARN,
		},
		queue.MessageCertificateDistroID: &sqs.MessageAttributeValue{
			DataType:    aws.String("String"),
			StringValue: aws.String(distroID),
		},
		queue.MessageCertificateSubDomain: &sqs.MessageAttributeValue{
			DataType:    aws.String("String"),
			StringValue: aws.String(subdomain),
		},
	})

	return true
}

func (c *CDN) updateCDNCertificate(message queue.QueueMessage) (ack bool) {
	log.Println("RECEIVING CERT UPDATE MESSAGE")
	distroID, ok := message[queue.MessageCertificateDistroID]
	if !ok {
		log.Println("Handling queue message for CDN failed. attribute:", queue.MessageCertificateDistroID, " does not exist on message")
		return ack
	}

	certificateARN, ok := message[queue.MessageCertificateARN]
	if !ok {
		log.Println("Handling queue message for CDN failed. attribute:", queue.MessageCertificateARN, " does not exist on message")
		return ack
	}

	subdomain, ok := message[queue.MessageCertificateSubDomain]
	if !ok {
		log.Println("Handling queue message for CDN failed. attribute:", queue.MessageCertificateSubDomain, " does not exist on message")
		return ack
	}

	getDistroInput := &cloudfront.GetDistributionInput{
		Id: aws.String(*distroID.StringValue),
	}

	output, err := c.cdn.GetDistribution(getDistroInput)
	if err != nil {
		log.Println("Error! getting distro", err)
		return false
	}

	conf := *output.Distribution.DistributionConfig
	conf.ViewerCertificate = constructCertificate(*certificateARN.StringValue)
	conf.Aliases = constructAliases(*subdomain.StringValue)

	input := &cloudfront.UpdateDistributionInput{
		DistributionConfig: &conf,
		IfMatch:            output.ETag,
		Id:                 output.Distribution.Id,
	}

	_, err = c.cdn.UpdateDistribution(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			fmt.Println(aerr.Error())

			return false
		}
	}

	ack = true
	return ack
}
