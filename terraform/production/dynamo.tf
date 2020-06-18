resource "aws_dynamodb_table" "janus_dynamo_users" {
  name           = "Users"
  billing_mode   = "PROVISIONED"
  read_capacity  = 2
  write_capacity = 2
  hash_key       = "user"

  attribute {
    name = "user"
    type = "S"
  }

  ttl {
    attribute_name = "ttl"
    enabled        = true
  }

  tags = {
    Name        = "janus-dynamo-users"
    Application = "janus"
  }
}

resource "aws_dynamodb_table" "janus_dynamo_stacks" {
  name           = "Stacks"
  billing_mode   = "PROVISIONED"
  read_capacity  = 2
  write_capacity = 2
  hash_key       = "user"

  attribute {
    name = "user"
    type = "S"
  }

  ttl {
    attribute_name = "ttl"
    enabled        = true
  }

  tags = {
    Name        = "janus-dynamo-stacks"
    Application = "janus"
  }
}

output "dynamo_users" {
  value = aws_dynamodb_table.janus_dynamo_users.arn
}

output "dynamo_stacks" {
  value = aws_dynamodb_table.janus_dynamo_stacks.arn
}
