package server

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestMetricHandler(t *testing.T) {
	type want struct {
		code int
		body []string
	}
	r := Router()
	server := httptest.NewServer(r)
	defer server.Close()

	tests := []struct {
		description string
		requestURL  string
		method      string
		expected    want
	}{
		{
			description: "200 Success gauge number with dots",
			requestURL:  "/update/gauge/numberMetric/100500.000",
			method:      http.MethodPost,
			expected:    want{code: 200},
		},
		{
			description: "200 Success gauge number without dots",
			requestURL:  "/update/gauge/numberMetric/80",
			method:      http.MethodPost,
			expected:    want{code: 200},
		},
		{
			description: "200 Success Counter",
			requestURL:  "/update/counter/PollCount/5",
			method:      http.MethodPost,
			expected:    want{code: 200},
		},
		{
			description: "200 Get Counter",
			requestURL:  "/value/counter/PollCount",
			method:      http.MethodGet,
			expected: want{
				code: 200,
				body: []string{"5"},
			},
		},
		{
			description: "200 Update Counter again",
			requestURL:  "/update/counter/PollCount/1",
			method:      http.MethodPost,
			expected:    want{code: 200},
		},
		{
			description: "200 Get Counter +1",
			requestURL:  "/value/counter/PollCount",
			method:      http.MethodGet,
			expected: want{
				code: 200,
				body: []string{"6"},
			},
		},
		{
			description: "400 Parse Error",
			requestURL:  "/update/gauge/stringMetric/aaa",
			method:      http.MethodPost,
			expected:    want{code: http.StatusBadRequest},
		},
		{
			description: "400 Parse Error",
			requestURL:  "/update/counter/PollCount/665g6",
			method:      http.MethodPost,
			expected:    want{code: http.StatusBadRequest},
		},
		{
			description: "501 No such metric",
			requestURL:  "/update/wrong/doSomeThingElse/123",
			method:      http.MethodPost,
			expected:    want{code: http.StatusNotImplemented},
		},
		{
			description: "400 short uri on update",
			requestURL:  "/update/shortURI/doSomeThingElse",
			method:      http.MethodPost,
			expected:    want{code: http.StatusBadRequest},
		},
		{
			description: "get unknown gauge",
			method:      http.MethodGet,
			requestURL:  "/value/gauge/lol",
			expected:    want{code: http.StatusNotFound},
		},
		{
			description: "get unknown counter",
			method:      http.MethodGet,
			requestURL:  "/value/counter/lol",
			expected:    want{code: http.StatusNotFound},
		},
		{
			description: "400 get unknown type",
			method:      http.MethodGet,
			requestURL:  "/value/wrongType/name",
			expected:    want{code: http.StatusBadRequest},
		},
		{
			description: "501 update unknown type",
			method:      http.MethodPost,
			requestURL:  "/update/unknown/testCounter/100",
			expected:    want{code: http.StatusNotImplemented},
		},
	}
	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			resp, body := executeRequest(t, server, tt.method, tt.requestURL, "")
			defer resp.Body.Close()
			assert.Equal(t, tt.expected.code, resp.StatusCode)
			for _, s := range tt.expected.body {
				assert.Equal(t, body, s)
			}
			assert.Equal(t, tt.expected.code, resp.StatusCode)
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
			handlerFunc := http.HandlerFunc(NotFoundHandler)
			handlerFunc.ServeHTTP(recorder, request)
			result := recorder.Result()
			if result.StatusCode != url.expectedCode {
				t.Errorf("Expected status code %d, but got %d", url.expectedCode, recorder.Code)
			}
			err := result.Body.Close()
			if err != nil {
				return
			}
		})
	}
}

func TestGetAllHandler(t *testing.T) {
	tests := []struct {
		name string
		url  string
	}{
		{
			name: "Check body isn't empty",
			url:  "/",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			request := httptest.NewRequest(http.MethodGet, test.url, nil)
			handlerFunc := http.HandlerFunc(AllMetricsHandler)
			handlerFunc.ServeHTTP(recorder, request)
			result := recorder.Result()

			if result.Body == http.NoBody {
				t.Error("Body is empty")
			}

			if result.StatusCode != http.StatusOK {
				t.Errorf("StatusCode must be %d, but got %d", http.StatusOK, result.StatusCode)
			}

			err := result.Body.Close()
			if err != nil {
				return
			}
		})
	}
}

func TestMetricJsonPostHandler(t *testing.T) {
	type want struct {
		code int
		body []string
	}
	r := Router()
	server := httptest.NewServer(r)
	defer server.Close()

	tests := []struct {
		description string
		requestURL  string
		method      string
		body        string
		expected    want
	}{
		{
			description: "200 update json counter delta",
			method:      http.MethodPost,
			requestURL:  "/update/",
			body:        `{"id":"poll","type":"counter","delta":5}`,
			expected:    want{code: http.StatusOK},
		},
		{
			description: "200 update json counter delta again",
			method:      http.MethodPost,
			requestURL:  "/update/",
			body:        `{"id":"poll","type":"counter","delta":5}`,
			expected:    want{code: http.StatusOK},
		},
		{
			description: "200 get counter delta",
			method:      http.MethodPost,
			requestURL:  "/value/",
			body:        `{"id":"poll","type":"counter"}`,
			expected: want{
				code: http.StatusOK,
				body: []string{
					`{"id":"poll","type":"counter","delta":10}`,
				},
			},
		},
		{
			description: "200 update json gauge value",
			method:      http.MethodPost,
			requestURL:  "/update/",
			body:        `{"id":"alloc","type":"gauge","value":10000}`,
			expected:    want{code: http.StatusOK},
		},
		{
			description: "200 get gauge value",
			method:      http.MethodPost,
			requestURL:  "/value/",
			body:        `{"id":"alloc","type":"gauge"}`,
			expected: want{
				code: http.StatusOK,
				body: []string{
					`{"id":"alloc","type":"gauge","value":10000}`,
				},
			},
		},
		{
			description: "501 update json unknown type",
			method:      http.MethodPost,
			requestURL:  "/update/",
			body:        `{"id":"other","type":"other","value":1}`,
			expected:    want{code: http.StatusNotImplemented},
		},
	}
	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			resp, body := executeRequest(t, server, tt.method, tt.requestURL, tt.body)
			defer resp.Body.Close()
			assert.Equal(t, tt.expected.code, resp.StatusCode)
			for _, s := range tt.expected.body {
				assert.Contains(t, body, s)
			}
			assert.Equal(t, tt.expected.code, resp.StatusCode)
		})
	}
}

func executeRequest(t *testing.T, ts *httptest.Server, method string, query string, body string) (*http.Response, string) {
	reader := strings.NewReader(body)
	req, err := http.NewRequest(method, ts.URL+query, reader)
	req.Header.Add("Content-Type", "application/json")

	require.NoError(t, err)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	respBody, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)

	return resp, string(respBody)
}
