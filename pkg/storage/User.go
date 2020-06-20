package storage

import (
	"encoding/json"
	"log"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/go-redis/redis/v7"
)

type UserModel struct {
	User          string `json:"user"`
	Token         string `json:"token"`
	IsPro         bool   `json:"isPro"`
	Type          string `json:"type"`
	Billing       *UserBilling
	ThirdPartyAWS *ThirdPartyAWS
}

type AllowedUserSettings struct {
	IsPro         bool   `json:"isPro"`
	Type          string `json:"type"`
	Name          string `json:"name"`
	ThirdPartyAWS *ThirdPartyAWS
}

type AllowedOrgaSettings struct {
	UserMemberStatus string `json:"userMemberStatus"`
	AllowedUserSettings
}

type ThirdPartyAWS struct {
	AccessKey    string `json:"accessKey"`
	SecretKey    string `json:"secretKey"`
	Domain       string `json:"domain"`
	LambdaARN    string `json:"lambdaARN"`
	HostedZoneID string `json:"hostedZoneId"`
}

type UserBilling struct {
	SubscriptionID string `json:"subscriptionId"`
}

type user struct {
	rC *redis.Client
	dC *dynamoClient
}

const dynamoUserPrimaryKey = "user"
const (
	TypeOrganization = "Organization"
	TypeUser         = "User"
)

func newUserDB(db db, s *session.Session) *user {
	return &user{
		rC: connectRedis(db.redisID),
		dC: connectDynamo(s, db.dynamoTable, dynamoUserPrimaryKey),
	}
}

func (u *user) Get(key string) (*UserModel, error) {
	result, err := u.rC.Get(key).Result()
	if err == redis.Nil {
		log.Println("Initalizing empty redis 'User' DB")
		go u.Set(key, &UserModel{})

		return &UserModel{}, err
	}

	if result == "" {
		return u.getFromDynamo(key)
	}

	model := &UserModel{}
	err = json.Unmarshal([]byte(result), model)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return model, nil
}

func (u *user) getFromDynamo(key string) (*UserModel, error) {
	dbResult := u.dC.GetItem(key)

	model := UserModel{}
	err := dynamodbattribute.UnmarshalMap(dbResult.Item, &model)
	if err != nil {
		log.Println("could not get dynamo", err)
		return &UserModel{}, err
	}

	if model.User == "" {
		return &UserModel{}, nil
	}

	return &model, nil
}

func (u *user) Set(key string, value *UserModel) error {
	jsonValue, err := json.Marshal(value)
	if err != nil {
		log.Println("cannot marshal user", err)
		return err
	}

	err = u.rC.Set(key, jsonValue, 0).Err()
	if err != nil {
		log.Println("failed to store user in redis", err)
		return err
	}

	err = u.dC.SetItem(value)
	if err != nil {
		log.Println("cloud not set to dynamo", err)
	}

	return nil
}
