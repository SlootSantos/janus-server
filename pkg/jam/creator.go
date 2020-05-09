package jam

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/SlootSantos/janus-server/pkg/api/auth"
	"github.com/SlootSantos/janus-server/pkg/storage"
)

const bucketPrefix = "janus-bucket-"

// builds entire JAM-Stack
// including a Bucket, a CDN and a Githook
func (j *Creator) build(ctx context.Context, repository string) ([]byte, error) {
	jamID := strconv.FormatInt(time.Now().UnixNano(), 10)
	creationParam := &CreationParam{
		ID: jamID,
		Bucket: struct{ ID string }{
			ID: bucketPrefix + jamID,
		},
		Repo: struct{ Name string }{
			Name: repository,
		},
	}

	creationResult, err := j.buildStack(ctx, creationParam)
	creationStack := Stack(*creationResult)
	if err != nil {
		return nil, err
	}

	userName := ctx.Value(auth.ContextKeyUserName).(string)

	updatedList, err := storeJAM(&creationStack, userName)
	if err != nil {
		return []byte{}, err
	}

	return updatedList, nil
}

func (j *Creator) delete(ctx context.Context, stackID string) ([]byte, error) {
	user := ctx.Value(auth.ContextKeyUserName).(string)

	stack := getStackByID(stackID, user)
	if stack == (Stack{}) {
		err := errors.New("No Stack ID: " + stackID + " for User: " + user)
		log.Println(err.Error())

		return []byte{}, err
	}

	deletionParam := DeletionParam(stack)
	j.destroyStack(ctx, &deletionParam)

	updatedList, err := removeJAM(&stack, user)
	if err != nil {
		log.Println(err)
		return []byte{}, err
	}

	return updatedList, nil
}

func (j *Creator) list(ctx context.Context) ([]Stack, error) {
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

func storeJAM(stack *Stack, user string) (updatedList []byte, err error) {
	jamList := getAllStacks(user)
	updatedStackList := append(jamList, *stack)

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

func removeJAM(stack *Stack, user string) (updatedList []byte, err error) {
	jamList := getAllStacks(user)
	if len(jamList) == 0 {
		return updatedList, nil
	}

	for idx, jam := range jamList {
		if jam.ID != stack.ID {
			continue
		}

		jamList = append(jamList[:idx], jamList[idx+1:]...)
		break
	}

	if err := storage.Store.Stack.Set(user, jamList); err != nil {
		log.Println(err)
		return updatedList, err
	}

	updatedList, err = json.Marshal(jamList)
	if err != nil {
		return updatedList, err
	}

	return updatedList, nil
}

func getAllStacks(user string) []Stack {
	jamList, _ := storage.Store.Stack.Get(user)

	return jamList
}

func getStackByID(jamID string, user string) Stack {
	stack := Stack{}
	jamListArr := getAllStacks(user)

	for _, jam := range jamListArr {
		if jam.ID != jamID {
			continue
		}

		stack = jam
	}

	return stack
}
