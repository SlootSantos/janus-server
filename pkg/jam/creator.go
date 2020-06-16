package jam

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"sync"

	"github.com/SlootSantos/janus-server/pkg/api/auth"
	"github.com/SlootSantos/janus-server/pkg/storage"
	"github.com/google/go-github/github"
)

const bucketPrefix = "janus-bucket-"

// builds entire JAM-Stack
// including a Bucket, a CDN and a Githook
func (j *Creator) Build(ctx context.Context, config StackCreateConfig) ([]byte, error) {
	stackID := generateRandomID()
	subdomain := config.CustomSubDomain

	if subdomain == "" {
		subdomain = stackID
	}

	creationParam := &CreationParam{
		ID: stackID,
		CDN: StackCDN{
			Subdomain: subdomain,
		},
		Bucket: struct{ ID string }{
			ID: bucketPrefix + stackID,
		},
		Repo: config.Repository,
	}

	creationResult, err := j.buildStack(ctx, creationParam)
	creationStack := Stack(*creationResult)
	if err != nil {
		return nil, err
	}

	creationStack.IsThirdParty = config.IsThirdParty
	updatedList, err := storeStack(ctx, &creationStack, config.Repository.Owner)
	if err != nil {
		return []byte{}, err
	}

	return updatedList, nil
}

func (j *Creator) Delete(ctx context.Context, stackConf StackDestroyConfig) ([]byte, error) {
	owner := stackConf.Repository.Owner

	stack := getStackByID(ctx, owner, stackConf.ID)
	if stack == (Stack{}) {
		err := errors.New("No Stack ID: " + stackConf.ID + " for User: " + owner)
		log.Println(err.Error())

		return []byte{}, err
	}

	deletionParam := DeletionParam(stack)
	j.destroyStack(ctx, &deletionParam)

	updatedList, err := removeStack(ctx, &stack, owner)
	if err != nil {
		log.Println(err)
		return []byte{}, err
	}

	return updatedList, nil
}

func (j *Creator) List(ctx context.Context) ([]Stack, error) {
	user := ctx.Value(auth.ContextKeyUserName).(string)
	stacks := getAllStacks(ctx, user, true)

	return stacks, nil
}

func (j *Creator) buildStack(ctx context.Context, creationParam *CreationParam) (*OutputParam, error) {
	out := &OutputParam{ID: creationParam.ID, Repo: &StackRepo{}}

	var wg sync.WaitGroup
	for _, resource := range j.resources {
		wg.Add(1)
		go execCreate(ctx, resource, creationParam, out, &wg)
	}

	wg.Wait()

	return out, nil
}

func (j *Creator) destroyStack(ctx context.Context, deletionParam *DeletionParam) {
	var wg sync.WaitGroup
	for _, resource := range j.resources {
		wg.Add(1)
		go execDestroy(ctx, resource, deletionParam, &wg)
	}

	wg.Wait()

	return
}

func execCreate(ctx context.Context, sr stackResource, creationParam *CreationParam, out *OutputParam, wg *sync.WaitGroup) error {
	defer wg.Done()

	_, err := sr.Create(ctx, creationParam, out)
	if err != nil {
		log.Println("Error creating stackresource", err.Error())
		return err
	}

	return nil
}

func execDestroy(ctx context.Context, sr stackResource, deletionParam *DeletionParam, wg *sync.WaitGroup) error {
	defer wg.Done()

	err := sr.Destroy(ctx, deletionParam)
	if err != nil {
		log.Println("Error destroying stackresource", err.Error())
		return err
	}

	return nil
}

func storeStack(ctx context.Context, stack *Stack, user string) (updatedList []byte, err error) {
	stackList := getAllStacks(ctx, user, false)
	updatedStackList := append(stackList, *stack)

	if err := storage.Store.Stack.Set(user, updatedStackList); err != nil {
		log.Println(err)
		return updatedList, err
	}

	updatedList, err = json.Marshal(updatedStackList)
	if err != nil {
		return updatedList, err
	}

	return updatedList, nil
}

func removeStack(ctx context.Context, stack *Stack, user string) (updatedList []byte, err error) {
	stackList := getAllStacks(ctx, user, false)
	if len(stackList) == 0 {
		return updatedList, nil
	}

	for idx, s := range stackList {
		if s.ID != stack.ID {
			continue
		}

		stackList = append(stackList[:idx], stackList[idx+1:]...)
		break
	}

	if err := storage.Store.Stack.Set(user, stackList); err != nil {
		log.Println(err)
		return updatedList, err
	}

	updatedList, err = json.Marshal(stackList)
	if err != nil {
		return updatedList, err
	}

	return updatedList, nil
}

func getAllStacks(ctx context.Context, user string, includeOrgs bool) []Stack {
	stackList, _ := storage.Store.Stack.Get(user)

	if includeOrgs {
		stackList = addOrgaRepos(ctx, stackList)
	}

	return stackList
}

func addOrgaRepos(ctx context.Context, existingStackList []storage.StackModel) []storage.StackModel {
	client := auth.AuthenticateUser(ctx.Value(auth.ContextKeyToken).(string))
	userName := auth.AuthenticateUser(ctx.Value(auth.ContextKeyUserName).(string))

	orgas, _, err := client.Organizations.ListOrgMemberships(ctx, &github.ListOrgMembershipsOptions{
		ListOptions: github.ListOptions{},
	})
	if err != nil {
		log.Println("could not fetch Orga for user", userName)
		return existingStackList
	}

	for _, org := range orgas {
		log.Println("ORG", *org.Organization.ID)
		orgStacks, _ := storage.Store.Stack.Get(*org.Organization.Login)
		existingStackList = append(existingStackList, orgStacks...)
	}

	return existingStackList
}

func getStackByID(ctx context.Context, user string, stackID string) Stack {
	stack := Stack{}
	stackListArr := getAllStacks(ctx, user, false)

	for _, s := range stackListArr {
		if s.ID != stackID {
			continue
		}

		stack = s
	}

	return stack
}

func generateRandomID() string {
	random := make([]byte, 16)
	rand.Read(random)

	return fmt.Sprintf("%x", random)
}
