package routes

import (
	"github.com/simonwep/genisis/core"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"testing"
)

func setup(t *testing.T) string {
	core.DropDatabase()
	var token string

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
			token = response.Header().Get("Authorization")[7:]
		},
	})

	return token
}

func TestGetAllData(t *testing.T) {
	token := setup(t)

	tryAuthorizedGet(AuthorizedGetConfig{
		Url:   "/data",
		Token: token,
		Handler: func(response *httptest.ResponseRecorder) {
			assert.Equal(t, 200, response.Code)
			assert.Equal(t, "{}", response.Body.String())
		},
	})

	tryAuthorizedPut(AuthorizedPostConfig{
		Url:   "/data/bar",
		Body:  "{\"hello\": \"world!\"}",
		Token: token,
		Handler: func(response *httptest.ResponseRecorder) {
			assert.Equal(t, 200, response.Code)
		},
	})

	tryAuthorizedPut(AuthorizedPostConfig{
		Url:   "/data/foo",
		Body:  "{\"test\": 200}",
		Token: token,
		Handler: func(response *httptest.ResponseRecorder) {
			assert.Equal(t, 200, response.Code)
		},
	})

	tryAuthorizedGet(AuthorizedGetConfig{
		Url:   "/data",
		Token: token,
		Handler: func(response *httptest.ResponseRecorder) {
			assert.Equal(t, 200, response.Code)
			assert.Equal(t, "{\"bar\":{\"hello\":\"world!\"},\"foo\":{\"test\":200}}", response.Body.String())
		},
	})
}

func TestSingleDataSet(t *testing.T) {
	token := setup(t)

	tryAuthorizedPut(AuthorizedPostConfig{
		Url:   "/data/bar",
		Body:  "{\"hello\": \"world!\"}",
		Token: token,
		Handler: func(response *httptest.ResponseRecorder) {
			assert.Equal(t, 200, response.Code)
		},
	})

	tryAuthorizedGet(AuthorizedGetConfig{
		Url:   "/data/bar",
		Token: token,
		Handler: func(response *httptest.ResponseRecorder) {
			assert.Equal(t, 200, response.Code)
			assert.Equal(t, "{\"hello\":\"world!\"}", response.Body.String())
		},
	})
}
