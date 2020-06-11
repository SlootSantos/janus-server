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
		Repo: struct{ Name string }{
			Name: config.Repository,
		},
	}

	creationResult, err := j.buildStack(ctx, creationParam)
	creationStack := Stack(*creationResult)
	if err != nil {
		return nil, err
	}

	userName := ctx.Value(auth.ContextKeyUserName).(string)
	creationStack.IsThirdParty = config.IsThirdParty

	updatedList, err := storeStack(&creationStack, userName)
	if err != nil {
		return []byte{}, err
	}

	return updatedList, nil
}

func (j *Creator) Delete(ctx context.Context, stackID string) ([]byte, error) {
	user := ctx.Value(auth.ContextKeyUserName).(string)

	stack := getStackByID(stackID, user)
	if stack == (Stack{}) {
		err := errors.New("No Stack ID: " + stackID + " for User: " + user)
		log.Println(err.Error())

		return []byte{}, err
	}

	deletionParam := DeletionParam(stack)
	j.destroyStack(ctx, &deletionParam)

	updatedList, err := removeStack(&stack, user)
	if err != nil {
		log.Println(err)
		return []byte{}, err
	}

	return updatedList, nil
}

func (j *Creator) List(ctx context.Context) ([]Stack, error) {
	stacks := getAllStacks(ctx.Value(auth.ContextKeyUserName).(string))

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

func storeStack(stack *Stack, user string) (updatedList []byte, err error) {
	stackList := getAllStacks(user)
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

func removeStack(stack *Stack, user string) (updatedList []byte, err error) {
	stackList := getAllStacks(user)
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

func getAllStacks(user string) []Stack {
	stackList, _ := storage.Store.Stack.Get(user)

	return stackList
}

func getStackByID(stackID string, user string) Stack {
	stack := Stack{}
	stackListArr := getAllStacks(user)

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
