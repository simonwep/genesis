package routes

import (
	"github.com/simonwep/genisis/core"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"testing"
)

func TestRegisterAndLogin(t *testing.T) {
	core.DropDatabase()

	tryUnauthorizedPost(UnauthorizedPostConfig{
		Url:  "/register",
		Body: "{\"user\": \"foo\", \"password\": \"test\"}",
		Handler: func(response *httptest.ResponseRecorder) {
			assert.Equal(t, 201, response.Code)
			assert.Equal(t, 0, response.Body.Len())
		},
	})

	tryUnauthorizedPost(UnauthorizedPostConfig{
		Url:  "/login",
		Body: "{\"user\": \"foo\", \"password\": \"test\"}",
		Handler: func(response *httptest.ResponseRecorder) {
			assert.Equal(t, 200, response.Code)
			assert.Equal(t, 0, response.Body.Len())
		},
	})

	tryUnauthorizedPost(UnauthorizedPostConfig{
		Url:  "/login",
		Body: "{\"user\": \"foo\", \"password\": \"test2\"}",
		Handler: func(response *httptest.ResponseRecorder) {
			assert.Equal(t, 401, response.Code)
			assert.Equal(t, 0, response.Body.Len())
		},
	})
}

func TestRegisterIncorrect(t *testing.T) {
	core.DropDatabase()

	tryUnauthorizedPost(UnauthorizedPostConfig{
		Url:  "/register",
		Body: "{\"user\": \"foo2\", \"password\": \"test\"}",
		Handler: func(response *httptest.ResponseRecorder) {
			assert.Equal(t, 401, response.Code)
			assert.Equal(t, 0, response.Body.Len())
		},
	})
}

func TestRegisterDuplicate(t *testing.T) {
	core.DropDatabase()

	tryUnauthorizedPost(UnauthorizedPostConfig{
		Url:  "/register",
		Body: "{\"user\": \"foo\", \"password\": \"test\"}",
		Handler: func(response *httptest.ResponseRecorder) {
			assert.Equal(t, 201, response.Code)
			assert.Equal(t, 0, response.Body.Len())
		},
	})

	tryUnauthorizedPost(UnauthorizedPostConfig{
		Url:  "/register",
		Body: "{\"user\": \"foo\", \"password\": \"test\"}",
		Handler: func(response *httptest.ResponseRecorder) {
			assert.Equal(t, 401, response.Code)
			assert.Equal(t, 0, response.Body.Len())
		},
	})
}
