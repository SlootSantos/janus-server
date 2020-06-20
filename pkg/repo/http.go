package repo

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/SlootSantos/janus-server/pkg/api/auth"
	"github.com/SlootSantos/janus-server/pkg/storage"
	"github.com/go-redis/redis/v7"
	"github.com/google/go-github/github"
)

// RoutePrefix is the REST endpoint for repo
const RoutePrefix = "/repo"

// RouteSyncPrefix is the REST endpoint for syncing repos
const RouteSyncPrefix = "/repo/sync"

var methodHandlerMap = map[string]http.HandlerFunc{
	http.MethodGet:  handleGET,
	http.MethodPost: handlePOST,
}

// HandleHTTP handles the Repo http endpoint
func HandleHTTP(w http.ResponseWriter, req *http.Request) {
	if handler, ok := methodHandlerMap[req.Method]; ok {
		handler(w, req)
		return
	}

	w.WriteHeader(http.StatusMethodNotAllowed)
	w.Write([]byte("Method not allowed"))
}

// HandleSyncHTTP handles the Repo http endpoint
func HandleSyncHTTP(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Method not allowed"))
		return
	}

	repos := sync(req.Context())
	w.Write(repos)
}

func handleGET(w http.ResponseWriter, req *http.Request) {
	repos := list(req.Context())

	w.Write(repos)
}
func handlePOST(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "POST")
}

func list(ctx context.Context) []byte {
	reposStorage, err := storage.Store.Repo.Get(ctx.Value(auth.ContextKeyUserName).(string))

	if err == redis.Nil {
		log.Println("Empty Repo Redis Instance")
		reposStorage = ""
	}

	if reposStorage == "" {
		reposJSON := fetchReposJSON(ctx)
		_, err = storage.Store.Repo.Set(ctx.Value(auth.ContextKeyUserName).(string), reposJSON)
		if err != nil {
			log.Println(err)
		}

		return reposJSON
	}

	return []byte(reposStorage)
}

func sync(ctx context.Context) []byte {
	reposJSON := fetchReposJSON(ctx)
	_, err := storage.Store.Repo.Set(ctx.Value(auth.ContextKeyUserName).(string), reposJSON)
	if err != nil {
		log.Println(err)
	}

	return reposJSON
}

func fetchReposJSON(ctx context.Context) []byte {
	client := auth.AuthenticateUser(ctx.Value(auth.ContextKeyToken).(string))

	lsOpt := &github.RepositoryListOptions{
		// Affiliation: "oranization_member",
		ListOptions: github.ListOptions{
			PerPage: 100,
		},
	}

	reps, _, err := client.Repositories.List(ctx, "", lsOpt)
	if err != nil {
		fmt.Println(err)
	}

	repoList := []*storage.RepoModel{}
	for _, r := range reps {
		newRepo := &storage.RepoModel{
			Name:  *r.Name,
			Owner: *r.Owner.Login,
			Type:  *r.Owner.Type,
		}

		repoList = append(repoList, newRepo)
	}

	reposJSON, _ := json.Marshal(repoList)
	return reposJSON
}
