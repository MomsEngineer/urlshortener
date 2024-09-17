package app

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

type Database interface {
	SaveLink(id, link string)
	GetLink(id string) (string, bool)
}

type MockDB struct{}

func (m *MockDB) SaveLink(id, link string) {
}

func (m *MockDB) GetLink(id string) (link string, exists bool) {
	if id == "abc123" {
		link, exists = "https://example.com", true
	} else {
		link, exists = "", false
	}
	return
}

func TestHandler(t *testing.T) {
	mockDB := new(MockDB)
	handler := Handler(mockDB)

	tests := []struct {
		name           string
		method         string
		url            string
		body           []byte
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "POST request",
			method:         http.MethodPost,
			url:            "/",
			body:           []byte("https://example.com"),
			expectedStatus: http.StatusCreated,
			expectedBody:   "http://localhost:8080/",
		},
		{
			name:           "Successful GET request",
			method:         http.MethodGet,
			url:            "/abc123",
			body:           nil,
			expectedStatus: http.StatusTemporaryRedirect,
			expectedBody:   "",
		},
		{
			name:           "Failed GET request",
			method:         http.MethodGet,
			url:            "/abc",
			body:           nil,
			expectedStatus: http.StatusNotFound,
			expectedBody:   "",
		},
		{
			name:           "Invalid method",
			method:         http.MethodPut,
			url:            "/",
			body:           nil,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Only GET and POST requests are allowed!",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.url, bytes.NewBuffer(tt.body))
			rr := httptest.NewRecorder()

			handler.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)
			assert.Contains(t, rr.Body.String(), tt.expectedBody)
		})
	}
}
