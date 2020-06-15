package repo

import (
	"context"
	"fmt"
	"log"
	"os"

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

func (r *Repo) Create(ctx context.Context, params *jam.CreationParam, out *jam.OutputParam) (string, error) {
	hookURL := os.Getenv("GIT_HOOK_URL")
	repoName := params.Repo.Name
	repoOwner := params.Repo.Owner

	client := auth.AuthenticateUser(ctx.Value(auth.ContextKeyToken).(string))

	config := hookConfig{
		"name":         &hookName,
		"url":          &hookURL,
		"content_type": "json",
	}

	gitHook := &github.Hook{
		Name:   &hookName,
		URL:    &hookURL,
		Config: config,
		Events: []string{
			"push",
			"pull_request",
			"release",
		},
	}

	hook, _, err := client.Repositories.CreateHook(ctx, repoOwner, repoName, gitHook)
	if err != nil {
		return "broken", err
	}

	out.Repo = &params.Repo
	out.Repo.ID = *hook.ID

	log.Println("DONE: creating git hook", *hook.ID)
	return "nil", nil
}

func (r *Repo) Destroy(ctx context.Context, params *jam.DeletionParam) error {
	fmt.Println("START: destroying git hook", params.Repo.Name)

	repoName := params.Repo.Name
	repoOwner := params.Repo.Owner

	client := auth.AuthenticateUser(ctx.Value(auth.ContextKeyToken).(string))

	// backward compatible: if no owner on repo => fetch user from token
	if repoOwner == "" {
		user, _, err := client.Users.Get(context.Background(), "")
		if err != nil {
			fmt.Println(err)
			return err
		}

		repoOwner = *user.Login
	}

	_, err := client.Repositories.DeleteHook(context.Background(), repoOwner, repoName, params.Repo.ID)
	if err != nil {
		return err
	}

	log.Println("DONE: destroying git hook")
	return nil
}

func (r *Repo) List(ctx context.Context) string {
	return "nil"
}
