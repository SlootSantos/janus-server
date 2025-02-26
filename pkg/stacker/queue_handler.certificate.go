package stacker

import (
	"log"
	"os"
	"strconv"

	"github.com/SlootSantos/janus-server/pkg/cdn"
	"github.com/SlootSantos/janus-server/pkg/queue"
	"github.com/SlootSantos/janus-server/pkg/session"
	"github.com/SlootSantos/janus-server/pkg/storage"
)

func updateCDNCertificate(message queue.QueueMessage) (ack bool) {
	distroID, ok := message[queue.MessageCertificateDistroID]
	if !ok {
		log.Println("Handling queue message for CDN failed. attribute:", queue.MessageCertificateDistroID, " does not exist on message")
		return ack
	}

	certificateARN, ok := message[queue.MessageCertificateARN]
	if !ok {
		log.Println("Handling queue message for CDN failed. attribute:", queue.MessageCertificateARN, " does not exist on message")
		return ack
	}

	subdomain, ok := message[queue.MessageCertificateSubDomain]
	if !ok {
		log.Println("Handling queue message for CDN failed. attribute:", queue.MessageCertificateSubDomain, " does not exist on message")
		return ack
	}

	username, ok := message[queue.MessageCommonUser]
	if !ok {
		log.Println("Handling queue message for Bucket-Policy failed. attribute:", queue.MessageCommonUser, " does not exist on message")
		return ack
	}

	isThirdPartyStr, ok := message[queue.MessageCommonIsThirdParty]
	if !ok {
		log.Println("Handling queue message for CDN failed. attribute:", queue.MessageCommonIsThirdParty, " does not exist on message")
		return ack
	}

	isThridParty, err := strconv.ParseBool(*isThirdPartyStr.StringValue)
	if err != nil {
		log.Printf("Error parsing isThirdParty to bool failed. org value: %s", *isThirdPartyStr.StringValue)
		isThridParty = false
	}

	awsSess, _ := session.AWSSession()
	cdnConf := &cdn.CreateCDNParams{
		HostedZoneID: os.Getenv("DOMAIN_ZONE_ID"),
		CertARN:      os.Getenv("DOMAIN_CERT_ARN"),
		Session:      awsSess,
		Domain:       os.Getenv("DOMAIN_HOST"),
	}

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

		cdnConf = &cdn.CreateCDNParams{
			Domain:       user.ThirdPartyAWS.Domain,
			HostedZoneID: user.ThirdPartyAWS.HostedZoneID,
			Session:      awsSess,
		}
	}
	c := cdn.New(cdnConf)

	return c.HandleQueueMessageCertificate(*distroID.StringValue, *certificateARN.StringValue, *subdomain.StringValue)
}
