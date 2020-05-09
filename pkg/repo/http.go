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
	"github.com/google/go-github/github"
)

// RoutePrefix is the REST endpoint for repo
const RoutePrefix = "/repo"

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

func handleGET(w http.ResponseWriter, req *http.Request) {
	repos := list(req.Context())

	w.Write(repos)
}
func handlePOST(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "POST")
}

func list(ctx context.Context) []byte {
	reposStorage, err := storage.Store.Repo.Get(ctx.Value(auth.ContextKeyUserName).(string))
	if err != nil {
		log.Println(err)
	}

	if reposStorage == "" {
		client := auth.AuthenticateUser(ctx.Value(auth.ContextKeyToken).(string))

		lsOpt := &github.RepositoryListOptions{
			ListOptions: github.ListOptions{
				PerPage: 50,
			},
		}

		reps, _, err := client.Repositories.List(ctx, "", lsOpt)
		if err != nil {
			fmt.Println(err)
		}

		repoList := []string{}
		for _, r := range reps {
			repoList = append(repoList, *r.Name)
		}

		reposJSON, _ := json.Marshal(repoList)
		err = storage.Store.Repo.Set(ctx.Value(auth.ContextKeyUserName).(string), reposJSON)

		return reposJSON
	}

	return []byte(reposStorage)
}
