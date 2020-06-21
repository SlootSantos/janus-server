package storage

import (
	"context"
	"encoding/json"
	"log"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/go-redis/redis/v7"
	"github.com/google/go-github/github"
)

type UserModel struct {
	User          string `json:"user"`
	Token         string `json:"token"`
	IsPro         bool   `json:"isPro"`
	Type          string `json:"type"`
	Billing       *UserBilling
	ThirdPartyAWS *ThirdPartyAWS
}

type AllowedUserSettings struct {
	IsPro         bool   `json:"isPro"`
	Type          string `json:"type"`
	Name          string `json:"name"`
	ThirdPartyAWS *ThirdPartyAWS
}

type AllowedOrgaSettings struct {
	UserMemberStatus string `json:"userMemberStatus"`
	AllowedUserSettings
}

type ThirdPartyAWS struct {
	AccessKey    string `json:"accessKey"`
	SecretKey    string `json:"secretKey"`
	Domain       string `json:"domain"`
	LambdaARN    string `json:"lambdaARN"`
	HostedZoneID string `json:"hostedZoneId"`
}

type UserBilling struct {
	SubscriptionID string `json:"subscriptionId"`
}

type user struct {
	rC *redis.Client
	dC *dynamoClient
}

const dynamoUserPrimaryKey = "user"
const (
	TypeOrganization = "Organization"
	TypeUser         = "User"
)

func newUserDB(db db, s *session.Session) *user {
	return &user{
		rC: connectRedis(db.redisID),
		dC: connectDynamo(s, db.dynamoTable, dynamoUserPrimaryKey),
	}
}

func (u *user) Get(key string) (*UserModel, error) {
	result, err := u.rC.Get(key).Result()
	if err == redis.Nil {
		log.Println("Initalizing empty redis 'User' DB")
		go u.Set(key, &UserModel{})

		return &UserModel{}, err
	}

	if result == "" {
		return u.getFromDynamo(key)
	}

	model := &UserModel{}
	err = json.Unmarshal([]byte(result), model)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return model, nil
}

func (u *user) getFromDynamo(key string) (*UserModel, error) {
	dbResult := u.dC.GetItem(key)

	model := UserModel{}
	err := dynamodbattribute.UnmarshalMap(dbResult.Item, &model)
	if err != nil {
		log.Println("could not get dynamo", err)
		return &UserModel{}, err
	}

	if model.User == "" {
		return &UserModel{}, nil
	}

	return &model, nil
}

func (u *user) Set(key string, value *UserModel) error {
	jsonValue, err := json.Marshal(value)
	if err != nil {
		log.Println("cannot marshal user", err)
		return err
	}

	err = u.rC.Set(key, jsonValue, 0).Err()
	if err != nil {
		log.Println("failed to store user in redis", err)
		return err
	}

	err = u.dC.SetItem(value)
	if err != nil {
		log.Println("cloud not set to dynamo", err)
	}

	return nil
}

func (u *UserModel) GetAllowedSettings(name string) *AllowedUserSettings {
	var thirdParty *ThirdPartyAWS

	if u.ThirdPartyAWS != nil {
		orgAccess := u.ThirdPartyAWS.AccessKey
		maskedAccess := orgAccess[0:2] + "*******" + orgAccess[len(orgAccess)-3:]

		orgSecret := u.ThirdPartyAWS.SecretKey
		maskedSecret := orgSecret[0:2] + "*******" + orgSecret[len(orgSecret)-3:]
		thirdParty = &ThirdPartyAWS{
			AccessKey:    maskedAccess,
			SecretKey:    maskedSecret,
			Domain:       u.ThirdPartyAWS.Domain,
			LambdaARN:    u.ThirdPartyAWS.LambdaARN,
			HostedZoneID: u.ThirdPartyAWS.HostedZoneID,
		}
	}

	allowedSettings := &AllowedUserSettings{
		Type:          u.Type,
		IsPro:         u.IsPro,
		ThirdPartyAWS: thirdParty,
		Name:          name,
	}

	return allowedSettings
}

func (u *UserModel) GetAllowedOrgaSettings(ctx context.Context, client *github.Client, username string) []*AllowedOrgaSettings {
	orgas, _, err := client.Organizations.ListOrgMemberships(ctx, &github.ListOrgMembershipsOptions{
		ListOptions: github.ListOptions{},
	})
	if err != nil {
		log.Println("could not fetch Orga for user", username)
		return []*AllowedOrgaSettings{}
	}

	orgSettings := []*AllowedOrgaSettings{}

	for _, org := range orgas {
		orgSetting, _ := Store.User.Get(*org.Organization.Login)
		orgMember, _, err := client.Organizations.GetOrgMembership(ctx, username, *org.Organization.Login)
		if err != nil {
			log.Println("uuops", err.Error())
		}
		s := &AllowedOrgaSettings{
			*orgMember.Role,
			*orgSetting.GetAllowedSettings(*org.Organization.Login),
		}
		orgSettings = append(orgSettings, s)
	}

	return orgSettings
}
