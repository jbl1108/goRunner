package services

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jbl1108/goRunner/usecases/datamodel"
)

func TestNewKeyValueRepository(t *testing.T) {
	service := NewKeyValueRepository("http://localhost:8080")

	if service.requestUrl != "http://localhost:8080" {
		t.Errorf("expected requestUrl 'http://localhost:8080', got '%s'", service.requestUrl)
	}
}

func TestKeyValueStorageService_GetAllTrainings(t *testing.T) {
	// Create a test server that returns mock trainings
	trainings := []datamodel.Training{
		{Uid: "uid1", Week: 1, Dayofweek: 1, Activity: "run", Done: false},
		{Uid: "uid2", Week: 2, Dayofweek: 2, Activity: "swim", Done: true},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(trainings)
	}))
	defer server.Close()

	service := NewKeyValueRepository(server.URL)

	result, err := service.GetAllTrainings()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if len(result) != 2 {
		t.Errorf("expected 2 trainings, got %d", len(result))
	}

	if result[0].Uid != "uid1" || result[1].Uid != "uid2" {
		t.Error("trainings not returned correctly")
	}
}

func TestKeyValueStorageService_GetAllTrainings_HTTPError(t *testing.T) {
	// Create a test server that returns an error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	service := NewKeyValueRepository(server.URL)

	_, err := service.GetAllTrainings()
	if err == nil {
		t.Error("expected error from HTTP request")
	}
}

func TestKeyValueStorageService_GetAllTrainings_InvalidJSON(t *testing.T) {
	// Create a test server that returns invalid JSON
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("invalid json"))
	}))
	defer server.Close()

	service := NewKeyValueRepository(server.URL)

	_, err := service.GetAllTrainings()
	if err == nil {
		t.Error("expected error from invalid JSON")
	}
}
