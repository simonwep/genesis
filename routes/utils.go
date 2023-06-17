package routes

import (
	"net/http"
	"net/http/httptest"
	"strings"
)

type Config struct {
	Method  string
	Url     string
	Body    string
	Handler func(*httptest.ResponseRecorder)
}

func tryRoute(config Config) {
	router := SetupRoutes()

	response := httptest.NewRecorder()
	request, _ := http.NewRequest(config.Method, config.Url, strings.NewReader(config.Body))
	request.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(response, request)
	config.Handler(response)
}
