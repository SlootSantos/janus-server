package stacker

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/SlootSantos/janus-server/pkg/api/auth"
	"github.com/SlootSantos/janus-server/pkg/jam"
	"github.com/SlootSantos/janus-server/pkg/session"
	"github.com/SlootSantos/janus-server/pkg/storage"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/lambda"
)

func (s *Stacker) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	methodHandlerMap := map[string]http.HandlerFunc{
		http.MethodGet:    s.handleGET,
		http.MethodPost:   s.handlePOST,
		http.MethodDelete: s.handleDELETE,
	}

	if handler, ok := methodHandlerMap[req.Method]; ok {
		handler(w, req)
		return
	}

	w.WriteHeader(http.StatusMethodNotAllowed)
}

func (s *Stacker) handlePOST(w http.ResponseWriter, req *http.Request) {
	var config jam.StackCreateConfig

	err := json.NewDecoder(req.Body).Decode(&config)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if config.Repository == "" {
		http.Error(w, "Missing Repository param", http.StatusBadRequest)
		return
	}

	ctx := context.WithValue(req.Context(), auth.ContextKeyIsThirdParty, config.IsThirdParty)

	creator, err := s.launchCreator(ctx, config.IsThirdParty)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	updatedList, err := creator.Build(ctx, config)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(updatedList)
}

func (s *Stacker) handleGET(w http.ResponseWriter, req *http.Request) {
	creator, err := s.launchCreator(req.Context(), false)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	stacks, _ := creator.List(req.Context())
	stackJSON, err := json.Marshal(stacks)
	if err != nil {
		log.Println(err)
	}

	w.Write(stackJSON)
}

func (s *Stacker) handleDELETE(w http.ResponseWriter, req *http.Request) {
	var config jam.StackDestroyConfig

	err := json.NewDecoder(req.Body).Decode(&config)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if config.ID == "" {
		http.Error(w, "Missing Stack ID", http.StatusBadRequest)
		return
	}

	ctx := context.WithValue(req.Context(), auth.ContextKeyIsThirdParty, config.IsThirdParty)

	creator, err := s.launchCreator(ctx, config.IsThirdParty)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	newlist, err := creator.Delete(ctx, config.ID)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Write(newlist)
}

const RouteCredentialsPrefix = "/creds"

type ThirdPartyAWSCredentials struct {
	SecretKey string
	AccessKey string
	Domain    string
}

func SetThirdPartyAWSCredentials(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var creds ThirdPartyAWSCredentials

	err := json.NewDecoder(req.Body).Decode(&creds)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if creds.AccessKey == "" || creds.SecretKey == "" || creds.Domain == "" {
		http.Error(w, "AccessKey or SecretKey or Domain missing in request", http.StatusBadRequest)
		return
	}

	user, err := storage.Store.User.Get(req.Context().Value(auth.ContextKeyUserName).(string))
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "User not existant", http.StatusBadRequest)
		return
	}

	user.ThirdPartyAWS = &storage.ThirdPartyAWS{
		AccessKey: creds.AccessKey,
		SecretKey: creds.SecretKey,
		Domain:    creds.Domain,
	}

	lambdaARN, _ := setupThirdPartyAccount(req.Context(), user.ThirdPartyAWS)
	user.ThirdPartyAWS.LambdaARN = lambdaARN

	err = storage.Store.User.Set(req.Context().Value(auth.ContextKeyUserName).(string), user)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func setupThirdPartyAccount(ctx context.Context, thirdPartyCreds *storage.ThirdPartyAWS) (string, error) {
	thirdPartySession, _ := session.AWSSessionThirdParty(thirdPartyCreds.AccessKey, thirdPartyCreds.SecretKey)
	i := iam.New(thirdPartySession)

	roleName := "role-with-policy-11"
	var roleARN string
	role, err := i.CreateRole(&iam.CreateRoleInput{
		Path:                     aws.String("/service-role/"),
		AssumeRolePolicyDocument: aws.String("{\"Version\": \"2012-10-17\",\"Statement\": [{\"Effect\": \"Allow\",\"Principal\": {\"Service\": [\"lambda.amazonaws.com\",\"edgelambda.amazonaws.com\"]},\"Action\": \"sts:AssumeRole\"}]}"),
		RoleName:                 aws.String(roleName),
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case iam.ErrCodeEntityAlreadyExistsException:
				log.Printf("User: already has role named: %s", roleName)
				r, _ := i.GetRole(&iam.GetRoleInput{
					RoleName: aws.String(roleName),
				})
				roleARN = *r.Role.Arn
			default:
				log.Println("Could not handle AWS err", err.Error())
				return "", fmt.Errorf("Can not create role for AWS account: %s", err.Error())
			}
		}
	} else {
		roleARN = *role.Role.Arn
		i.PutRolePolicy(&iam.PutRolePolicyInput{
			PolicyName:     aws.String("stackers-lambda-exec-policy"),
			PolicyDocument: aws.String("{\"Version\": \"2012-10-17\", \"Statement\": [ { \"Effect\": \"Allow\", \"Action\": [ \"logs:CreateLogGroup\", \"logs:CreateLogStream\", \"logs:PutLogEvents\" ], \"Resource\": [ \"arn:aws:logs:*:*:*\" ] } ] }"),
			RoleName:       aws.String(roleName),
		})

		time.Sleep(time.Second * 60)
	}

	var functionARN string
	lam := lambda.New(thirdPartySession)
	res, err := lam.CreateFunction(&lambda.CreateFunctionInput{
		FunctionName: aws.String("stackers-handler-routing-2"),
		Handler:      aws.String("index.handler"),
		Runtime:      aws.String("nodejs12.x"),
		Role:         aws.String(roleARN),
		Code: &lambda.FunctionCode{
			S3Bucket: aws.String("stackers.io-lambda-public-functions"),
			S3Key:    aws.String("exx.zip"),
		},
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case lambda.ErrCodeResourceConflictException:
				log.Printf("User: already has function named: %s", "stackers-handler-routing-2")
				lam.GetFunction(&lambda.GetFunctionInput{
					FunctionName: aws.String("stackers-handler-routing-2"),
				})
			default:
				log.Println("Could not handle AWS err", err.Error())
				return "", fmt.Errorf("Can not create role for AWS account: %s", err.Error())
			}
		}
	} else {
		functionARN = *res.FunctionArn
		lam.PublishVersion(&lambda.PublishVersionInput{
			FunctionName: aws.String(functionARN),
			Description:  aws.String("Stackers.io Routing Lambda Function"),
		})
	}

	return functionARN, nil
}

func setLaunchParam(ctx context.Context, isThirdParty bool) (*launchParam, error) {
	creatorParam := launchParam{}
	username := ctx.Value(auth.ContextKeyUserName).(string)

	if isThirdParty {
		user, err := storage.Store.User.Get(username)
		if err != nil {
			return nil, err
		}

		creatorParam = launchParam{
			hasOwnCreds: true,
			creds: awsCredentials{
				accessKey: user.ThirdPartyAWS.AccessKey,
				secretKey: user.ThirdPartyAWS.SecretKey,
			},
			domain: user.ThirdPartyAWS.Domain,
		}
	}

	return &creatorParam, nil
}
