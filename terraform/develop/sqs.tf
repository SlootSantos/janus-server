resource "aws_sqs_queue" "deadletter_queue" {
  name = "janus-sqs-deadletter"
}

resource "aws_sqs_queue" "AccessID" {
  name           = "janus-access-id-q"
  redrive_policy = "{\"deadLetterTargetArn\":\"${aws_sqs_queue.deadletter_queue.arn}\",\"maxReceiveCount\": 5}"
}

resource "aws_sqs_queue" "DestroyCDN" {
  name           = "janus-destroy-cdn-q"
  redrive_policy = "{\"deadLetterTargetArn\":\"${aws_sqs_queue.deadletter_queue.arn}\",\"maxReceiveCount\": 5}"
}
resource "aws_sqs_queue" "Certificate" {
  name           = "janus-certificate-q"
  redrive_policy = "{\"deadLetterTargetArn\":\"${aws_sqs_queue.deadletter_queue.arn}\",\"maxReceiveCount\": 10}"
}

output "queue_destroyCDN_URL" {
  value = aws_sqs_queue.DestroyCDN.id
}

output "queue_accessID_URL" {
  value = aws_sqs_queue.AccessID.id
}

output "queue_certificate_URL" {
  value = aws_sqs_queue.Certificate.id
}

