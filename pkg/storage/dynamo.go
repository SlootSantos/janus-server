package storage

import (
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type dynamoClient struct {
	dynamo     *dynamodb.DynamoDB
	table      string
	primaryKey string
}

type dynamoItem struct{}

func (d *dynamoClient) SetItem(item interface{}) error {
	av, err := dynamodbattribute.MarshalMap(item)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String(d.table),
		Item:      av,
	}

	_, err = d.dynamo.PutItem(input)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	return nil
}

func (d *dynamoClient) GetItem(selector string) *dynamodb.GetItemOutput {
	result, err := d.dynamo.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(d.table),
		Key: map[string]*dynamodb.AttributeValue{
			d.primaryKey: {
				S: aws.String(selector),
			},
		},
	})

	if err != nil {
		log.Println("FAILED GETTING", err)
	}

	return result
}

func connectDynamo(sess *session.Session, table string, primaryKey string) *dynamoClient {
	dynamo := dynamodb.New(sess)

	return &dynamoClient{dynamo, table, primaryKey}
}
