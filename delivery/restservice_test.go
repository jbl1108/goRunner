package delivery

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/jbl1108/github.com/jbl1108/goSecret/usecases/datamodel"
)

// mockKeyValInputPort implements input.KeyValInputPort for testing.
type mockKeyValInputPort struct {
	getKeyResult []byte
	getKeyErr    error
	setKeyErr    error
	lastSetMsg   datamodel.Message
}

func (m *mockKeyValInputPort) GetKey(key string) ([]byte, error) {
	return m.getKeyResult, m.getKeyErr
}

func (m *mockKeyValInputPort) SetKey(message datamodel.Message) error {
	m.lastSetMsg = message
	return m.setKeyErr
}

// Test_handleGetKey_Success verifies successful key retrieval.
func Test_handleGetKey_Success(t *testing.T) {
	mock := &mockKeyValInputPort{
		getKeyResult: []byte("test-value"),
	}
	svc := NewKeyValueRestService("", mock)
	mux := http.NewServeMux()
	svc.RegisterRoutes(mux)

	req := httptest.NewRequest("GET", "/key/topic1/mykey", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
	if string(w.Body.Bytes()) != "test-value" {
		t.Errorf("expected body 'test-value', got %q", w.Body.String())
	}
	if ct := w.Header().Get("Content-Type"); ct != "application/octet-stream" {
		t.Errorf("expected Content-Type application/octet-stream, got %q", ct)
	}
}

// Test_handleGetKey_NotFound verifies error handling when key not found.
func Test_handleGetKey_NotFound(t *testing.T) {
	mock := &mockKeyValInputPort{
		getKeyErr: errors.New("key not found"),
	}
	svc := NewKeyValueRestService("", mock)
	mux := http.NewServeMux()
	svc.RegisterRoutes(mux)

	req := httptest.NewRequest("GET", "/key/topic1/missing", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", w.Code)
	}
	if !strings.Contains(w.Body.String(), "key not found") {
		t.Errorf("expected error message in body, got %q", w.Body.String())
	}
}

// Test_handleSetKey_Success verifies successful key creation.
func Test_handleSetKey_Success(t *testing.T) {
	mock := &mockKeyValInputPort{}
	svc := NewKeyValueRestService("", mock)
	mux := http.NewServeMux()
	svc.RegisterRoutes(mux)

	req := httptest.NewRequest("POST", "/key/topic1/mykey", strings.NewReader("test-data"))
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected 201 Created, got %d", w.Code)
	}
	if mock.lastSetMsg.Topic != "topic1" {
		t.Errorf("expected topic topic1, got %q", mock.lastSetMsg.Topic)
	}
	if mock.lastSetMsg.Data.Key != "mykey" {
		t.Errorf("expected key mykey, got %q", mock.lastSetMsg.Data.Key)
	}
	if mock.lastSetMsg.Data.Value != "test-data" {
		t.Errorf("expected value test-data, got %q", mock.lastSetMsg.Data.Value)
	}
}

// Test_handleSetKey_EmptyBody verifies handling of empty request body.
func Test_handleSetKey_EmptyBody(t *testing.T) {
	mock := &mockKeyValInputPort{}
	svc := NewKeyValueRestService("", mock)
	mux := http.NewServeMux()
	svc.RegisterRoutes(mux)

	req := httptest.NewRequest("POST", "/key/topic1/key1", strings.NewReader(""))
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected 201, got %d", w.Code)
	}
	if mock.lastSetMsg.Data.Value != "" {
		t.Errorf("expected empty value, got %q", mock.lastSetMsg.Data.Value)
	}
}

// Test_handleSetKey_LargeBody verifies handling of large request body.
func Test_handleSetKey_LargeBody(t *testing.T) {
	mock := &mockKeyValInputPort{}
	svc := NewKeyValueRestService("", mock)
	mux := http.NewServeMux()
	svc.RegisterRoutes(mux)

	largeData := strings.Repeat("x", 10000)
	req := httptest.NewRequest("POST", "/key/topic1/key1", strings.NewReader(largeData))
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected 201, got %d", w.Code)
	}
	if mock.lastSetMsg.Data.Value != largeData {
		t.Errorf("expected large data to be stored")
	}
}

// Test_handleSetKey_UsecaseError verifies error propagation from usecase.
func Test_handleSetKey_UsecaseError(t *testing.T) {
	mock := &mockKeyValInputPort{
		setKeyErr: errors.New("storage error"),
	}
	svc := NewKeyValueRestService("", mock)
	mux := http.NewServeMux()
	svc.RegisterRoutes(mux)

	req := httptest.NewRequest("POST", "/key/topic1/key1", strings.NewReader("data"))
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", w.Code)
	}
	if !strings.Contains(w.Body.String(), "storage error") {
		t.Errorf("expected storage error in response, got %q", w.Body.String())
	}
}

// Test_handleSetKey_BodyReadError verifies handling of malformed request bodies.
func Test_handleSetKey_BodyReadError(t *testing.T) {
	mock := &mockKeyValInputPort{}
	svc := NewKeyValueRestService("", mock)
	mux := http.NewServeMux()
	svc.RegisterRoutes(mux)

	// Create a request with a body that will error on read
	body := &errorReader{}
	req := httptest.NewRequest("POST", "/key/topic1/key1", body)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
	if !strings.Contains(w.Body.String(), "Error reading request body") {
		t.Errorf("expected body read error message, got %q", w.Body.String())
	}
}

// Test_HealthEndpoint verifies the health check endpoint.
func Test_HealthEndpoint(t *testing.T) {
	mock := &mockKeyValInputPort{}
	svc := NewKeyValueRestService("", mock)
	mux := http.NewServeMux()
	svc.RegisterRoutes(mux)

	req := httptest.NewRequest("GET", "/health/", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
	if string(w.Body.Bytes()) != "OK" {
		t.Errorf("expected body 'OK', got %q", w.Body.String())
	}
}

// Test_RootEndpoint verifies the welcome page.
func Test_RootEndpoint(t *testing.T) {
	mock := &mockKeyValInputPort{}
	svc := NewKeyValueRestService("", mock)
	mux := http.NewServeMux()
	svc.RegisterRoutes(mux)

	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
	body := w.Body.String()
	if !strings.Contains(body, "Welcome to the KeyValue REST Service") {
		t.Errorf("expected welcome message, got %q", body)
	}
	if !strings.Contains(body, "GET /key/{topic}/{key}") {
		t.Errorf("expected GET endpoint description, got %q", body)
	}
}

// errorReader implements io.Reader that always fails.
type errorReader struct{}

func (er *errorReader) Read(p []byte) (n int, err error) {
	return 0, io.ErrUnexpectedEOF
}
