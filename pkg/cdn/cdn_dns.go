package cdn

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/route53"
)

const (
	dnsActionUpsert            = "UPSERT"
	dnsActionDelete            = "DELETE"
	aliasPrefixGreenDeployment = "green."
	aliasPrefixDevelopmentEnv  = "dev."
	aliasPrefixStageEnv        = "stage."
	aliasPrefixPRPreview       = "*.pr."
)

func (c *CDN) createDNSRecord(distroDomain string, subdomain string) {
	c.handleCommonRoute53Change(dnsActionUpsert, subdomain, distroDomain)
}

func (c *CDN) destroyDNSRecord(distroDomain string, subdomain string) {
	c.handleCommonRoute53Change(dnsActionDelete, subdomain, distroDomain)
}

func (c *CDN) handleCommonRoute53Change(action string, subdomain string, target string) {
	log.Println("SUBDOMAIN IN DNS", subdomain)
	recordParams := &route53.ChangeResourceRecordSetsInput{
		ChangeBatch: &route53.ChangeBatch{
			Changes: []*route53.Change{
				{
					Action: aws.String(action),
					ResourceRecordSet: &route53.ResourceRecordSet{
						Name: aws.String(subdomain + "." + c.config.domain),
						Type: aws.String("CNAME"),
						ResourceRecords: []*route53.ResourceRecord{
							{
								Value: aws.String(target),
							},
						},
						TTL:           aws.Int64(60),
						Weight:        aws.Int64(1),
						SetIdentifier: aws.String("Custom Domain CNAME for stackers CDN: " + target),
					},
				},
				{
					Action: aws.String(action),
					ResourceRecordSet: &route53.ResourceRecordSet{
						Name: aws.String(aliasPrefixGreenDeployment + subdomain + "." + c.config.domain),
						Type: aws.String("CNAME"),
						ResourceRecords: []*route53.ResourceRecord{
							{
								Value: aws.String(target),
							},
						},
						TTL:           aws.Int64(60),
						Weight:        aws.Int64(1),
						SetIdentifier: aws.String("Custom Domain CNAME for stackers CDN: " + target),
					},
				},
				{
					Action: aws.String(action),
					ResourceRecordSet: &route53.ResourceRecordSet{
						Name: aws.String(aliasPrefixDevelopmentEnv + subdomain + "." + c.config.domain),
						Type: aws.String("CNAME"),
						ResourceRecords: []*route53.ResourceRecord{
							{
								Value: aws.String(target),
							},
						},
						TTL:           aws.Int64(60),
						Weight:        aws.Int64(1),
						SetIdentifier: aws.String("Custom Domain CNAME for stackers CDN: " + target),
					},
				},
				{
					Action: aws.String(action),
					ResourceRecordSet: &route53.ResourceRecordSet{
						Name: aws.String(aliasPrefixStageEnv + subdomain + "." + c.config.domain),
						Type: aws.String("CNAME"),
						ResourceRecords: []*route53.ResourceRecord{
							{
								Value: aws.String(target),
							},
						},
						TTL:           aws.Int64(60),
						Weight:        aws.Int64(1),
						SetIdentifier: aws.String("Custom Domain CNAME for stackers CDN: " + target),
					},
				},
				{
					Action: aws.String(action),
					ResourceRecordSet: &route53.ResourceRecordSet{
						Name: aws.String(aliasPrefixPRPreview + subdomain + "." + c.config.domain),
						Type: aws.String("CNAME"),
						ResourceRecords: []*route53.ResourceRecord{
							{
								Value: aws.String(target),
							},
						},
						TTL:           aws.Int64(60),
						Weight:        aws.Int64(1),
						SetIdentifier: aws.String("Custom Domain CNAME for stackers CDN: " + target),
					},
				},
			},
			Comment: aws.String("Sample update."),
		},
		HostedZoneId: aws.String(c.config.hostedZoneID),
		// HostedZoneId: aws.String("/hostedzone/" + os.Getenv("DOMAIN_ZONE_ID")),
	}

	_, err := c.dns.ChangeResourceRecordSets(recordParams)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}
