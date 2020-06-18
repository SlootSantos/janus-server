resource "aws_vpc" "janus-vpc" {
  cidr_block           = "10.0.0.0/16"
  enable_dns_hostnames = true
}

# Define the public subnet
resource "aws_subnet" "janus-public-subnet" {
  vpc_id            = aws_vpc.janus-vpc.id
  cidr_block        = "10.0.1.0/24"
  availability_zone = "us-east-1a"

  tags = {
    Name        = "janus-public-subnet"
    Application = "janus"
  }
}

# Define the internet gateway
resource "aws_internet_gateway" "janus-gw" {
  vpc_id = aws_vpc.janus-vpc.id

  tags = {
    Name        = "janus-internet-gateway"
    Application = "janus"
  }
}

# Define the route table
resource "aws_route_table" "web-public-rt" {
  vpc_id = aws_vpc.janus-vpc.id

  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.janus-gw.id
  }

  tags = {
    Name        = "janus-public-routing-table"
    Application = "janus"
  }
}

# Assign the route table to the public Subnet
resource "aws_route_table_association" "web-public-rt" {
  subnet_id      = aws_subnet.janus-public-subnet.id
  route_table_id = aws_route_table.web-public-rt.id
}

# Define the security group for public subnet
resource "aws_security_group" "janus_redis_cluster_sg" {
  name        = "janus_redis_cluster_sg"
  description = "Allow Redis access from same subnet"

  ingress {
    from_port   = 6379
    to_port     = 6379
    protocol    = "tcp"
    cidr_blocks = [aws_subnet.janus-public-subnet.cidr_block]
  }

  vpc_id = aws_vpc.janus-vpc.id

  tags = {
    project = "janus"
  }
}

# Define the security group for public subnet
resource "aws_security_group" "janus_public_sg" {
  name        = "vpc_test_web"
  description = "Allow incoming HTTP connections & SSH access"

  ingress {
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = ["10.0.1.0/24"]
  }

  ingress {
    from_port   = 443
    to_port     = 443
    protocol    = "tcp"
    cidr_blocks = ["10.0.1.0/24"]
  }

  ingress {
    from_port   = -1
    to_port     = -1
    protocol    = "icmp"
    cidr_blocks = ["10.0.1.0/24"]
  }

  ingress {
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = ["10.0.1.0/24"]
  }

  vpc_id = aws_vpc.janus-vpc.id

  tags = {
    Name        = "janus-public-security-group"
    Application = "janus"
  }
}

