package routes

import (
	"github.com/simonwep/genisis/core"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRegisterAndLogin(t *testing.T) {
	core.DropDatabase()

	tryUnauthorizedPost("/register", UnauthorizedBodyConfig{
		Body: "{\"user\": \"foo\", \"password\": \"test\"}",
		Handler: func(response *httptest.ResponseRecorder) {
			assert.Equal(t, http.StatusCreated, response.Code)
			assert.Equal(t, 0, response.Body.Len())
		},
	})

	tryUnauthorizedPost("/login", UnauthorizedBodyConfig{
		Body: "{\"user\": \"foo\", \"password\": \"test\"}",
		Handler: func(response *httptest.ResponseRecorder) {
			assert.Equal(t, http.StatusOK, response.Code)
			assert.Equal(t, 0, response.Body.Len())
		},
	})

	tryUnauthorizedPost("/login", UnauthorizedBodyConfig{
		Body: "{\"user\": \"foo\", \"password\": \"test2\"}",
		Handler: func(response *httptest.ResponseRecorder) {
			assert.Equal(t, http.StatusUnauthorized, response.Code)
			assert.Equal(t, 0, response.Body.Len())
		},
	})
}

func TestRegisterIncorrect(t *testing.T) {
	core.DropDatabase()

	tryUnauthorizedPost("/register", UnauthorizedBodyConfig{
		Body: "{\"user\": \"foo2\", \"password\": \"test\"}",
		Handler: func(response *httptest.ResponseRecorder) {
			assert.Equal(t, http.StatusUnauthorized, response.Code)
			assert.Equal(t, 0, response.Body.Len())
		},
	})
}

func TestRegisterDuplicate(t *testing.T) {
	core.DropDatabase()

	tryUnauthorizedPost("/register", UnauthorizedBodyConfig{
		Body: "{\"user\": \"foo\", \"password\": \"test\"}",
		Handler: func(response *httptest.ResponseRecorder) {
			assert.Equal(t, http.StatusCreated, response.Code)
			assert.Equal(t, 0, response.Body.Len())
		},
	})

	tryUnauthorizedPost("/register", UnauthorizedBodyConfig{
		Body: "{\"user\": \"foo\", \"password\": \"test\"}",
		Handler: func(response *httptest.ResponseRecorder) {
			assert.Equal(t, http.StatusUnauthorized, response.Code)
			assert.Equal(t, 0, response.Body.Len())
		},
	})
}
