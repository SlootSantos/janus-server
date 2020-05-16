resource "aws_iam_instance_profile" "janus_beanstalk_profile" {
  name = "test_profile"
  role = aws_iam_role.janus_beanstalk_role.name
}

resource "aws_iam_role" "janus_beanstalk_role" {
  name = "janus-beanstalk-role"
assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Principal": {
        "Service": "ec2.amazonaws.com"
      },
      "Effect": "Allow",
      "Sid": ""
    }
  ]
}
EOF
}


data "aws_iam_policy" "ECR_READ" {
  arn = "arn:aws:iam::aws:policy/AmazonEC2ContainerRegistryReadOnly"
}

resource "aws_iam_role_policy_attachment" "beanstalk_role_policy_attach_ECR_READ" {
  role       = aws_iam_role.janus_beanstalk_role.name
  policy_arn = data.aws_iam_policy.ECR_READ.arn
}

data "aws_iam_policy" "BEANSTALK_WEB" {
  arn = "arn:aws:iam::aws:policy/AWSElasticBeanstalkWebTier"
}

resource "aws_iam_role_policy_attachment" "beanstalk_role_policy_attach_BEANSTALK_WEB" {
  role       = aws_iam_role.janus_beanstalk_role.name
  policy_arn = data.aws_iam_policy.BEANSTALK_WEB.arn
}

data "aws_iam_policy" "BEANSTALK_WORK" {
  arn = "arn:aws:iam::aws:policy/AWSElasticBeanstalkWorkerTier"
}

resource "aws_iam_role_policy_attachment" "beanstalk_role_policy_attach_BEANSTALK_WORK" {
  role       = aws_iam_role.janus_beanstalk_role.name
  policy_arn = data.aws_iam_policy.BEANSTALK_WORK.arn
}

data "aws_iam_policy" "BEANSTALK_MULTI_DOCKER" {
  arn = "arn:aws:iam::aws:policy/AWSElasticBeanstalkMulticontainerDocker"
}

resource "aws_iam_role_policy_attachment" "beanstalk_role_policy_attach_BEANSTALK_MULTI_DOCKER" {
  role       = aws_iam_role.janus_beanstalk_role.name
  policy_arn = data.aws_iam_policy.BEANSTALK_MULTI_DOCKER.arn
}