resource "aws_elasticache_subnet_group" "redis_subnet_group" {
  name       = "example-subnet-group"
  subnet_ids = [aws_subnet.janus-public-subnet.id]
}  

resource "aws_elasticache_cluster" "janus_redis" {
  cluster_id           = "cluster-janus-redis"
  engine               = "redis"
  node_type            = "cache.t2.micro"
  num_cache_nodes      = 1
  parameter_group_name = "default.redis3.2"
  engine_version       = "3.2.10"
  port                 = 6379
  subnet_group_name    = aws_elasticache_subnet_group.redis_subnet_group.name
  security_group_ids   = [aws_security_group.janus_redis_cluster_sg.id]
  apply_immediately    = true
}

output "redis_cluster_endpoint" {
  value = aws_elasticache_cluster.janus_redis.cache_nodes.0.address
}