package routes

import (
	"github.com/simonwep/genisis/core"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func setup(t *testing.T) string {
	core.DropDatabase()
	var token string

	tryUnauthorizedPost(UnauthorizedBodyConfig{
		Url:  "/register",
		Body: "{\"user\": \"foo\", \"password\": \"test\"}",
		Handler: func(response *httptest.ResponseRecorder) {
			assert.Equal(t, http.StatusCreated, response.Code)
			assert.Equal(t, 0, response.Body.Len())
		},
	})

	tryUnauthorizedPost(UnauthorizedBodyConfig{
		Url:  "/login",
		Body: "{\"user\": \"foo\", \"password\": \"test\"}",
		Handler: func(response *httptest.ResponseRecorder) {
			assert.Equal(t, http.StatusOK, response.Code)
			assert.Equal(t, 0, response.Body.Len())
			token = response.Header().Get("Authorization")[7:]
		},
	})

	return token
}

func TestGetAllData(t *testing.T) {
	token := setup(t)

	tryAuthorizedGet(AuthorizedConfig{
		Url:   "/data",
		Token: token,
		Handler: func(response *httptest.ResponseRecorder) {
			assert.Equal(t, http.StatusOK, response.Code)
			assert.Equal(t, "{}", response.Body.String())
		},
	})

	tryAuthorizedPut(AuthorizedBodyConfig{
		Url:   "/data/bar",
		Body:  "{\"hello\": \"world!\"}",
		Token: token,
		Handler: func(response *httptest.ResponseRecorder) {
			assert.Equal(t, http.StatusOK, response.Code)
		},
	})

	tryAuthorizedPut(AuthorizedBodyConfig{
		Url:   "/data/foo",
		Body:  "{\"test\": 200}",
		Token: token,
		Handler: func(response *httptest.ResponseRecorder) {
			assert.Equal(t, http.StatusOK, response.Code)
		},
	})

	tryAuthorizedGet(AuthorizedConfig{
		Url:   "/data",
		Token: token,
		Handler: func(response *httptest.ResponseRecorder) {
			assert.Equal(t, http.StatusOK, response.Code)
			assert.Equal(t, "{\"bar\":{\"hello\":\"world!\"},\"foo\":{\"test\":200}}", response.Body.String())
		},
	})
}

func TestSingleObject(t *testing.T) {
	token := setup(t)

	tryAuthorizedPut(AuthorizedBodyConfig{
		Url:   "/data/bar",
		Body:  "{\"hello\": \"world!\"}",
		Token: token,
		Handler: func(response *httptest.ResponseRecorder) {
			assert.Equal(t, http.StatusOK, response.Code)
		},
	})

	tryAuthorizedGet(AuthorizedConfig{
		Url:   "/data/bar",
		Token: token,
		Handler: func(response *httptest.ResponseRecorder) {
			assert.Equal(t, http.StatusOK, response.Code)
			assert.Equal(t, "{\"hello\":\"world!\"}", response.Body.String())
		},
	})
}

func TestSingleArray(t *testing.T) {
	token := setup(t)

	tryAuthorizedPut(AuthorizedBodyConfig{
		Url:   "/data/bar",
		Body:  "[1, 2, 3, 4, 5, 6]",
		Token: token,
		Handler: func(response *httptest.ResponseRecorder) {
			assert.Equal(t, http.StatusOK, response.Code)
		},
	})

	tryAuthorizedGet(AuthorizedConfig{
		Url:   "/data/bar",
		Token: token,
		Handler: func(response *httptest.ResponseRecorder) {
			assert.Equal(t, http.StatusOK, response.Code)
			assert.Equal(t, "[1,2,3,4,5,6]", response.Body.String())
		},
	})
}

func TestInvalidKeys(t *testing.T) {
	token := setup(t)

	tryAuthorizedPut(AuthorizedBodyConfig{
		Url:   "/data/HsRLrSgCFylAK77aJmvRon0ubjXjzPFtd", // too long
		Body:  "{\"hello\": \"world!\"}",
		Token: token,
		Handler: func(response *httptest.ResponseRecorder) {
			assert.Equal(t, http.StatusBadRequest, response.Code)
		},
	})

	tryAuthorizedPut(AuthorizedBodyConfig{
		Url:   "/data/🦧", // invalid characters
		Body:  "{\"hello\": \"world!\"}",
		Token: token,
		Handler: func(response *httptest.ResponseRecorder) {
			assert.Equal(t, http.StatusBadRequest, response.Code)
		},
	})
}

func TestDeleteData(t *testing.T) {
	token := setup(t)

	tryAuthorizedPut(AuthorizedBodyConfig{
		Url:   "/data/bar",
		Body:  "{\"hello\": \"world!\"}",
		Token: token,
		Handler: func(response *httptest.ResponseRecorder) {
			assert.Equal(t, http.StatusOK, response.Code)
		},
	})

	tryAuthorizedGet(AuthorizedConfig{
		Url:   "/data/bar",
		Token: token,
		Handler: func(response *httptest.ResponseRecorder) {
			assert.Equal(t, http.StatusOK, response.Code)
			assert.Equal(t, "{\"hello\":\"world!\"}", response.Body.String())
		},
	})

	tryAuthorizedDelete(AuthorizedConfig{
		Url:   "/data/bar",
		Token: token,
		Handler: func(response *httptest.ResponseRecorder) {
			assert.Equal(t, http.StatusOK, response.Code)
		},
	})

	tryAuthorizedGet(AuthorizedConfig{
		Url:   "/data/bar",
		Token: token,
		Handler: func(response *httptest.ResponseRecorder) {
			assert.Equal(t, http.StatusNotFound, response.Code)
		},
	})
}

func TestTooBig(t *testing.T) {
	token := setup(t)

	body1 := strings.Repeat("{\"hello\": \"world!\"},", 10)
	tryAuthorizedPut(AuthorizedBodyConfig{
		Url:   "/data/bar",
		Body:  "[" + body1[:len(body1)-1] + "]",
		Token: token,
		Handler: func(response *httptest.ResponseRecorder) {
			assert.Equal(t, http.StatusOK, response.Code)
		},
	})

	body2 := strings.Repeat("{\"hello\": \"world!\"},", 50)
	tryAuthorizedPut(AuthorizedBodyConfig{
		Url:   "/data/bar",
		Body:  "[" + body2[:len(body2)-1] + "]",
		Token: token,
		Handler: func(response *httptest.ResponseRecorder) {
			assert.Equal(t, http.StatusRequestEntityTooLarge, response.Code)
		},
	})
}

func TestTooMany(t *testing.T) {
	token := setup(t)

	tryAuthorizedPut(AuthorizedBodyConfig{
		Url:   "/data/bar1",
		Body:  "{\"hello\": \"world!\"}",
		Token: token,
		Handler: func(response *httptest.ResponseRecorder) {
			assert.Equal(t, http.StatusOK, response.Code)
		},
	})

	tryAuthorizedPut(AuthorizedBodyConfig{
		Url:   "/data/bar2",
		Body:  "{\"hello\": \"world!\"}",
		Token: token,
		Handler: func(response *httptest.ResponseRecorder) {
			assert.Equal(t, http.StatusOK, response.Code)
		},
	})

	tryAuthorizedPut(AuthorizedBodyConfig{
		Url:   "/data/bar3",
		Body:  "{\"hello\": \"world!\"}",
		Token: token,
		Handler: func(response *httptest.ResponseRecorder) {
			assert.Equal(t, http.StatusOK, response.Code)
		},
	})

	// this one should fail
	tryAuthorizedPut(AuthorizedBodyConfig{
		Url:   "/data/bar4",
		Body:  "{\"hello\": \"world!\"}",
		Token: token,
		Handler: func(response *httptest.ResponseRecorder) {
			assert.Equal(t, http.StatusForbidden, response.Code)
		},
	})

	// make it succeed
	tryAuthorizedDelete(AuthorizedConfig{
		Url:   "/data/bar3",
		Token: token,
		Handler: func(response *httptest.ResponseRecorder) {
			assert.Equal(t, http.StatusOK, response.Code)
		},
	})

	tryAuthorizedPut(AuthorizedBodyConfig{
		Url:   "/data/bar4",
		Body:  "{\"hello\": \"world!\"}",
		Token: token,
		Handler: func(response *httptest.ResponseRecorder) {
			assert.Equal(t, http.StatusOK, response.Code)
		},
	})
}
