package routes

import (
	"github.com/simonwep/genisis/core"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestInvalidLogin(t *testing.T) {
	core.DropDatabase()

	tryUnauthorizedPost("/login", UnauthorizedBodyConfig{
		Body: "{\"user\": \"foo\", \"password\": \"test\"}",
		Handler: func(response *httptest.ResponseRecorder) {
			assert.Equal(t, http.StatusUnauthorized, response.Code)
		},
	})
}

func TestSuccessfulLogin(t *testing.T) {
	core.DropDatabase()
	core.CreateUser("foo", "test")

	tryUnauthorizedPost("/login", UnauthorizedBodyConfig{
		Body: "{\"user\": \"foo\", \"password\": \"test\"}",
		Handler: func(response *httptest.ResponseRecorder) {
			assert.Equal(t, http.StatusOK, response.Code)
			assert.IsType(t, "", response.Header().Get("Authorization"))
		},
	})
}
