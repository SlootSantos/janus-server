package auth

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/SlootSantos/janus-server/pkg/storage"
	"github.com/golang/mock/gomock"
)

func TestHandleLogin(t *testing.T) {
	t.Run("should redirect for login", func(t *testing.T) {
		// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
		// pass 'nil' as the third parameter.
		req, err := http.NewRequest("GET", Login, nil)
		if err != nil {
			t.Fatal(err)
		}
		// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(HandleLogin)
		// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
		// directly and pass in our Request and ResponseRecorder.
		handler.ServeHTTP(rr, req)
		// Check the status code is what we expect.
		if status := rr.Code; status != http.StatusTemporaryRedirect {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusTemporaryRedirect)
		}
	})
}

func TestHandleLoginCheck(t *testing.T) {
	t.Run("should set correct headers", func(t *testing.T) {
		req, err := http.NewRequest("GET", LoginCheck, nil)
		if err != nil {
			t.Fatal(err)
		}

		req.Header.Set("Origin", "http://localhost:3000")

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(HandleLoginCheck)

		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusOK)
		}

		expectedAccessControlOriginHeader := "http://localhost:3000"
		if rr.Header().Get("Access-Control-Allow-Origin") != expectedAccessControlOriginHeader {
			t.Errorf("handler returned unexpected body: got %v want %v",
				rr.Header().Get("Access-Control-Allow-Origin"), expectedAccessControlOriginHeader)
		}

		expectedAccessControlCredHeader := "true"
		if rr.Header().Get("Access-Control-Allow-Credentials") != expectedAccessControlCredHeader {
			t.Errorf("handler returned unexpected body: got %v want %v",
				rr.Header().Get("Access-Control-Allow-Credentials"), expectedAccessControlCredHeader)
		}

		expectedContentTypeHeader := "application/json"
		if rr.Header().Get("Content-Type") != expectedContentTypeHeader {
			t.Errorf("handler returned unexpected body: got %v want %v",
				rr.Header().Get("Content-Type"), expectedContentTypeHeader)
		}
	})

	t.Run("should return false if no cookie", func(t *testing.T) {
		req, err := http.NewRequest("GET", LoginCheck, nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(HandleLoginCheck)

		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusOK)
		}

		expectedBody := &loginCheckResponse{}
		json.Unmarshal(rr.Body.Bytes(), expectedBody)
		if expectedBody.LoggedIn {
			t.Errorf("handler returned unexpected body: got %v want %v",
				expectedBody.LoggedIn, false)
		}

		if expectedBody.User != nil {
			t.Errorf("handler returned unexpected body: got %v want %v",
				expectedBody.User, nil)
		}
	})

	t.Run("should return true if cookie", func(t *testing.T) {
		req, err := http.NewRequest("GET", LoginCheck, nil)
		if err != nil {
			t.Fatal(err)
		}

		cookieValue, _ := CreateJWT(&authUser{
			Name: "SlootSantos",
		})

		req.AddCookie(&http.Cookie{Name: OAuthCookieName, Value: cookieValue})

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(HandleLoginCheck)

		ctrl := gomock.NewController(t)
		userMock, _, _ := storage.MockInit(ctrl)

		userGetReturn := &storage.UserModel{}
		userMock.EXPECT().Get("SlootSantos").Times(1).Return(userGetReturn, nil)

		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusOK)
		}

		expectedBody := &loginCheckResponse{}
		json.Unmarshal(rr.Body.Bytes(), expectedBody)
		if !expectedBody.LoggedIn {
			t.Errorf("handler returned unexpected body: got %v want %v",
				expectedBody.LoggedIn, true)
		}

		if expectedBody.User == nil {
			t.Errorf("handler returned unexpected body: got %v want %v",
				expectedBody.User, "something but nil")
		}
	})
}
