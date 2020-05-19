package repo

import (
	"context"
	"fmt"
	"log"

	"github.com/SlootSantos/janus-server/pkg/api/auth"
	"github.com/SlootSantos/janus-server/pkg/jam"
	"github.com/google/go-github/github"
)

// Repo contains all data to interact w/ AWS Cloudfront
type Repo struct{}

// New creates a new CDN creator
func New() *Repo {
	log.Print("DONE: setting up Repo-Creator")

	return &Repo{}
}

type hookConfig map[string]interface{}

var hookName = "web"
var hookURL = "https://312c7c53.ngrok.io/hook" // FROM ENV?

func (r *Repo) Create(ctx context.Context, params *jam.CreationParam, out *jam.OutputParam) (string, error) {
	repoName := params.Repo.Name
	client := auth.AuthenticateUser(ctx.Value(auth.ContextKeyToken).(string))

	user, _, err := client.Users.Get(ctx, "")
	if err != nil {
		fmt.Println(err)
	}

	config := hookConfig{
		"name":         &hookName,
		"url":          &hookURL,
		"content_type": "json",
	}

	gitHook := &github.Hook{
		Name:   &hookName,
		URL:    &hookURL,
		Config: config,
	}

	hook, _, err := client.Repositories.CreateHook(ctx, *user.Login, repoName, gitHook)
	if err != nil {
		return "broken", err
	}

	out.Repo = &jam.StackRepo{
		ID:   *hook.ID,
		Name: repoName,
	}

	log.Println("DONE: creating git hook", *hook.ID)
	return "nil", nil
}

func (r *Repo) Destroy(ctx context.Context, params *jam.DeletionParam) error {
	fmt.Println("START: destroying git hook", params.Repo.Name)
	repoName := params.Repo.Name
	client := auth.AuthenticateUser(ctx.Value(auth.ContextKeyToken).(string))

	user, _, err := client.Users.Get(context.Background(), "")
	if err != nil {
		fmt.Println(err)
		return err
	}

	_, err = client.Repositories.DeleteHook(context.Background(), *user.Login, repoName, params.Repo.ID)
	if err != nil {
		return err
	}

	log.Println("DONE: destroying git hook")
	return nil
}
func (r *Repo) List(ctx context.Context) string {
	return "nil"
}
