package handlers

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"

	"tracker/models"
)

func TestDateHandler(t *testing.T) {
	os.Setenv("TEST_MODE", "true")
	defer os.Unsetenv("TEST_MODE")
	// Mock FetchDates function for testing
	originalFetchDateFunc := fetchDatesFunc
	defer func() { fetchDatesFunc = originalFetchDateFunc }()

	fetchDatesFunc = func(id string) (models.Date, error) {
		idNum, _ := strconv.Atoi(id)
		if idNum == 1 {
			return models.Date{
				Id: 1, Dates: []string{"2023-01-01"},
			}, nil
		}
		return models.Date{}, fmt.Errorf("error fetching dates")
	}
	tests := []struct {
		name               string
		method             string
		urlPath            string
		queryParams        string
		expectedStatusCode int
	}{
		{
			name:               "Valid Request",
			method:             http.MethodGet,
			urlPath:            "/dates",
			queryParams:        "?id=1",
			expectedStatusCode: http.StatusOK,
		},
		{
			name:               "Invalid Path",
			method:             http.MethodGet,
			urlPath:            "/invalid",
			expectedStatusCode: http.StatusNotFound,
		},
		{
			name:               "Invalid Method",
			method:             http.MethodPost,
			urlPath:            "/dates",
			expectedStatusCode: http.StatusMethodNotAllowed,
		},
		{
			name:               "Invalid ID - Out of Range",
			method:             http.MethodGet,
			urlPath:            "/dates",
			queryParams:        "?id=100",
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:               "Invalid ID - Non-numeric",
			method:             http.MethodGet,
			urlPath:            "/dates",
			queryParams:        "?id=abc",
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:               "Internal Server Error",
			method:             http.MethodGet,
			urlPath:            "/dates",
			queryParams:        "?id=2",
			expectedStatusCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Prepare the request
			req := httptest.NewRequest(tt.method, tt.urlPath+tt.queryParams, nil)
			w := httptest.NewRecorder()

			DateHandler(w, req)

			res := w.Result()
			if res.StatusCode != tt.expectedStatusCode {
				t.Errorf("expected status code %d, got %d", tt.expectedStatusCode, res.StatusCode)
			}
		})
	}
}

func TestLocationHandler(t *testing.T) {
	os.Setenv("TEST_MODE", "true")
	defer os.Unsetenv("TEST_MODE")
	// Mock FetchDates function for testing
	originalFetchLocationsFunc := fetchLocationsFunc
	defer func() { fetchLocationsFunc = originalFetchLocationsFunc }()

	fetchLocationsFunc = func(id string) (models.Location, error) {
		idNum, _ := strconv.Atoi(id)
		if idNum == 1 {
			return models.Location{
				ArtistId: 1, Locations: []string{"2023-01-01"},
				Date: "",
			}, nil
		}
		return models.Location{}, fmt.Errorf("error fetching dates")
	}
	tests := []struct {
		name               string
		method             string
		urlPath            string
		queryParams        string
		expectedStatusCode int
	}{
		{
			name:               "Valid Request",
			method:             http.MethodGet,
			urlPath:            "/locations",
			queryParams:        "?id=1",
			expectedStatusCode: http.StatusOK,
		},
		{
			name:               "Invalid Path",
			method:             http.MethodGet,
			urlPath:            "/invalid",
			expectedStatusCode: http.StatusNotFound,
		},
		{
			name:               "Invalid Method",
			method:             http.MethodPost,
			urlPath:            "/locations",
			expectedStatusCode: http.StatusMethodNotAllowed,
		},
		{
			name:               "Invalid ID - Out of Range",
			method:             http.MethodGet,
			urlPath:            "/locations",
			queryParams:        "?id=100",
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:               "Invalid ID - Non-numeric",
			method:             http.MethodGet,
			urlPath:            "/locations",
			queryParams:        "?id=abc",
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:               "Internal Server Error",
			method:             http.MethodGet,
			urlPath:            "/locations",
			queryParams:        "?id=2",
			expectedStatusCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Prepare the request
			req := httptest.NewRequest(tt.method, tt.urlPath+tt.queryParams, nil)
			w := httptest.NewRecorder()

			LocationHandler(w, req)

			res := w.Result()
			if res.StatusCode != tt.expectedStatusCode {
				t.Errorf("expected status code %d, got %d", tt.expectedStatusCode, res.StatusCode)
			}
		})
	}
}
