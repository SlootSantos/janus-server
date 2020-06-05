package storage

import "github.com/aws/aws-sdk-go/aws/session"

type BuildModel struct {
	Latest string `json:"latest"`
}
type Storage struct {
	Stack *stack
	User  *user
	Repo  *repo
}

type db struct {
	redisID     int
	dynamoTable string
}

const stacks = "Stacks"
const users = "Users"
const repos = "Repos"

var databasesMap = map[string]db{
	users:  {1, users},
	stacks: {2, stacks},
	repos:  {3, repos},
}

var Store *Storage

func Init(sess *session.Session) *Storage {
	store := &Storage{
		Stack: newStackDB(databasesMap[stacks], sess),
		User:  newUserDB(databasesMap[users], sess),
		Repo:  newRepoDB(databasesMap[repos], sess),
	}

	Store = store
	return store
}
