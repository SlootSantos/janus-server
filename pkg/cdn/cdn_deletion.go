package cdn

import (
	"log"

	"github.com/SlootSantos/janus-server/pkg/queue"
	"github.com/aws/aws-sdk-go/service/cloudfront"
)

func (c *CDN) deleteDisabledDistro(message queue.QueueMessage) (ack bool) {
	distroID, ok := message[queue.MessageAccessDistroID]
	if !ok {
		log.Println("Handling queue message for CDN failed. attribute:", queue.MessageAccessDistroID, " does not exist on message")
		return ack
	}

	etag, ok := message[queue.MessageAccessEtag]
	if !ok {
		log.Println("Handling queue message for CDN failed. attribute:", queue.MessageAccessEtag, " does not exist on message")
		return ack
	}

	deleteDistroInput := &cloudfront.DeleteDistributionInput{
		Id:      distroID.StringValue,
		IfMatch: etag.StringValue,
	}

	_, err := c.cdn.DeleteDistribution(deleteDistroInput)
	if err != nil {
		return ack
	}

	ack = true
	return ack
}
