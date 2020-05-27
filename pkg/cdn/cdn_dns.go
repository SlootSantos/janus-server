package cdn

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/route53"
)

const (
	dnsActionUpsert       = "UPSERT"
	dnsActionDelete       = "DELETE"
	greenDeploymentPrefix = "green-"
)

func (c *CDN) createDNSRecord(distroDomain string, stackID string) {
	c.handleCommonRoute53Change(dnsActionUpsert, stackID, distroDomain)
}

func (c *CDN) destroyDNSRecord(distroDomain string, stackID string) {
	c.handleCommonRoute53Change(dnsActionDelete, stackID, distroDomain)
}

func (c *CDN) handleCommonRoute53Change(action string, subdomain string, target string) {
	recordParams := &route53.ChangeResourceRecordSetsInput{
		ChangeBatch: &route53.ChangeBatch{
			Changes: []*route53.Change{
				{
					Action: aws.String(action),
					ResourceRecordSet: &route53.ResourceRecordSet{
						Name: aws.String(subdomain + "." + os.Getenv("DOMAIN_HOST")),
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
						Name: aws.String(greenDeploymentPrefix + subdomain + "." + os.Getenv("DOMAIN_HOST")),
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
		HostedZoneId: aws.String("/hostedzone/" + os.Getenv("DOMAIN_ZONE_ID")),
	}

	_, err := c.dns.ChangeResourceRecordSets(recordParams)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}
