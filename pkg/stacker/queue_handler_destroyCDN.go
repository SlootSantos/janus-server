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

func deleteDisabledDistro(message queue.QueueMessage) (ack bool) {
	log.Println("RECEIVING DELETION MESSAGE")
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

	certARN, ok := message[queue.MessageCertificateARN]
	if !ok {
		log.Println("Handling queue message for CDN failed. attribute:", queue.MessageCertificateARN, " does not exist on message")
		return ack
	}

	username, ok := message[queue.MessageCommonUser]
	if !ok {
		log.Println("Handling queue message for CDN failed. attribute:", queue.MessageCommonUser, " does not exist on message")
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
			Domain:  user.ThirdPartyAWS.Domain,
			CertARN: os.Getenv("DOMAIN_CERT_ARN"),
			// HostedZoneID: os.Getenv("DOMAIN_ZONE_ID"),
			HostedZoneID: "/hostedzone/Z0475209DXX3SXD7W99L",
			Session:      awsSess,
		}
	}

	c := cdn.New(cdnConf)

	return c.HandleQueueMessaeDestroyCDN(*distroID.StringValue, *etag.StringValue, *certARN.StringValue)
}
