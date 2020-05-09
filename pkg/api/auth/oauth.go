package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	githuboauth "golang.org/x/oauth2/github"
)

type authUser struct {
	Name     string
	CLIToken string
}

var (
	oauthConf        *oauth2.Config
	oauthStateString string
	oauthSigningKey  []byte
	OAuthCookieName  = "janus_cookie"
)

func OauthConf() *oauth2.Config {
	if oauthConf == nil {
		oauthConf = &oauth2.Config{
			ClientID:     os.Getenv("OAUTH_CLIENT_ID"),
			ClientSecret: os.Getenv("OAUTH_CLIENT_SECRET"),
			Scopes:       []string{"user:email", "repo"},
			Endpoint:     githuboauth.Endpoint,
		}
	}

	return oauthConf
}

func OauthStateString() string {
	return os.Getenv("OAUTH_CLIENT_STATE")
}

func OauthSigningKey() []byte {
	return []byte(os.Getenv("OAUTH_CLIENT_SIGNING_KEY"))
}

// AuthenticateUser authes the user & gets a new github client
func AuthenticateUser(tokenStr string) *github.Client {
	token, _ := tokenFromJSON(tokenStr)
	oauthClient := OauthConf().Client(oauth2.NoContext, token)
	client := github.NewClient(oauthClient)

	return client
}

func GetUserFromCookie(cookie *http.Cookie) (*authUser, error) {
	return decodeJWT(cookie.Value)
}

func decodeJWT(tokenStr string) (*authUser, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return OauthSigningKey(), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		user := &authUser{
			Name: fmt.Sprintf("%v", claims["user_name"]),
		}

		return user, nil
	}

	return nil, err
}

func CreateJWT(user *authUser) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_name": user.Name,
		"nbf":       time.Date(2015, 10, 10, 12, 0, 0, 0, time.UTC).Unix(),
	})

	return token.SignedString(OauthSigningKey())
}

func TokenToJSON(token *oauth2.Token) (string, error) {
	d, err := json.Marshal(token)

	if err != nil {
		return "", err
	}

	fmt.Println(d)
	return string(d), nil
}

func tokenFromJSON(jsonStr string) (*oauth2.Token, error) {
	var token oauth2.Token
	if err := json.Unmarshal([]byte(jsonStr), &token); err != nil {
		return nil, err
	}

	return &token, nil
}
