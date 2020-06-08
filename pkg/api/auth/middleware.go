package auth

import (
	"context"
	"log"
	"net/http"

	"github.com/SlootSantos/janus-server/pkg/storage"
)

type key string

// ContextKeyToken constant for context
const ContextKeyToken key = "token"

// ContextKeyUserName constant for context
const ContextKeyUserName key = "userName"

// WithCredentials is a middleware function
func WithCredentials(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		origin := req.Header.Get("Origin") // TODO! Whitelist

		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Allow-Methods", "GET,POST,OPTIONS,DELETE")

		if req.Method == http.MethodOptions {
			w.WriteHeader(200)
			return
		}

		cookie, _ := req.Cookie(OAuthCookieName)

		if cookie == nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Missing Authentication Cookie"))
			return
		}

		newReq := createCtxWithToken(cookie, req)
		next(w, newReq)
	}
}

func createCtxWithToken(cookie *http.Cookie, req *http.Request) *http.Request {
	authUser, _ := GetUserFromCookie(cookie)
	userModel, _ := storage.Store.User.Get(authUser.Name)

	log.Printf("%+v", userModel.Billing)

	ctx := context.WithValue(req.Context(), ContextKeyToken, userModel.Token)
	ctx = context.WithValue(ctx, ContextKeyUserName, authUser.Name)

	return req.WithContext(ctx)
}
