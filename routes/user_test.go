package routes

import (
	"encoding/json"
	"github.com/simonwep/genesis/core"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func loginAdmin(t *testing.T) string {
	core.ResetDatabase()
	var token string

	tryUnauthorizedPost("/login", UnauthorizedBodyConfig{
		Body: "{\"user\": \"bar\", \"password\": \"EczUR8dn\"}",
		Handler: func(response *httptest.ResponseRecorder) {
			assert.Equal(t, http.StatusOK, response.Code)
			data := LoginResponse{}
			json.Unmarshal(response.Body.Bytes(), &data)
			token = data.Token
		},
	})

	return token
}

func TestUnauthorizedAccess(t *testing.T) {
	token := loginUser(t)

	tryAuthorizedPost("/user", AuthorizedBodyConfig{
		Body:  "{\"user\": \"test\", \"password\": \"foobar1234\", \"admin\": false}",
		Token: token,
		Handler: func(response *httptest.ResponseRecorder) {
			assert.Equal(t, http.StatusUnauthorized, response.Code)
		},
	})

	tryAuthorizedPost("/user/foo", AuthorizedBodyConfig{
		Body:  "{\"user\": \"test\", \"password\": \"foobar1234\", \"admin\": false}",
		Token: token,
		Handler: func(response *httptest.ResponseRecorder) {
			assert.Equal(t, http.StatusUnauthorized, response.Code)
		},
	})

	tryAuthorizedGet("/user", AuthorizedConfig{
		Token: token,
		Handler: func(response *httptest.ResponseRecorder) {
			assert.Equal(t, http.StatusUnauthorized, response.Code)
		},
	})

	tryAuthorizedDelete("/user/foo", AuthorizedConfig{
		Token: token,
		Handler: func(response *httptest.ResponseRecorder) {
			assert.Equal(t, http.StatusUnauthorized, response.Code)
		},
	})
}

func TestRetrieveUsers(t *testing.T) {
	token := loginAdmin(t)

	tryAuthorizedGet("/user", AuthorizedConfig{
		Token: token,
		Handler: func(response *httptest.ResponseRecorder) {
			bar := "{\"user\":\"bar\",\"admin\":true}"
			foo := "{\"user\":\"foo\",\"admin\":false}"
			assert.Equal(t, "["+bar+","+foo+"]", response.Body.String())
		},
	})
}

func TestDeleteUser(t *testing.T) {
	token := loginAdmin(t)

	tryUnauthorizedPost("/login", UnauthorizedBodyConfig{
		Body: "{\"user\": \"foo\", \"password\": \"hgEiPCZP\"}",
		Handler: func(response *httptest.ResponseRecorder) {
			assert.Equal(t, http.StatusOK, response.Code)
		},
	})

	tryAuthorizedDelete("/user/foo", AuthorizedConfig{
		Token: token,
		Handler: func(response *httptest.ResponseRecorder) {
			assert.Equal(t, http.StatusOK, response.Code)
		},
	})

	tryUnauthorizedPost("/login", UnauthorizedBodyConfig{
		Body: "{\"user\": \"foo\", \"password\": \"hgEiPCZP\"}",
		Handler: func(response *httptest.ResponseRecorder) {
			assert.Equal(t, http.StatusUnauthorized, response.Code)
		},
	})
}

func TestCreateUser(t *testing.T) {
	token := loginAdmin(t)

	tryAuthorizedPost("/user", AuthorizedBodyConfig{
		Token: token,
		Body:  "{\"user\":\"test2\",\"password\":\"foobar1235\",\"admin\":true}",
		Handler: func(response *httptest.ResponseRecorder) {
			assert.Equal(t, http.StatusCreated, response.Code)
		},
	})

	tryAuthorizedPost("/user", AuthorizedBodyConfig{
		Token: token,
		Body:  "{\"user\":\"test//2\",\"password\":\"foobar1235\",\"admin\":true}",
		Handler: func(response *httptest.ResponseRecorder) {
			assert.Equal(t, http.StatusBadRequest, response.Code)
		},
	})

	tryUnauthorizedPost("/login", UnauthorizedBodyConfig{
		Body: "{\"user\": \"test2\", \"password\": \"foobar1235\"}",
		Handler: func(response *httptest.ResponseRecorder) {
			assert.Equal(t, http.StatusOK, response.Code)
		},
	})
}

func TestUpdateUser(t *testing.T) {
	token := loginAdmin(t)

	tryUnauthorizedPost("/login", UnauthorizedBodyConfig{
		Body: "{\"user\":\"foo\", \"password\":\"hgEiPCZP\"}",
		Handler: func(response *httptest.ResponseRecorder) {
			assert.Equal(t, http.StatusOK, response.Code)
		},
	})

	tryAuthorizedPost("/user/foo", AuthorizedBodyConfig{
		Token: token,
		Body:  "{\"user\":\"foo\",\"password\":\"wK8iVkRO\",\"admin\":true}",
		Handler: func(response *httptest.ResponseRecorder) {
			assert.Equal(t, http.StatusOK, response.Code)
		},
	})

	tryUnauthorizedPost("/login", UnauthorizedBodyConfig{
		Body: "{\"user\":\"foo\", \"password\":\"hgEiPCZP\"}",
		Handler: func(response *httptest.ResponseRecorder) {
			assert.Equal(t, http.StatusUnauthorized, response.Code)
		},
	})

	tryUnauthorizedPost("/login", UnauthorizedBodyConfig{
		Body: "{\"user\":\"foo\", \"password\":\"wK8iVkRO\"}",
		Handler: func(response *httptest.ResponseRecorder) {
			assert.Equal(t, http.StatusOK, response.Code)
		},
	})
}
