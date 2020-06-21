package storage

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/golang/mock/gomock"
)

type BuildModel struct {
	Latest string `json:"latest"`
}
type Storage struct {
	User  UserIface
	Repo  RepoIface
	Stack StackIface
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
	users:  {1, users},
	stacks: {2, stacks},
	repos:  {3, repos},
}

var Store *Storage

func Init(sess *session.Session) *Storage {
	store := &Storage{
		User:  newUserDB(databasesMap[users], sess),
		Repo:  newRepoDB(databasesMap[repos], sess),
		Stack: newStackDB(databasesMap[stacks], sess),
	}

	Store = store
	return store
}

func MockInit(ctrl *gomock.Controller) (*MockUserIface, *MockRepoIface, *MockStackIface) {
	userMock := NewMockUserIface(ctrl)
	repoMock := NewMockRepoIface(ctrl)
	stackMock := NewMockStackIface(ctrl)

	store := &Storage{
		User:  userMock,
		Repo:  repoMock,
		Stack: stackMock,
	}

	Store = store
	return userMock, repoMock, stackMock
}
