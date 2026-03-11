package delivery

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jbl1108/goRunner/usecases/datamodel"
)

// Mock TrainingInputPort for testing
type mockTrainingInputPort struct {
	trainings   []datamodel.Training
	getAllErr   error
	getErr      error
	addErr      error
	updateErr   error
	deleteErr   error
	lastAdded   datamodel.Training
	lastUpdated datamodel.Training
	lastDeleted string
}

func (m *mockTrainingInputPort) GetAllTrainings() ([]datamodel.Training, error) {
	if m.getAllErr != nil {
		return nil, m.getAllErr
	}
	return m.trainings, nil
}

func (m *mockTrainingInputPort) GetTraining(uid string) (datamodel.Training, error) {
	if m.getErr != nil {
		return datamodel.Training{}, m.getErr
	}
	for _, t := range m.trainings {
		if t.Uid == uid {
			return t, nil
		}
	}
	return datamodel.Training{}, &mockError{msg: "training not found"}
}

func (m *mockTrainingInputPort) AddTraining(training datamodel.Training) error {
	if m.addErr != nil {
		return m.addErr
	}
	m.lastAdded = training
	m.trainings = append(m.trainings, training)
	return nil
}

func (m *mockTrainingInputPort) UpdateTraining(training datamodel.Training) error {
	if m.updateErr != nil {
		return m.updateErr
	}
	m.lastUpdated = training
	for i, t := range m.trainings {
		if t.Uid == training.Uid {
			m.trainings[i] = training
			return nil
		}
	}
	return &mockError{msg: "training not found"}
}

func (m *mockTrainingInputPort) DeleteTraining(uid string) error {
	if m.deleteErr != nil {
		return m.deleteErr
	}
	m.lastDeleted = uid
	for i, t := range m.trainings {
		if t.Uid == uid {
			m.trainings = append(m.trainings[:i], m.trainings[i+1:]...)
			return nil
		}
	}
	return &mockError{msg: "training not found"}
}

type mockError struct {
	msg string
}

func (m *mockError) Error() string {
	return m.msg
}

func TestNewTrainingRestService(t *testing.T) {
	mockPort := &mockTrainingInputPort{}
	service := NewTrainingRestService(":8080", mockPort)

	if service.address != ":8080" {
		t.Errorf("expected address ':8080', got '%s'", service.address)
	}

	if service.trainingHandlingUsecase != mockPort {
		t.Error("trainingHandlingUsecase not set correctly")
	}
}

func TestTrainingRestService_RootEndpoint(t *testing.T) {
	mockPort := &mockTrainingInputPort{}
	service := NewTrainingRestService(":8080", mockPort)

	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	mux := http.NewServeMux()
	service.RegisterRoutes(mux)

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	if w.Body.Len() == 0 {
		t.Error("expected non-empty response body")
	}
}

func TestTrainingRestService_HealthEndpoint(t *testing.T) {
	mockPort := &mockTrainingInputPort{}
	service := NewTrainingRestService(":8080", mockPort)

	req := httptest.NewRequest("GET", "/health/", nil)
	w := httptest.NewRecorder()

	mux := http.NewServeMux()
	service.RegisterRoutes(mux)

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	if w.Body.String() != "OK" {
		t.Errorf("expected 'OK', got '%s'", w.Body.String())
	}
}

func TestTrainingRestService_GetAllTrainings(t *testing.T) {
	mockPort := &mockTrainingInputPort{
		trainings: []datamodel.Training{
			{Uid: "uid1", Week: 1, Dayofweek: 1, Activity: "run", Done: false},
		},
	}
	service := NewTrainingRestService(":8080", mockPort)

	req := httptest.NewRequest("GET", "/training", nil)
	w := httptest.NewRecorder()

	mux := http.NewServeMux()
	service.RegisterRoutes(mux)

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var trainings []datamodel.Training
	err := json.NewDecoder(w.Body).Decode(&trainings)
	if err != nil {
		t.Errorf("failed to decode response: %v", err)
	}

	if len(trainings) != 1 {
		t.Errorf("expected 1 training, got %d", len(trainings))
	}
}

func TestTrainingRestService_GetTraining(t *testing.T) {
	mockPort := &mockTrainingInputPort{
		trainings: []datamodel.Training{
			{Uid: "test-uid", Week: 1, Dayofweek: 1, Activity: "run", Done: false},
		},
	}
	service := NewTrainingRestService(":8080", mockPort)

	req := httptest.NewRequest("GET", "/training/test-uid", nil)
	w := httptest.NewRecorder()

	mux := http.NewServeMux()
	service.RegisterRoutes(mux)

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var training datamodel.Training
	err := json.NewDecoder(w.Body).Decode(&training)
	if err != nil {
		t.Errorf("failed to decode response: %v", err)
	}

	if training.Uid != "test-uid" {
		t.Errorf("expected uid 'test-uid', got '%s'", training.Uid)
	}
}

func TestTrainingRestService_PostTraining(t *testing.T) {
	mockPort := &mockTrainingInputPort{}
	service := NewTrainingRestService(":8080", mockPort)

	training := datamodel.Training{
		Uid:       "new-uid",
		Week:      1,
		Dayofweek: 1,
		Activity:  "run",
		Done:      false,
	}

	body, _ := json.Marshal(training)
	req := httptest.NewRequest("POST", "/training/", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	mux := http.NewServeMux()
	service.RegisterRoutes(mux)

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected status 201, got %d", w.Code)
	}

	if mockPort.lastAdded.Uid != "new-uid" {
		t.Error("training was not added through the port")
	}
}

func TestTrainingRestService_PutTraining(t *testing.T) {
	mockPort := &mockTrainingInputPort{
		trainings: []datamodel.Training{
			{Uid: "test-uid", Week: 1, Dayofweek: 1, Activity: "run", Done: false},
		},
	}
	service := NewTrainingRestService(":8080", mockPort)

	updatedTraining := datamodel.Training{
		Uid:       "test-uid",
		Week:      1,
		Dayofweek: 1,
		Activity:  "jog",
		Done:      true,
	}

	body, _ := json.Marshal(updatedTraining)
	req := httptest.NewRequest("PUT", "/training/test-uid", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	mux := http.NewServeMux()
	service.RegisterRoutes(mux)

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	if mockPort.lastUpdated.Uid != "test-uid" || mockPort.lastUpdated.Activity != "jog" {
		t.Error("training was not updated correctly")
	}
}

func TestTrainingRestService_DeleteTraining(t *testing.T) {
	mockPort := &mockTrainingInputPort{
		trainings: []datamodel.Training{
			{Uid: "test-uid", Week: 1, Dayofweek: 1, Activity: "run", Done: false},
		},
	}
	service := NewTrainingRestService(":8080", mockPort)

	req := httptest.NewRequest("DELETE", "/training/test-uid", nil)
	w := httptest.NewRecorder()

	mux := http.NewServeMux()
	service.RegisterRoutes(mux)

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("expected status 204, got %d", w.Code)
	}

	if mockPort.lastDeleted != "test-uid" {
		t.Error("training was not deleted through the port")
	}
}
