package stacker

import (
	"log"
	"strconv"

	"github.com/SlootSantos/janus-server/pkg/bucket"
	"github.com/SlootSantos/janus-server/pkg/queue"
	"github.com/SlootSantos/janus-server/pkg/session"
	"github.com/SlootSantos/janus-server/pkg/storage"
)

func setAccessID(message queue.QueueMessage) (ack bool) {
	log.Println("SETTING ACCESS ID!")
	bucketID, ok := message[queue.MessageDestroyBucketID]
	if !ok {
		log.Println("Handling queue message for Bucket-Policy failed. attribute:", queue.MessageDestroyBucketID, " does not exist on message")
		return ack
	}

	accessID, ok := message[queue.MessageDestroyAccessID]
	if !ok {
		log.Println("Handling queue message for Bucket-Policy failed. attribute:", queue.MessageDestroyAccessID, " does not exist on message")
		return ack
	}

	username, ok := message[queue.MessageCommonUser]
	if !ok {
		log.Println("Handling queue message for Bucket-Policy failed. attribute:", queue.MessageCommonUser, " does not exist on message")
		return ack
	}

	isThirdPartyStr, ok := message[queue.MessageCommonIsThirdParty]
	if !ok {
		log.Println("Handling queue message for Bucket-Policy failed. attribute:", queue.MessageCommonIsThirdParty, " does not exist on message")
		return ack
	}

	isThridParty, err := strconv.ParseBool(*isThirdPartyStr.StringValue)
	if err != nil {
		log.Printf("Error parsing isThirdParty to bool failed. org value: %s", *isThirdPartyStr.StringValue)
		isThridParty = false
	}

	awsSess, _ := session.AWSSession()

	if isThridParty {
		user, err := storage.Store.User.Get(*username.StringValue)
		if err != nil {
			return false
		}

		awsSess, err = session.AWSSessionThirdParty(user.ThirdPartyAWS.AccessKey, user.ThirdPartyAWS.SecretKey)
		if err != nil {
			log.Println("Could not create Session--- access")
			return
		}
	}

	bucket := bucket.New(awsSess, nil)

	return bucket.HandleQueueMessageAccessID(*bucketID.StringValue, *accessID.StringValue)
}
