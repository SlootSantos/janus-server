package storage

import (
	"github.com/aws/aws-sdk-go/aws/session"
)

type BuildModel struct {
	Latest string `json:"latest"`
}
type Storage struct {
	User         *user
	Repo         *repo
	Stack        *stack
	Organization *organization
}

type db struct {
	redisID     int
	dynamoTable string
}

const (
	users         = "Users"
	repos         = "Repos"
	stacks        = "Stacks"
	organizations = "Organizations"
)

var databasesMap = map[string]db{
	users:         {1, users},
	stacks:        {2, stacks},
	repos:         {3, repos},
	organizations: {4, organizations},
}

var Store *Storage

func Init(sess *session.Session) *Storage {
	store := &Storage{
		User:         newUserDB(databasesMap[users], sess),
		Repo:         newRepoDB(databasesMap[repos], sess),
		Stack:        newStackDB(databasesMap[stacks], sess),
		Organization: newOrgaDB(databasesMap[organizations], sess),
	}

	Store = store
	return store
}
