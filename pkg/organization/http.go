package organization

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/SlootSantos/janus-server/pkg/api/auth"
	"github.com/SlootSantos/janus-server/pkg/stacker"
	"github.com/SlootSantos/janus-server/pkg/storage"
)

type OrganizationCreateConfiguration struct {
	Name          string
	ThirdPartyAWS storage.ThirdPartyAWS
}

const githubRoleAdmin = "admin"

var methodHandlerMap = map[string]http.HandlerFunc{
	http.MethodPost: handlePOST,
	http.MethodGet:  handleGET,
	// http.MethodDelete: handleDELETE,
}

func HandleHTTP(w http.ResponseWriter, req *http.Request) {
	handler, ok := methodHandlerMap[req.Method]
	if !ok {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	handler(w, req)
}

func isOrgAdmin(ctx context.Context, orgName string) (allowed bool, status int) {
	username := ctx.Value(auth.ContextKeyUserName).(string)
	token := ctx.Value(auth.ContextKeyToken).(string)

	client := auth.AuthenticateUser(token)

	orgMember, _, err := client.Organizations.GetOrgMembership(ctx, username, orgName)
	if err != nil {
		return false, http.StatusBadRequest
	}

	if *orgMember.Role != githubRoleAdmin {
		return false, http.StatusUnauthorized
	}

	return true, http.StatusOK
}

func handlePOST(w http.ResponseWriter, req *http.Request) {
	var orgConfig OrganizationCreateConfiguration

	err := json.NewDecoder(req.Body).Decode(&orgConfig)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ok, status := isOrgAdmin(req.Context(), orgConfig.Name)
	if !ok {
		w.WriteHeader(status)
		return
	}

	thirdpartyAccountOut, _ := stacker.SetupThirdPartyAccount(req.Context(), &orgConfig.ThirdPartyAWS)

	orgConfig.ThirdPartyAWS.LambdaARN = thirdpartyAccountOut.LambdaARN
	orgConfig.ThirdPartyAWS.HostedZoneID = thirdpartyAccountOut.HostedZoneID
	orga := &storage.UserModel{
		ThirdPartyAWS: &orgConfig.ThirdPartyAWS,
		User:          orgConfig.Name,
		Type:          storage.TypeOrganization,
	}

	storage.Store.User.Set(orgConfig.Name, orga)

	w.WriteHeader(http.StatusOK)
}

func handleGET(w http.ResponseWriter, req *http.Request) {
	query := req.URL.Query()
	orgName := query.Get("orgaName")

	ok, status := isOrgAdmin(req.Context(), orgName)
	if !ok {
		w.WriteHeader(status)
		return
	}

	res, err := storage.Store.User.Get(orgName)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	log.Printf("%+v", res)
	log.Printf("%+v", res.ThirdPartyAWS)
	w.WriteHeader(http.StatusOK)
}
