package routes

import (
	"net/http"
	"net/http/httptest"
	"strings"
)

type UnauthorizedBodyConfig struct {
	Url     string
	Body    string
	Handler func(*httptest.ResponseRecorder)
}

type AuthorizedConfig struct {
	Url     string
	Token   string
	Handler func(*httptest.ResponseRecorder)
}

type AuthorizedBodyConfig struct {
	Url     string
	Body    string
	Token   string
	Handler func(*httptest.ResponseRecorder)
}

func tryRequest(config AuthorizedConfig, method string, body string) {
	router := SetupRoutes()

	response := httptest.NewRecorder()
	request, _ := http.NewRequest(method, config.Url, strings.NewReader(body))

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+config.Token)

	router.ServeHTTP(response, request)
	config.Handler(response)
}

func tryUnauthorizedPost(config UnauthorizedBodyConfig) {
	tryRequest(AuthorizedConfig{
		Url:     config.Url,
		Token:   "",
		Handler: config.Handler,
	}, "POST", config.Body)
}

func tryAuthorizedPut(config AuthorizedBodyConfig) {
	tryRequest(AuthorizedConfig{
		Url:     config.Url,
		Token:   config.Token,
		Handler: config.Handler,
	}, "PUT", config.Body)
}

func tryAuthorizedDelete(config AuthorizedConfig) {
	tryRequest(config, "DELETE", "")
}

func tryAuthorizedGet(config AuthorizedConfig) {
	tryRequest(config, "GET", "")
}
