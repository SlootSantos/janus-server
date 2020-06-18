data "aws_elastic_beanstalk_hosted_zone" "current" {}

data "aws_route53_zone" "janus_dns" {
  name         = "${var.janus_dns_host_name}."
}

data "aws_acm_certificate" "janus_dns_cert_euc_wild" {
  domain = "${var.janus_dns_host_name}"
}

resource "aws_route53_record" "janus_api_dns" {
  zone_id = data.aws_route53_zone.janus_dns.zone_id
  name    = "${var.janus_dns_api_name}.${var.janus_dns_host_name}"
  type    = "CNAME"
  ttl     = "300"
  records = [aws_elastic_beanstalk_environment.janus-server-production.cname]
}

output "domain_certificate_arn" {
  value = data.aws_acm_certificate.janus_dns_cert_euc_wild.arn
}

