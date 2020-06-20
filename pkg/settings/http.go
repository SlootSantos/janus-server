package settings

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/SlootSantos/janus-server/pkg/api/auth"
	"github.com/SlootSantos/janus-server/pkg/storage"
	"github.com/google/go-github/github"
)

type Settings struct {
	User  *storage.AllowedUserSettings   `json:"user"`
	Orgas []*storage.AllowedOrgaSettings `json:"orgas"`
}

func HandleHTTP(w http.ResponseWriter, req *http.Request) {
	userName := req.Context().Value(auth.ContextKeyUserName).(string)
	user, _ := storage.Store.User.Get(userName)

	allowedSettings := sanitizeSettingsOutput(user, userName)
	allSettings := addOrgaRepos(req.Context())
	set := Settings{
		User:  allowedSettings,
		Orgas: allSettings,
	}

	settingJSON, _ := json.Marshal(set)

	w.Write(settingJSON)
}

func sanitizeSettingsOutput(user *storage.UserModel, name string) *storage.AllowedUserSettings {
	var thirdParty *storage.ThirdPartyAWS

	if user.ThirdPartyAWS != nil {
		orgAccess := user.ThirdPartyAWS.AccessKey
		maskedAccess := orgAccess[0:2] + "*******" + orgAccess[len(orgAccess)-3:]

		orgSecret := user.ThirdPartyAWS.SecretKey
		maskedSecret := orgSecret[0:2] + "*******" + orgSecret[len(orgSecret)-3:]
		thirdParty = &storage.ThirdPartyAWS{
			AccessKey:    maskedAccess,
			SecretKey:    maskedSecret,
			Domain:       user.ThirdPartyAWS.Domain,
			LambdaARN:    user.ThirdPartyAWS.LambdaARN,
			HostedZoneID: user.ThirdPartyAWS.HostedZoneID,
		}
	}

	allowedSettings := &storage.AllowedUserSettings{
		Type:          user.Type,
		IsPro:         user.IsPro,
		ThirdPartyAWS: thirdParty,
		Name:          name,
	}

	return allowedSettings
}

func addOrgaRepos(ctx context.Context) []*storage.AllowedOrgaSettings {
	client := auth.AuthenticateUser(ctx.Value(auth.ContextKeyToken).(string))
	username := ctx.Value(auth.ContextKeyUserName).(string)

	orgas, _, err := client.Organizations.ListOrgMemberships(ctx, &github.ListOrgMembershipsOptions{
		ListOptions: github.ListOptions{},
	})
	if err != nil {
		// log.Println("could not fetch Orga for user", userName)
		return []*storage.AllowedOrgaSettings{}
	}

	orgSettings := []*storage.AllowedOrgaSettings{}

	for _, org := range orgas {
		orgSetting, _ := storage.Store.User.Get(*org.Organization.Login)
		orgMember, _, err := client.Organizations.GetOrgMembership(ctx, username, *org.Organization.Login)
		if err != nil {
			log.Println("uuops", err.Error())
		}
		s := &storage.AllowedOrgaSettings{
			*orgMember.Role,
			*sanitizeSettingsOutput(orgSetting, *org.Organization.Login),
		}
		orgSettings = append(orgSettings, s)
	}

	return orgSettings
}
