package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"
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

		expected := "{\"LoggedIn\":false}\n"
		if rr.Body.String() != expected {
			t.Errorf("handler returned unexpected body: got %v want %v",
				rr.Body.String(), expected)
		}
	})

	t.Run("should return true if cookie", func(t *testing.T) {
		req, err := http.NewRequest("GET", LoginCheck, nil)
		if err != nil {
			t.Fatal(err)
		}

		req.AddCookie(&http.Cookie{Name: OAuthCookieName})

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(HandleLoginCheck)

		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusOK)
		}

		expected := "{\"LoggedIn\":true}\n"
		if rr.Body.String() != expected {
			t.Errorf("handler returned unexpected body: got %v want %v",
				rr.Body.String(), expected)
		}
	})
}
