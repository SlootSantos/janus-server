data "aws_elastic_beanstalk_hosted_zone" "current" {}

data "aws_route53_zone" "janus_dns" {
  name         = "stackers.io."
}

data "aws_acm_certificate" "janus_dns_cert_euc_wild" {
  domain = "*.stackers.io"
}

# resource "aws_route53_record" "janus_website_cdn_dns" {
#   zone_id = data.aws_route53_zone.janus_dns.zone_id
#   name    = data.aws_route53_zone.janus_dns.name
#   type    = "A"

#   alias {
#    name                   = aws_cloudfront_distribution.janus_cdn.domain_name
#    zone_id                = aws_cloudfront_distribution.janus_cdn.hosted_zone_id
#    evaluate_target_health = false
#   }
# }

resource "aws_route53_record" "janus_api_dns" {
  zone_id = data.aws_route53_zone.janus_dns.zone_id
  name    = "api.stackers.io"
  type    = "CNAME"
  ttl     = "300"
  records = [aws_elastic_beanstalk_environment.janus-server-production.cname]
}

output "domain_certificate_arn" {
  value = data.aws_acm_certificate.janus_dns_cert_euc_wild.arn
}

