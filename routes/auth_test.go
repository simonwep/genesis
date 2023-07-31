package routes

import (
	"github.com/simonwep/genesis/core"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestInvalidLogin(t *testing.T) {
	core.ResetDatabase()

	tryUnauthorizedPost("/login", UnauthorizedBodyConfig{
		Body: "{\"user\": \"foo123\", \"password\": \"hgEiPCZP\"}",
		Handler: func(response *httptest.ResponseRecorder) {
			assert.Equal(t, http.StatusUnauthorized, response.Code)
		},
	})
}

func TestLogin(t *testing.T) {
	core.ResetDatabase()

	tryUnauthorizedPost("/login", UnauthorizedBodyConfig{
		Body: "{\"user\": \"foo\", \"password\": \"hgEiPCZP\"}",
		Handler: func(response *httptest.ResponseRecorder) {
			assert.Equal(t, http.StatusOK, response.Code)
		},
	})
}

func TestLogout(t *testing.T) {
	core.ResetDatabase()

	var cookie1 string
	var cookie2 string

	tryUnauthorizedPost("/login", UnauthorizedBodyConfig{
		Body: "{\"user\": \"foo\", \"password\": \"hgEiPCZP\"}",
		Handler: func(response *httptest.ResponseRecorder) {
			assert.Equal(t, http.StatusOK, response.Code)
			cookie1 = response.Header().Get("Set-Cookie")
		},
	})

	// invalidate first cookie by rotating the key
	tryAuthorizedPost("/login/refresh", AuthorizedBodyConfig{
		Cookie: cookie1,
		Handler: func(response *httptest.ResponseRecorder) {
			assert.Equal(t, http.StatusOK, response.Code)
			cookie2 = response.Header().Get("Set-Cookie")
		},
	})

	// invalidate the second one directly
	tryAuthorizedPost("/logout", AuthorizedBodyConfig{
		Cookie: cookie2,
		Handler: func(response *httptest.ResponseRecorder) {
			assert.Equal(t, http.StatusOK, response.Code)
		},
	})

	// both tokens should be blacklisted
	tryAuthorizedPost("/login/refresh", AuthorizedBodyConfig{
		Cookie: cookie1,
		Handler: func(response *httptest.ResponseRecorder) {
			assert.Equal(t, http.StatusUnauthorized, response.Code)
		},
	})

	tryAuthorizedPost("/login/refresh", AuthorizedBodyConfig{
		Cookie: cookie2,
		Handler: func(response *httptest.ResponseRecorder) {
			assert.Equal(t, http.StatusUnauthorized, response.Code)
		},
	})
}
