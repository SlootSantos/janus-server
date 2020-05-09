package storage

import (
	"encoding/json"
	"log"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/go-redis/redis/v7"
)

// StackModel represents an entire JAM Stack
type StackModel struct {
	ID       string `json:"id"`
	BucketID string `json:"bucketId"`
	CDN      *StackCDNModel
	Repo     *RepoModel
}

// StackCDNModel contains all stack relevant information about the CDN
type StackCDNModel struct {
	ID       string `json:"id"`
	Domain   string `json:"domain"`
	AccessID string `json:"accessId"`
}

type stack struct {
	rC *redis.Client
	dC *dynamoClient
}

type dynamoStack struct {
	User   string `json:"user"`
	Stacks []StackModel
}

const dynamoStackPrimaryKey = "user"

func newStackDB(db db, s *session.Session) *stack {
	return &stack{
		rC: connectRedis(db.redisID),
		dC: connectDynamo(s, db.dynamoTable, dynamoStackPrimaryKey),
	}
}

func (s *stack) Get(key string) ([]StackModel, error) {
	result, err := s.rC.Get(key).Result()
	if err != nil {
		return nil, err
	}

	if result == "" {
		return s.getFromDynamo(key)
	}

	model := []StackModel{}
	err = json.Unmarshal([]byte(result), &model)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return model, nil
}

func (s *stack) getFromDynamo(key string) ([]StackModel, error) {
	dbResult := s.dC.GetItem(key)

	model := dynamoStack{}
	err := dynamodbattribute.UnmarshalMap(dbResult.Item, &model)
	if err != nil {
		log.Println("could not get dynamo", err)
		return []StackModel{}, err
	}

	if model.Stacks == nil {
		return []StackModel{}, nil
	}

	go s.Set(key, model.Stacks)

	return model.Stacks, nil
}

func (s *stack) Set(key string, value []StackModel) error {
	jsonValue, err := json.Marshal(value)
	if err != nil {
		log.Println("cannot marshal user", err)
		return err
	}

	err = s.rC.Set(key, jsonValue, 0).Err()
	if err != nil {
		log.Println("failed to store user in redis", err)
		return err
	}

	dynamoValue := dynamoStack{User: key, Stacks: value}
	err = s.dC.SetItem(dynamoValue)
	if err != nil {
		log.Println("cloud not set to dynamo", err)
	}

	return nil
}
