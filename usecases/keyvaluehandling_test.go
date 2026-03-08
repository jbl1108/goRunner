package usecases

import (
	"errors"
	"testing"

	"github.com/jbl1108/github.com/jbl1108/goSecret/usecases/datamodel"
)

// mockStorage implements output.KeyValueStorage for testing.
type mockStorage struct {
	openErr     error
	closeErr    error
	setErr      error
	getErr      error
	getResult   []byte
	lastSetKey  string
	lastSetData []byte
}

func (m *mockStorage) Open() error  { return m.openErr }
func (m *mockStorage) Close() error { return m.closeErr }
func (m *mockStorage) Get(key string) ([]byte, error) {
	return m.getResult, m.getErr
}
func (m *mockStorage) Set(key string, data []byte) error {
	m.lastSetKey = key
	m.lastSetData = data
	return m.setErr
}

// Test_SetKey_Success verifies successful key storage.
func Test_SetKey_Success(t *testing.T) {
	mock := &mockStorage{}
	uc := NewKeyValueHandling(mock)

	msg := datamodel.Message{
		Topic: "topic1",
		Data: datamodel.KeyValue{
			Key:   "key1",
			Value: "test-value",
		},
	}
	err := uc.SetKey(msg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if mock.lastSetKey != "topic1:key1" {
		t.Errorf("expected key topic1:key1, got %q", mock.lastSetKey)
	}
	if string(mock.lastSetData) != "test-value" {
		t.Errorf("expected data test-value, got %q", mock.lastSetData)
	}
}

// Test_SetKey_OpenError verifies error handling on storage open failure.
func Test_SetKey_OpenError(t *testing.T) {
	mock := &mockStorage{
		openErr: errors.New("connection failed"),
	}
	uc := NewKeyValueHandling(mock)

	msg := datamodel.Message{
		Topic: "topic1",
		Data:  datamodel.KeyValue{Key: "key1", Value: "data"},
	}
	err := uc.SetKey(msg)
	if err == nil || err.Error() != "connection failed" {
		t.Errorf("expected connection failed error, got %v", err)
	}
}

// Test_SetKey_WriteError verifies error handling on storage write failure.
func Test_SetKey_WriteError(t *testing.T) {
	mock := &mockStorage{
		setErr: errors.New("write failed"),
	}
	uc := NewKeyValueHandling(mock)

	msg := datamodel.Message{
		Topic: "topic1",
		Data:  datamodel.KeyValue{Key: "key1", Value: "data"},
	}
	err := uc.SetKey(msg)
	if err == nil || err.Error() != "write failed" {
		t.Errorf("expected write failed error, got %v", err)
	}
}

// Test_GetKey_Success verifies successful key retrieval.
func Test_GetKey_Success(t *testing.T) {
	mock := &mockStorage{
		getResult: []byte("test-data"),
	}
	uc := NewKeyValueHandling(mock)

	result, err := uc.GetKey("topic1:key1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(result) != "test-data" {
		t.Errorf("expected test-data, got %q", result)
	}
}

// Test_GetKey_NotFound verifies error handling when key is not found.
func Test_GetKey_NotFound(t *testing.T) {
	mock := &mockStorage{
		getErr: errors.New("key not found"),
	}
	uc := NewKeyValueHandling(mock)

	_, err := uc.GetKey("topic1:missing")
	if err == nil || err.Error() != "key not found" {
		t.Errorf("expected key not found error, got %v", err)
	}
}

// Test_GetKey_OpenError verifies error handling on storage open failure.
func Test_GetKey_OpenError(t *testing.T) {
	mock := &mockStorage{
		openErr: errors.New("connection lost"),
	}
	uc := NewKeyValueHandling(mock)

	_, err := uc.GetKey("topic1:key1")
	if err == nil || err.Error() != "connection lost" {
		t.Errorf("expected connection lost error, got %v", err)
	}
}
