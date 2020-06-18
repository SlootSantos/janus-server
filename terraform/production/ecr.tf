resource "aws_ecr_repository" "janus_server" {
  name                 = "janus/server"
  image_tag_mutability = "MUTABLE"
}

output "ecr_repo_url__janus_server" {
  value = aws_ecr_repository.janus_server.repository_url
}

