package routes

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
)

type UnauthorizedBodyConfig struct {
	Body    string
	Handler func(*httptest.ResponseRecorder)
}

type AuthorizedConfig struct {
	Token   string
	Cookie  string
	Handler func(*httptest.ResponseRecorder)
}

type AuthorizedBodyConfig struct {
	Body    string
	Token   string
	Cookie  string
	Handler func(*httptest.ResponseRecorder)
}

func tryRequest(url, method, body string, config AuthorizedConfig) {
	router := SetupRoutes()

	response := httptest.NewRecorder()
	request, _ := http.NewRequest(method, url, strings.NewReader(body))

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Content-Length", strconv.FormatInt(int64(len(body)), 10))
	request.Header.Set("Authorization", "Bearer "+config.Token)
	request.Header.Set("Cookie", config.Cookie)

	router.ServeHTTP(response, request)
	config.Handler(response)
}

func tryUnauthorizedPost(url string, config UnauthorizedBodyConfig) {
	tryRequest(url, "POST", config.Body, AuthorizedConfig{
		Token:   "",
		Cookie:  "",
		Handler: config.Handler,
	})
}

func tryAuthorizedPost(url string, config AuthorizedBodyConfig) {
	tryRequest(url, "POST", config.Body, AuthorizedConfig{
		Token:   config.Token,
		Cookie:  config.Cookie,
		Handler: config.Handler,
	})
}

func tryAuthorizedDelete(url string, config AuthorizedConfig) {
	tryRequest(url, "DELETE", "", config)
}

func tryAuthorizedGet(url string, config AuthorizedConfig) {
	tryRequest(url, "GET", "", config)
}
