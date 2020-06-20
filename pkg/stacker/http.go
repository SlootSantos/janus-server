package stacker

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/SlootSantos/janus-server/pkg/api/auth"
	"github.com/SlootSantos/janus-server/pkg/jam"
	"github.com/SlootSantos/janus-server/pkg/session"
	"github.com/SlootSantos/janus-server/pkg/storage"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/aws/aws-sdk-go/service/route53"
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

	if config.Repository.Name == "" {
		http.Error(w, "Missing Repository param", http.StatusBadRequest)
		return
	}

	ctx := context.WithValue(req.Context(), auth.ContextKeyIsThirdParty, config.IsThirdParty)

	creator, err := s.launchCreator(ctx, &launchParamOptions{
		isThirdParty: config.IsThirdParty,
		ownerType:    config.Repository.Type,
		ownerName:    config.Repository.Owner,
	})
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
	creator, err := s.launchCreator(req.Context(), &launchParamOptions{
		isThirdParty: false,
	})
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

	creator, err := s.launchCreator(ctx, &launchParamOptions{
		isThirdParty: config.IsThirdParty,
		ownerType:    config.Repository.Type,
		ownerName:    config.Repository.Owner,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	newlist, err := creator.Delete(ctx, config)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Write(newlist)
}

const RouteCredentialsPrefix = "/creds"

type ThirdPartyAWSCredentials struct {
	ThirdPartyAWS storage.ThirdPartyAWS
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

	if (creds.ThirdPartyAWS == storage.ThirdPartyAWS{}) || creds.ThirdPartyAWS.AccessKey == "" || creds.ThirdPartyAWS.SecretKey == "" || creds.ThirdPartyAWS.Domain == "" {
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
		AccessKey: creds.ThirdPartyAWS.AccessKey,
		SecretKey: creds.ThirdPartyAWS.SecretKey,
		Domain:    creds.ThirdPartyAWS.Domain,
	}

	thirdpartyAccountOut, _ := SetupThirdPartyAccount(req.Context(), user.ThirdPartyAWS)
	user.ThirdPartyAWS.LambdaARN = thirdpartyAccountOut.LambdaARN
	user.ThirdPartyAWS.HostedZoneID = thirdpartyAccountOut.HostedZoneID

	err = storage.Store.User.Set(req.Context().Value(auth.ContextKeyUserName).(string), user)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

type setupThirdPartyOutput struct {
	HostedZoneID string
	LambdaARN    string
}

func SetupThirdPartyAccount(ctx context.Context, thirdPartyCreds *storage.ThirdPartyAWS) (*setupThirdPartyOutput, error) {
	thirdPartySession, _ := session.AWSSessionThirdParty(thirdPartyCreds.AccessKey, thirdPartyCreds.SecretKey)

	r53 := route53.New(thirdPartySession)
	zones, _ := r53.ListHostedZonesByName(&route53.ListHostedZonesByNameInput{
		DNSName: &thirdPartyCreds.Domain,
	})

	zoneID := *zones.HostedZones[0].Id

	i := iam.New(thirdPartySession)
	roleName := "stackers-cdn-routing-lambda-policy"
	functionName := "stackers-cdn-routing-origin-lambda"
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
				return nil, fmt.Errorf("Can not create role for AWS account: %s", err.Error())
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
		FunctionName: aws.String(functionName),
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
				log.Printf("User: already has function named: %s", functionName)
				existingFunction, _ := lam.GetFunction(&lambda.GetFunctionInput{
					FunctionName: aws.String(functionName),
				})
				versions, _ := lam.ListVersionsByFunction(&lambda.ListVersionsByFunctionInput{
					FunctionName: &functionName,
				})
				latestVersion := strconv.Itoa(len(versions.Versions) - 1)

				functionARN = *existingFunction.Configuration.FunctionArn + ":" + latestVersion
			default:
				log.Println("Could not handle AWS err", err.Error())
				return nil, fmt.Errorf("Can not create role for AWS account: %s", err.Error())
			}
		}

	} else {
		_, err := lam.PublishVersion(&lambda.PublishVersionInput{
			FunctionName: aws.String(*res.FunctionArn),
			Description:  aws.String("Stackers.io Routing Lambda Function"),
		})
		if err != nil {
			log.Println("Could not publish function", err.Error())
			return nil, err
		}

		versions, _ := lam.ListVersionsByFunction(&lambda.ListVersionsByFunctionInput{
			FunctionName: &functionName,
		})
		latestVersion := strconv.Itoa(len(versions.Versions) - 1)

		functionARN = *res.FunctionArn + ":" + latestVersion
	}

	out := &setupThirdPartyOutput{
		LambdaARN:    functionARN,
		HostedZoneID: zoneID,
	}

	return out, nil
}

type launchParamOptions struct {
	isThirdParty bool
	ownerType    string
	ownerName    string
}

func setLaunchParam(ctx context.Context, options *launchParamOptions) (*launchParam, error) {
	creatorParam := launchParam{}
	username := ctx.Value(auth.ContextKeyUserName).(string)

	if options.ownerType == storage.TypeOrganization {
		username = options.ownerName
	}

	if options.isThirdParty {
		user, err := storage.Store.User.Get(username)
		if err != nil {
			return nil, err
		}

		creds := awsCredentials{
			accessKey: user.ThirdPartyAWS.AccessKey,
			secretKey: user.ThirdPartyAWS.SecretKey,
		}

		creatorParam = launchParam{
			hasOwnCreds:  true,
			creds:        creds,
			domain:       user.ThirdPartyAWS.Domain,
			lambdaARN:    user.ThirdPartyAWS.LambdaARN,
			hostedZoneID: user.ThirdPartyAWS.HostedZoneID,
		}

		log.Printf("creatorParam %+v", creatorParam)
	}

	return &creatorParam, nil
}
