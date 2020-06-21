package settings

import (
	"encoding/json"
	"net/http"

	"github.com/SlootSantos/janus-server/pkg/api/auth"
	"github.com/SlootSantos/janus-server/pkg/storage"
)

type Settings struct {
	User  *storage.AllowedUserSettings   `json:"user"`
	Orgas []*storage.AllowedOrgaSettings `json:"orgas"`
}

func HandleHTTP(w http.ResponseWriter, req *http.Request) {
	userName := req.Context().Value(auth.ContextKeyUserName).(string)
	token := req.Context().Value(auth.ContextKeyToken).(string)

	user, _ := storage.Store.User.Get(userName)
	client := auth.AuthenticateUser(token)

	allowedSettings := user.GetAllowedSettings(userName)
	allowedOrgaSettings := user.GetAllowedOrgaSettings(req.Context(), client, userName)

	set := Settings{
		User:  allowedSettings,
		Orgas: allowedOrgaSettings,
	}

	settingJSON, _ := json.Marshal(set)
	w.Write(settingJSON)
}
