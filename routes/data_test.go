package routes

import (
	"encoding/json"
	"github.com/simonwep/genisis/core"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func setup(t *testing.T) string {
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

	return token
}

func TestGetAllData(t *testing.T) {
	token := setup(t)

	tryAuthorizedGet("/data", AuthorizedConfig{
		Token: token,
		Handler: func(response *httptest.ResponseRecorder) {
			assert.Equal(t, http.StatusOK, response.Code)
			assert.Equal(t, "{}", response.Body.String())
		},
	})

	tryAuthorizedPost("/data/bar", AuthorizedBodyConfig{
		Body:  "{\"hello\": \"world!\"}",
		Token: token,
		Handler: func(response *httptest.ResponseRecorder) {
			assert.Equal(t, http.StatusOK, response.Code)
		},
	})

	tryAuthorizedPost("/data/foo", AuthorizedBodyConfig{
		Body:  "{\"test\": 200}",
		Token: token,
		Handler: func(response *httptest.ResponseRecorder) {
			assert.Equal(t, http.StatusOK, response.Code)
		},
	})

	tryAuthorizedGet("/data", AuthorizedConfig{
		Token: token,
		Handler: func(response *httptest.ResponseRecorder) {
			assert.Equal(t, http.StatusOK, response.Code)
			assert.Equal(t, "{\"bar\":{\"hello\":\"world!\"},\"foo\":{\"test\":200}}", response.Body.String())
		},
	})
}

func TestInvalidJSON(t *testing.T) {
	token := setup(t)

	tryAuthorizedPost("/data/bar", AuthorizedBodyConfig{
		Body:  "{\"hello\": \"world!\"",
		Token: token,
		Handler: func(response *httptest.ResponseRecorder) {
			assert.Equal(t, http.StatusBadRequest, response.Code)
		},
	})
}

func TestSingleObject(t *testing.T) {
	token := setup(t)

	tryAuthorizedPost("/data/bar", AuthorizedBodyConfig{
		Body:  "{\"hello\": 1e5}",
		Token: token,
		Handler: func(response *httptest.ResponseRecorder) {
			assert.Equal(t, http.StatusOK, response.Code)
		},
	})

	tryAuthorizedGet("/data/bar", AuthorizedConfig{
		Token: token,
		Handler: func(response *httptest.ResponseRecorder) {
			assert.Equal(t, http.StatusOK, response.Code)
			assert.Equal(t, "{\"hello\":1e5}", response.Body.String())
		},
	})
}

func TestEmpty(t *testing.T) {
	token := setup(t)

	tryAuthorizedGet("/data/bar", AuthorizedConfig{
		Token: token,
		Handler: func(response *httptest.ResponseRecorder) {
			assert.Equal(t, http.StatusNoContent, response.Code)
		},
	})
}

func TestSingleArray(t *testing.T) {
	token := setup(t)

	tryAuthorizedPost("/data/bar", AuthorizedBodyConfig{
		Body:  "[1, 2, 3, 4, 5, 6]",
		Token: token,
		Handler: func(response *httptest.ResponseRecorder) {
			assert.Equal(t, http.StatusOK, response.Code)
		},
	})

	tryAuthorizedGet("/data/bar", AuthorizedConfig{
		Token: token,
		Handler: func(response *httptest.ResponseRecorder) {
			assert.Equal(t, http.StatusOK, response.Code)
			assert.Equal(t, "[1,2,3,4,5,6]", response.Body.String())
		},
	})
}

func TestInvalidKeys(t *testing.T) {
	token := setup(t)

	tryAuthorizedPost("/data/HsRLrSgCFylAK77aJmvRon0ubjXjzPFtd", AuthorizedBodyConfig{ // too long
		Body:  "{\"hello\": \"world!\"}",
		Token: token,
		Handler: func(response *httptest.ResponseRecorder) {
			assert.Equal(t, http.StatusBadRequest, response.Code)
		},
	})

	tryAuthorizedPost("/data/🦧", AuthorizedBodyConfig{ // invalid characters
		Body:  "{\"hello\": \"world!\"}",
		Token: token,
		Handler: func(response *httptest.ResponseRecorder) {
			assert.Equal(t, http.StatusBadRequest, response.Code)
		},
	})
}

func TestDeleteData(t *testing.T) {
	token := setup(t)

	tryAuthorizedPost("/data/bar", AuthorizedBodyConfig{
		Body:  "{\"hello\": \"world!\"}",
		Token: token,
		Handler: func(response *httptest.ResponseRecorder) {
			assert.Equal(t, http.StatusOK, response.Code)
		},
	})

	tryAuthorizedGet("/data/bar", AuthorizedConfig{
		Token: token,
		Handler: func(response *httptest.ResponseRecorder) {
			assert.Equal(t, http.StatusOK, response.Code)
			assert.Equal(t, "{\"hello\":\"world!\"}", response.Body.String())
		},
	})

	tryAuthorizedDelete("/data/bar", AuthorizedConfig{
		Token: token,
		Handler: func(response *httptest.ResponseRecorder) {
			assert.Equal(t, http.StatusOK, response.Code)
		},
	})

	tryAuthorizedGet("/data/bar", AuthorizedConfig{
		Token: token,
		Handler: func(response *httptest.ResponseRecorder) {
			assert.Equal(t, http.StatusNoContent, response.Code)
		},
	})
}

func TestTooBig(t *testing.T) {
	token := setup(t)

	body1 := strings.Repeat("{\"hello\": \"world!\"},", 10)
	tryAuthorizedPost("/data/bar", AuthorizedBodyConfig{
		Body:  "[" + body1[:len(body1)-1] + "]",
		Token: token,
		Handler: func(response *httptest.ResponseRecorder) {
			assert.Equal(t, http.StatusOK, response.Code)
		},
	})

	body2 := strings.Repeat("{\"hello\": \"world!\"},", 50)
	tryAuthorizedPost("/data/bar", AuthorizedBodyConfig{
		Body:  "[" + body2[:len(body2)-1] + "]",
		Token: token,
		Handler: func(response *httptest.ResponseRecorder) {
			assert.Equal(t, http.StatusRequestEntityTooLarge, response.Code)
		},
	})
}

func TestTooMany(t *testing.T) {
	token := setup(t)

	tryAuthorizedPost("/data/bar1", AuthorizedBodyConfig{
		Body:  "{\"hello\": \"world!\"}",
		Token: token,
		Handler: func(response *httptest.ResponseRecorder) {
			assert.Equal(t, http.StatusOK, response.Code)
		},
	})

	tryAuthorizedPost("/data/bar2", AuthorizedBodyConfig{
		Body:  "{\"hello\": \"world!\"}",
		Token: token,
		Handler: func(response *httptest.ResponseRecorder) {
			assert.Equal(t, http.StatusOK, response.Code)
		},
	})

	tryAuthorizedPost("/data/bar3", AuthorizedBodyConfig{
		Body:  "{\"hello\": \"world!\"}",
		Token: token,
		Handler: func(response *httptest.ResponseRecorder) {
			assert.Equal(t, http.StatusOK, response.Code)
		},
	})

	// this one should fail
	tryAuthorizedPost("/data/bar4", AuthorizedBodyConfig{
		Body:  "{\"hello\": \"world!\"}",
		Token: token,
		Handler: func(response *httptest.ResponseRecorder) {
			assert.Equal(t, http.StatusForbidden, response.Code)
		},
	})

	// make it succeed
	tryAuthorizedDelete("/data/bar3", AuthorizedConfig{
		Token: token,
		Handler: func(response *httptest.ResponseRecorder) {
			assert.Equal(t, http.StatusOK, response.Code)
		},
	})

	tryAuthorizedPost("/data/bar4", AuthorizedBodyConfig{
		Body:  "{\"hello\": \"world!\"}",
		Token: token,
		Handler: func(response *httptest.ResponseRecorder) {
			assert.Equal(t, http.StatusOK, response.Code)
		},
	})
}
