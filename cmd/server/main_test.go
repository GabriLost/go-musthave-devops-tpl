package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetMetricHandlerHandler(t *testing.T) {
	tests := []struct {
		description  string
		requestURL   string
		expectedCode int
	}{
		{
			description:  "200 Success gauge number with dots",
			requestURL:   "/update/gauge/numberMetric/100500.000",
			expectedCode: 200,
		},
		{
			description:  "200 Success gauge number without dots",
			requestURL:   "/update/gauge/numberMetric/80",
			expectedCode: 200,
		},
		{
			description:  "400 Parse Error",
			requestURL:   "/update/gauge/stringMetric/aaa",
			expectedCode: 400,
		},
		{
			description:  "400 Parse Error",
			requestURL:   "/update/counter/PollCount/666",
			expectedCode: 400,
		},
		{
			description:  "400 No such metric",
			requestURL:   "/update/wrong/doSomeThingElse/123",
			expectedCode: 400,
		},
		{
			description:  "400 short uri on update",
			requestURL:   "/update/shortURI/doSomeThingElse",
			expectedCode: 400,
		},
	}
	for _, url := range tests {
		t.Run(url.description, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, url.requestURL, nil)
			recorder := httptest.NewRecorder()
			handlerFunc := http.HandlerFunc(GetMetricHandler)
			handlerFunc.ServeHTTP(recorder, request)
			result := recorder.Result()
			if result.StatusCode != url.expectedCode {
				t.Errorf("Expected status code %d, but got %d", url.expectedCode, recorder.Code)
			}
			result.Body.Close()
		})
	}
}

func TestNotFound(t *testing.T) {
	tests := []struct {
		description  string
		requestURL   string
		expectedCode int
	}{
		{
			description:  "404 Trying do something else on post 1",
			requestURL:   "/update/somethingElse",
			expectedCode: 404,
		},
		{
			description:  "404 Trying do something else on post 2",
			requestURL:   "/value",
			expectedCode: 404,
		},
	}
	for _, url := range tests {
		t.Run(url.description, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, url.requestURL, nil)
			recorder := httptest.NewRecorder()
			handlerFunc := http.HandlerFunc(NotFound)
			handlerFunc.ServeHTTP(recorder, request)
			result := recorder.Result()
			if result.StatusCode != url.expectedCode {
				t.Errorf("Expected status code %d, but got %d", url.expectedCode, recorder.Code)
			}
			result.Body.Close()
		})
	}
}

func TestGetAllHandler(t *testing.T) {
	tests := []struct {
		name string
		url  string
	}{
		{
			name: "One",
			url:  "/",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			request := httptest.NewRequest(http.MethodGet, test.url, nil)
			handlerFunc := http.HandlerFunc(GetAllHandler)
			handlerFunc.ServeHTTP(recorder, request)
			result := recorder.Result()

			if result.Body == http.NoBody {
				t.Error("Body is empty")
			}

			if result.StatusCode != http.StatusOK {
				t.Errorf("StatusCode must be %d, but got %d", http.StatusOK, result.StatusCode)
			}

			result.Body.Close()
		})
	}
}
