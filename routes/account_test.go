package routes

import (
	"encoding/json"
	"github.com/simonwep/genesis/core"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUpdatePassword(t *testing.T) {
	core.ResetDatabase()

	var token string
	tryUnauthorizedPost("/login", UnauthorizedBodyConfig{
		Body: "{\"user\": \"foo\", \"password\": \"hgEiPCZP\"}",
		Handler: func(response *httptest.ResponseRecorder) {
			assert.Equal(t, http.StatusOK, response.Code)
			data := LoginResponse{}
			json.Unmarshal(response.Body.Bytes(), &data)
			token = data.Token
		},
	})

	tryAuthorizedPost("/account/update", AuthorizedBodyConfig{
		Token: token,
		Body:  "{\"currentPassword\": \"hgEiPCZP\",\"newPassword\": \"hg235ZP\"}",
		Handler: func(response *httptest.ResponseRecorder) {
			assert.Equal(t, http.StatusOK, response.Code)
		},
	})

	tryAuthorizedPost("/account/update", AuthorizedBodyConfig{
		Token: token,
		Body:  "{\"currentPassword\": \"hgEiPCZP\",\"newPassword\": \"hg235ZP\"}",
		Handler: func(response *httptest.ResponseRecorder) {
			assert.Equal(t, http.StatusBadRequest, response.Code)
		},
	})

	tryUnauthorizedPost("/login", UnauthorizedBodyConfig{
		Body: "{\"user\": \"foo\", \"password\": \"hgEiPCZP\"}",
		Handler: func(response *httptest.ResponseRecorder) {
			assert.Equal(t, http.StatusUnauthorized, response.Code)
		},
	})

	tryUnauthorizedPost("/login", UnauthorizedBodyConfig{
		Body: "{\"user\": \"foo\", \"password\": \"hg235ZP\"}",
		Handler: func(response *httptest.ResponseRecorder) {
			assert.Equal(t, http.StatusOK, response.Code)
		},
	})
}
