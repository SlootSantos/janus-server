package pipeline

import (
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/SlootSantos/janus-server/pkg/jam"
	"github.com/SlootSantos/janus-server/pkg/storage"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

// Hook is the git hook route identifier
const Hook = "/hook"
const githubActionReleased = "released"

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
		case *github.PullRequestEvent:
			if *e.Action == "opened" || *e.Action == "synchronize" {
				log.Println("GOING FOR PR", *e.Action)

				stack, err := getStackByRepo(*e.Repo.Name, *e.Repo.Owner.Login)
				if err != nil {
					log.Println("Cannot build container for", *e.Repo.FullName, err.Error())
					return
				}

				buildID := generateRandomID()
				stack.Build = &storage.BuildModel{
					Latest: buildID,
				}

				go updateStacks(*e.Repo.Owner.Login, stack)

				Push(ContainerRunParams{
					AWSAccess: os.Getenv("PIPELINE_DEPLOYER_ACCESS"),
					AWSSecret: os.Getenv("PIPELINE_DEPLOYER_SECRET"),
					Branch:    getBranchFromRef("pr"),
					Bucket:    stack.BucketID,
					Repo:      stack.Repo.Name,
					Commit:    *e.PullRequest.Head.SHA,
					CDN:       stack.CDN.ID,
					User:      *e.Repo.Owner.Login,
					Token:     GetUserToken(*e.Sender.Login),
					buildID:   buildID,
					PrID:      strconv.Itoa(*e.Number),
				})
			}

		case *github.ReleaseEvent:
			stack, err := getStackByRepo(*e.Repo.Name, *e.Repo.Owner.Login)
			if err != nil {
				log.Println("Cannot build container for", *e.Repo.FullName, err.Error())
				return
			}

			awsAccess := os.Getenv("PIPELINE_DEPLOYER_ACCESS")
			awsSecrect := os.Getenv("PIPELINE_DEPLOYER_SECRET")

			if stack.IsThirdParty {
				user, err := storage.Store.User.Get(*e.Repo.Owner.Login)
				if err != nil {
					log.Println("Cannot build container for", err.Error())
					return
				}

				awsAccess = user.ThirdPartyAWS.AccessKey
				awsSecrect = user.ThirdPartyAWS.SecretKey
			}

			buildID := generateRandomID()
			stack.Build = &storage.BuildModel{
				Latest: buildID,
			}

			go updateStacks(*e.Repo.Owner.Login, stack)

			if *e.Action == githubActionReleased {
				Push(ContainerRunParams{
					AWSAccess: awsAccess,
					AWSSecret: awsSecrect,
					Branch:    getBranchFromRef(*e.Release.TargetCommitish) + ":tag",
					Bucket:    stack.BucketID,
					Repo:      stack.Repo.Name,
					CDN:       stack.CDN.ID,
					User:      *e.Repo.Owner.Login,
					Token:     GetUserToken(*e.Sender.Login),
					buildID:   buildID,
				})
			}

		case *github.PushEvent:
			if strings.Contains(*e.Ref, "tags") {
				log.Println("Skipping because TAG")
				return
			}

			stack, err := getStackByRepo(*e.Repo.Name, *e.Repo.Owner.Name)
			if err != nil {
				log.Println("Cannot build container for", *e.Repo.FullName, err.Error())
				return
			}

			awsAccess := os.Getenv("PIPELINE_DEPLOYER_ACCESS")
			awsSecrect := os.Getenv("PIPELINE_DEPLOYER_SECRET")

			if stack.IsThirdParty {
				user, err := storage.Store.User.Get(*e.Repo.Owner.Name)
				if err != nil {
					log.Println("Cannot build container for", err.Error())
					return
				}

				awsAccess = user.ThirdPartyAWS.AccessKey
				awsSecrect = user.ThirdPartyAWS.SecretKey
			}

			buildID := generateRandomID()
			stack.Build = &storage.BuildModel{
				Latest: buildID,
			}
			go updateStacks(*e.Repo.Owner.Name, stack)

			Push(ContainerRunParams{
				AWSAccess: awsAccess,
				AWSSecret: awsSecrect,
				Branch:    getBranchFromRef(*e.Ref),
				Commit:    *e.After,
				Bucket:    stack.BucketID,
				Repo:      stack.Repo.Name,
				CDN:       stack.CDN.ID,
				User:      *e.Repo.Owner.Name,
				Token:     GetUserToken(*e.Pusher.Name),
				buildID:   buildID,
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

func updateStacks(user string, newStack *storage.StackModel) error {
	jamList, _ := storage.Store.Stack.Get(user)

	for idx, jam := range jamList {
		log.Println("REPO:", jam.Repo.Name)
		if jam.Repo.Name != newStack.Repo.Name {
			continue
		}

		log.Println("found it!")
		jamList[idx] = *newStack
		break
	}

	storage.Store.Stack.Set(user, jamList)

	return nil
}

func GetUserToken(user string) string {
	u, _ := storage.Store.User.Get(user)

	var token oauth2.Token
	err := json.Unmarshal([]byte(u.Token), &token)
	if err != nil {
		panic("Not working ")
	}

	return token.AccessToken
}

func generateRandomID() string {
	random := make([]byte, 16)
	rand.Read(random)

	return fmt.Sprintf("%x", random)
}

func getBranchFromRef(ref string) string {
	refArr := strings.Split(ref, "/")
	branchName := refArr[len(refArr)-1]

	log.Println("BRANCH:", branchName)
	return branchName
}
