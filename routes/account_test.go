package routes

import (
	"github.com/simonwep/genisis/core"
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
			token = response.Header().Get("Authorization")[7:]
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
			assert.Equal(t, http.StatusUnauthorized, response.Code)
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
			token = response.Header().Get("Authorization")[7:]
		},
	})
}
