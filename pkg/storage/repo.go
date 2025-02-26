package storage

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/go-redis/redis/v7"
)

// RepoModel represents a Stacks repository
type RepoModel struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Owner string `json:"owner"`
	Type  string `json:"type"`
}

type repo struct {
	rC *redis.Client
}

func newRepoDB(db db, s *session.Session) *repo {
	return &repo{
		rC: connectRedis(db.redisID),
	}
}

func (r *repo) Get(key string) (string, error) {
	return r.rC.Get(key).Result()
}

func (r *repo) Set(key string, value []byte) (string, error) {
	return r.rC.Set(key, value, 0).Result()
}
