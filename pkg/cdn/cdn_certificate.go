package cdn

func (c *CDN) issueCertificate(domain string) {
	// res, err := c.acm.RequestCertificate(&acm.RequestCertificateInput{
	// 	// DomainName:       aws.String("*.mywrkspace.com"),
	// 	DomainName:       aws.String(domain),
	// 	ValidationMethod: aws.String("DNS"),
	// })

	// log.Println("err", err)
	// log.Println("RES", res.String())

	// next =>
	// go and get certificaate with arn from "res"
	// print "go and create cname for domain with <validation record>"
}
