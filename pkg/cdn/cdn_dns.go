package cdn

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/route53"
)

func (c *CDN) createDNSRecord(distroDomain string, stackID string) {
	c.handleCommonRoute53Change("UPSERT", stackID, distroDomain)
}

func (c *CDN) destroyDNSRecord(distroDomain string, stackID string) {
	c.handleCommonRoute53Change("DELETE", stackID, distroDomain)
}

func (c *CDN) handleCommonRoute53Change(action string, alias string, target string) {
	recordParams := &route53.ChangeResourceRecordSetsInput{
		ChangeBatch: &route53.ChangeBatch{
			Changes: []*route53.Change{
				{
					Action: aws.String(action),
					ResourceRecordSet: &route53.ResourceRecordSet{
						Name: aws.String(alias + "." + os.Getenv("DOMAIN_HOST")),
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
	resp, err := c.dns.ChangeResourceRecordSets(recordParams)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return
	}

	// Pretty-print the response data.
	fmt.Println("Change Response:")
	fmt.Println(resp)
}
