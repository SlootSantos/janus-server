//go:generate mockgen -source=$GOFILE -destination=mock.$GOFILE -package=$GOPACKAGE

package queue

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

type QueueMessage map[string]*sqs.MessageAttributeValue

type queueHandler interface {
	SendMessage(*sqs.SendMessageInput) (*sqs.SendMessageOutput, error)
	ReceiveMessage(*sqs.ReceiveMessageInput) (*sqs.ReceiveMessageOutput, error)
	DeleteMessage(*sqs.DeleteMessageInput) (*sqs.DeleteMessageOutput, error)
}

type Queue struct {
	name              string
	url               string
	sqs               queueHandler
	visibilityTimeout int64
	listenerFunc      func(QueueMessage) bool
}

type Q struct {
	AccessID   Queue
	DestroyCDN Queue
}

func New(sess *session.Session) Q {
	s := sqs.New(sess)

	urlDestroyCDN := os.Getenv("SQS_URL_DESTROY_CDN")
	urlAccessID := os.Getenv("SQS_URL_ACCESS_ID")

	q := Q{
		AccessID: Queue{
			name:              "AccessIDQueue",
			url:               urlAccessID,
			visibilityTimeout: 20,
			sqs:               s,
		},
		DestroyCDN: Queue{
			name:              "DestroyCDNQueue",
			visibilityTimeout: 300,
			url:               urlDestroyCDN,
			sqs:               s,
		},
	}

	return q
}

func (q *Queue) Push(message QueueMessage) {
	res, err := q.sqs.SendMessage(&sqs.SendMessageInput{
		DelaySeconds:      aws.Int64(10),
		MessageAttributes: message,
		MessageBody:       aws.String("."),
		QueueUrl:          &q.url,
	})

	log.Println(res, err)
}

func (q *Queue) Pull() {
	if q.listenerFunc == nil {
		return
	}

	log.Println("Initializing Polling for: ", q.name)

	for {
		q.pull()
	}
}

func (q *Queue) pull() {
	result, err := q.sqs.ReceiveMessage(&sqs.ReceiveMessageInput{
		AttributeNames: []*string{
			aws.String(sqs.MessageSystemAttributeNameSentTimestamp),
		},
		MessageAttributeNames: []*string{
			aws.String(sqs.QueueAttributeNameAll),
		},
		QueueUrl:            &q.url,
		MaxNumberOfMessages: aws.Int64(10),
		VisibilityTimeout:   aws.Int64(q.visibilityTimeout),
		WaitTimeSeconds:     aws.Int64(20),
	})

	if err != nil {
		log.Println("Error", err)
		return
	}

	if len(result.Messages) == 0 {
		return
	}

	for _, m := range result.Messages {
		qMessage := QueueMessage(m.MessageAttributes)

		messageAcknowledged := q.listenerFunc(qMessage)
		if !messageAcknowledged {
			return
		}

		_, err := q.sqs.DeleteMessage(&sqs.DeleteMessageInput{
			QueueUrl:      &q.url,
			ReceiptHandle: m.ReceiptHandle,
		})

		if err != nil {
			fmt.Println("Delete Error", err)
			return
		}
	}
}

func (q *Queue) SetListener(listener func(QueueMessage) bool) error {
	if q.listenerFunc != nil {
		return errors.New("Queue: has already registered listener " + q.name)
	}

	q.listenerFunc = listener
	go q.Pull()

	return nil
}

func NewMockQ(sqsMock *MockqueueHandler) Q {
	q := Q{
		AccessID: Queue{
			name:              "AccessIDQueue",
			url:               "url",
			visibilityTimeout: 20,
			sqs:               sqsMock,
		},
		DestroyCDN: Queue{
			name:              "DestroyCDNQueue",
			visibilityTimeout: 300,
			url:               "url",
			sqs:               sqsMock,
		},
	}

	return q
}
