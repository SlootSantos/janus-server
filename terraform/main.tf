terraform {
  backend "s3" {
    bucket = "janus-infrastructre-state"
    region = "us-east-1"
    key    = "state"
  }
}

provider "aws" {
  region = "us-east-1"
}

