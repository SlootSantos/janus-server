resource "aws_ssm_parameter" "janus_redis_cluster_add" {
  name  = "/janus/env/production/REDIS_CONN_HOSTNAME"
  type  = "String"
  value = aws_elasticache_cluster.janus_redis.cache_nodes.0.address
  overwrite = true
   tags = {
       "env"     = "production"
       "project" = "janus"
    }
}

resource "aws_ssm_parameter" "janus_redis_cluster_port" {
  name  = "/janus/env/production/REDIS_CONN_PORT"
  type  = "String"
  value = aws_elasticache_cluster.janus_redis.cache_nodes.0.port
  overwrite = true
   tags = {
       "env"     = "production"
       "project" = "janus"
    }
}

resource "aws_ssm_parameter" "janus_ecr_url" {
  name  = "/janus/env/production/ECR_URL"
  type  = "String"
  value = aws_ecr_repository.janus_server.repository_url
  overwrite = true
   tags = {
       "env"     = "production"
       "project" = "janus"
    }
}

resource "aws_ssm_parameter" "janus_sqs_access_id" {
  name  = "/janus/env/production/SQS_URL_ACCESS_ID"
  type  = "String"
  value = aws_sqs_queue.AccessID.id
  overwrite = true
   tags = {
       "env"     = "production"
       "project" = "janus"
    }
}

resource "aws_ssm_parameter" "janus_sqs_destroy_cdn" {
  name  = "/janus/env/production/SQS_URL_DESTROY_CDN"
  type  = "String"
  value = aws_sqs_queue.DestroyCDN.id
  overwrite = true
   tags = {
       "env"     = "production"
       "project" = "janus"
   }
}

resource "aws_ssm_parameter" "janus_domain_cert_arn" {
  name  = "/janus/env/production/DOMAIN_CERT_ARN"
  type  = "String"
  value = data.aws_acm_certificate.janus_dns_cert_euc_wild.arn
  overwrite = true
   tags = {
       "env"     = "production"
       "project" = "janus"
   }
}

resource "aws_ssm_parameter" "janus_domain_host" {
  name  = "/janus/env/production/DOMAIN_HOST"
  type  = "String"
  value = trimsuffix(data.aws_route53_zone.janus_dns.name, ".")
  overwrite = true
   tags = {
       "env"     = "production"
       "project" = "janus"
   }
}

resource "aws_ssm_parameter" "janus_domain_zone_id" {
  name  = "/janus/env/production/DOMAIN_ZONE_ID"
  type  = "String"
  value = data.aws_route53_zone.janus_dns.zone_id
  overwrite = true
   tags = {
       "env"     = "production"
       "project" = "janus"
   }
}

resource "aws_ssm_parameter" "janus_git_hook_url" {
  name  = "/janus/env/production/GIT_HOOK_URL"
  type  = "String"
  value = join("", ["https://api.",trimsuffix(data.aws_route53_zone.janus_dns.name, "."), "/hook"])
  overwrite = true
   tags = {
       "env"     = "production"
       "project" = "janus"
   }
}
