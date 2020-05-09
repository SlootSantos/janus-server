package pipeline

import (
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/SlootSantos/janus-server/pkg/jam"
	"github.com/SlootSantos/janus-server/pkg/storage"
	"github.com/google/go-github/github"
)

type repository struct {
	Fullname string `json:"full_name"`
}

type webhook struct {
	Repository repository
}

// Hook is the git hook route identifier
const Hook = "/hook"

// HandleHook handles incoming Github hooks
func HandleHook(w http.ResponseWriter, req *http.Request) {
	payload, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Printf("error reading request body: err=%s\n", err)
		return
	}

	defer req.Body.Close()

	go func() {
		event, err := github.ParseWebHook(github.WebHookType(req), payload)
		if err != nil {
			log.Printf("could not parse webhook: err=%s\n", err)
			return
		}

		switch e := event.(type) {
		case *github.PushEvent:
			stack, err := getStackByRepo(*e.Repo.Name, *e.Repo.Owner.Name)
			if err != nil {
				log.Println("Cannot build container for", *e.Repo.FullName, err.Error())
				return
			}

			BuildRepoSources(ContainerRunParams{
				AWSAccess: os.Getenv("PIPELINE_DEPLOYER_ACCESS"),
				AWSSecret: os.Getenv("PIPELINE_DEPLOYER_SECRET"),
				Bucket:    stack.BucketID,
				Repo:      stack.Repo.Name,
				CDN:       stack.CDN.ID,
				User:      *e.Repo.Owner.Name,
			})

		default:
			log.Printf("unknown event type %s\n", github.WebHookType(req))
			return
		}
	}()

	w.WriteHeader(200)
}

func getStackByRepo(repoName string, user string) (*jam.Stack, error) {
	jamList, _ := storage.Store.Stack.Get(user)

	for _, jam := range jamList {
		if jam.Repo.Name != repoName {
			continue
		}

		return &jam, nil
	}

	return nil, errors.New("can not find stack for repo" + repoName)
}
