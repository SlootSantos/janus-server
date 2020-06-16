package storage

import (
	"encoding/json"
	"log"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/go-redis/redis/v7"
)

type OrganizationModel struct {
	ThirdPartyAWS ThirdPartyAWS
	Billing       UserBilling
	Name          string `json:"name"`
	ID            string `json:"id"`
	IsPro         bool   `json:"isPro"`
}

type organization struct {
	rC *redis.Client
	dC *dynamoClient
}

const dynamoOrgaPrimaryKey = "user"

func newOrgaDB(db db, s *session.Session) *organization {
	return &organization{
		rC: connectRedis(db.redisID),
		dC: connectDynamo(s, db.dynamoTable, dynamoOrgaPrimaryKey),
	}
}

func (o *organization) Get(key string) (*OrganizationModel, error) {
	result, err := o.rC.Get(key).Result()
	if err == redis.Nil {
		log.Println("Initalizing empty redis 'Organizatoin' DB")
		go o.Set(key, &OrganizationModel{})

		return &OrganizationModel{}, err
	}

	model := &OrganizationModel{}
	err = json.Unmarshal([]byte(result), model)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return model, nil
}

func (o *organization) Set(key string, value *OrganizationModel) error {
	jsonValue, err := json.Marshal(value)
	if err != nil {
		log.Println("cannot marshal user", err)
		return err
	}

	_, err = o.rC.Set(key, jsonValue, 0).Result()
	return err
}

func (o *organization) getFromDynamo(key string) (*OrganizationModel, error) {
	dbResult := o.dC.GetItem(key)

	model := &OrganizationModel{}
	err := dynamodbattribute.UnmarshalMap(dbResult.Item, &model)
	if err != nil {
		log.Println("could not get dynamo", err)
		return &OrganizationModel{}, err
	}

	go o.Set(key, model)
	// dynamoValue := dynamoStack{User: key, Stacks: value}
	// err = s.dC.SetItem(dynamoValue)
	// if err != nil {
	// 	log.Println("cloud not set to dynamo", err)
	// }

	return model, nil
}
