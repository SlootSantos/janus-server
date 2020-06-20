package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/SlootSantos/janus-server/pkg/storage"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

// Login is the login route identifier
const Login = "/login"

// LoginCheck is the identifier for the route to check if user is logged-in
const LoginCheck = "/login/check"

// Callback is the callback route identifier
const Callback = "/callback"

func HandleLoginCheck(w http.ResponseWriter, req *http.Request) {
	origin := req.Header.Get("Origin")
	cookie, _ := req.Cookie(OAuthCookieName)
	cookieSet := true
	w.Header().Set("Access-Control-Allow-Origin", origin)
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Content-Type", "application/json")

	if cookie == nil {
		cookieSet = false
	}

	data := struct{ LoggedIn bool }{LoggedIn: cookieSet}
	json.NewEncoder(w).Encode(data)
}

// HandleLogin handles all of the Login logic
func HandleLogin(w http.ResponseWriter, req *http.Request) {
	url := OauthConf().AuthCodeURL(OauthStateString(), oauth2.AccessTypeOnline)

	http.Redirect(w, req, url, http.StatusTemporaryRedirect)
}

// HandleCallback handles teh github authentication callback
func HandleCallback(w http.ResponseWriter, req *http.Request) {
	state := req.FormValue("state")

	if state != OauthStateString() {
		fmt.Printf("invalid oauth state, expected '%s', got '%s'\n", OauthStateString(), state)
		http.Redirect(w, req, "/", http.StatusTemporaryRedirect)
		return
	}

	user, tokenStr := getUser(req)

	storeUser(user, tokenStr)
	setCookie(w, user)

	http.Redirect(w, req, os.Getenv("CLIENT_URL")+"/admin/dashboard", http.StatusTemporaryRedirect)
}

func getToken(code string) (string, error) {
	token, err := OauthConf().Exchange(oauth2.NoContext, code)
	if err != nil {
		log.Println(err)
	}

	return TokenToJSON(token)

}

func getUser(req *http.Request) (*github.User, string) {
	tokenStr, _ := getToken(req.FormValue("code"))
	client := AuthenticateUser(tokenStr)

	user, _, err := client.Users.Get(context.Background(), "")

	if err != nil {
		log.Printf("client.Users.Get() faled with '%s'\n", err)
	}

	return user, tokenStr
}

func storeUser(gUser *github.User, token string) error {
	user := &storage.UserModel{User: *gUser.Login, Token: token, Type: storage.TypeUser}

	err := storage.Store.User.Set(*gUser.Login, user)
	return err
}

func setCookie(w http.ResponseWriter, user *github.User) {
	cookieValue, _ := CreateJWT(&authUser{
		Name: *user.Login,
	})

	sessionCookie := &http.Cookie{
		SameSite: http.SameSiteLaxMode,
		Secure:   os.Getenv("ENV") != "local",
		Name:     OAuthCookieName,
		Value:    cookieValue,
		Path:     "/",
		Domain:   os.Getenv("SERVER_URL"),
		HttpOnly: true,
	}
	http.SetCookie(w, sessionCookie)
}
