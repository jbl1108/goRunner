package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"testing"

	"github.com/jbl1108/goRunner/delivery"
	"github.com/jbl1108/goRunner/usecases"
	"github.com/jbl1108/goRunner/usecases/datamodel"
)

// Mock implementations for testing
type mockTrainingPublisher struct {
	data         map[string]datamodel.Training
	lastTraining datamodel.Training
	lastUid      string
}

func newMockTrainingPublisher() *mockTrainingPublisher {
	return &mockTrainingPublisher{
		data: make(map[string]datamodel.Training),
	}
}

func (m *mockTrainingPublisher) TrainingAdded(training datamodel.Training) error {
	m.data[training.Uid] = training
	m.lastTraining = training
	log.Printf("Training added %v", training)
	return nil
}

func (m *mockTrainingPublisher) TrainingUpdated(training datamodel.Training) error {
	m.data[training.Uid] = training
	m.lastTraining = training
	log.Printf("Training updated %v", training)
	return nil
}

func (m *mockTrainingPublisher) TrainingDeleted(uid string) error {
	delete(m.data, uid)
	log.Printf("Training deleted %v", uid)
	m.lastUid = uid
	return nil
}

func (m *mockTrainingPublisher) Connect()    {}
func (m *mockTrainingPublisher) Disconnect() {}

type mockTrainingSynchronize struct {
	trainings []datamodel.Training
}

func (m *mockTrainingSynchronize) GetAllTrainings() ([]datamodel.Training, error) {
	return m.trainings, nil
}

func TestNewApplication(t *testing.T) {
	// This test assumes a config.conf file exists or uses defaults
	mockOutputPublisher := newMockTrainingPublisher()
	trainingDatamodel := datamodel.NewTrainingList()
	trainings := []datamodel.Training{
		{Uid: "uid1", Week: 1, Dayofweek: 1, Activity: "run", Done: false},
		{Uid: "uid2", Week: 2, Dayofweek: 2, Activity: "swim", Done: true},
	}
	for _, t := range trainings {
		trainingDatamodel.AddTraining(t)
	}
	tu := usecases.NewHandleTrainingUseCase(mockOutputPublisher, *trainingDatamodel)
	restService := delivery.NewTrainingRestService("localhost:9900", tu)
	go func() {
		err := restService.Start()
		if err != nil {
			t.Errorf("Failed to start REST service: %v", err)
			return
		}
	}()

	t.Run("Test Add", func(t *testing.T) {
		resp, err := http.Get("http://localhost:9900/training")
		if err != nil {
			t.Fatalf("Failed to get trainings: %v", err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("Expected status 200, got %d", resp.StatusCode)
		}
		defer resp.Body.Close()

		var responseTrainings []datamodel.Training
		err = json.NewDecoder(resp.Body).Decode(&responseTrainings)
		if err != nil {
			bodyBytes, _ := io.ReadAll(resp.Body)
			t.Fatalf("Failed to decode response: %v, body: %s", err, string(bodyBytes))
		}
		if len(responseTrainings) != 2 {
			t.Fatalf("Expected 2 trainings, got %d", len(responseTrainings))
		}
	})
	t.Run("Test Get by UID", func(t *testing.T) {
		resp, err := http.Get("http://localhost:9900/training/uid1")
		if err != nil {
			t.Fatalf("Failed to get training by UID: %v", err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("Expected status 200, got %d", resp.StatusCode)
		}
		defer resp.Body.Close()

		var training datamodel.Training
		err = json.NewDecoder(resp.Body).Decode(&training)
		if err != nil {
			bodyBytes, _ := io.ReadAll(resp.Body)
			t.Fatalf("Failed to decode response: %v, body: %s", err, string(bodyBytes))
		}

		if training.Uid != "uid1" {
			t.Fatalf("Expected UID 'uid1', got '%s'", training.Uid)
		}
	})

	t.Run("Test Add Training", func(t *testing.T) {
		newTraining := datamodel.Training{Uid: "uid3", Week: 3, Dayofweek: 3, Activity: "bike", Done: false}
		jsonData, err := json.Marshal(newTraining)
		if err != nil {
			t.Fatalf("Failed to marshal training: %v", err)
		}

		resp, err := http.Post("http://localhost:9900/training", "application/json", io.NopCloser(bytes.NewReader(jsonData)))
		if err != nil {
			t.Fatalf("Failed to add training: %v", err)
		}
		if resp.StatusCode != http.StatusCreated {
			t.Fatalf("Expected status 201, got %d", resp.StatusCode)
		}
		defer resp.Body.Close()
		var addedTraining datamodel.Training
		err = json.NewDecoder(resp.Body).Decode(&addedTraining)
		if err != nil {
			bodyBytes, _ := io.ReadAll(resp.Body)
			t.Fatalf("Failed to decode response: %v, body: %s", err, string(bodyBytes))
		}

	})
	t.Run("Test Update Training", func(t *testing.T) {
		updatedTraining := datamodel.Training{Week: 1, Dayofweek: 1, Activity: "run", Done: true}
		jsonData, err := json.Marshal(updatedTraining)
		if err != nil {
			t.Fatalf("Failed to marshal training: %v", err)
		}

		req, err := http.NewRequest(http.MethodPut, "http://localhost:9900/training/uid1", bytes.NewReader(jsonData))
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("Failed to update training: %v", err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("Expected status 200, got %d", resp.StatusCode)
		}
		defer resp.Body.Close()

		var updatedResponse datamodel.Training
		err = json.NewDecoder(resp.Body).Decode(&updatedResponse)
		if err != nil {
			bodyBytes, _ := io.ReadAll(resp.Body)
			t.Fatalf("Failed to decode response: %v, body: %s", err, string(bodyBytes))
		}

		if !updatedResponse.Done {
			t.Fatalf("Expected Done to be true, got false")
		}
	})
	t.Run("Test Get All Trainings After Add", func(t *testing.T) {
		resp, err := http.Get("http://localhost:9900/training")
		if err != nil {
			t.Fatalf("Failed to get trainings: %v", err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("Expected status 200, got %d", resp.StatusCode)
		}
		defer resp.Body.Close()

		var responseTrainings []datamodel.Training
		err = json.NewDecoder(resp.Body).Decode(&responseTrainings)
		if err != nil {
			bodyBytes, _ := io.ReadAll(resp.Body)
			t.Fatalf("Failed to decode response: %v, body: %s", err, string(bodyBytes))
		}
		log.Printf("Last response trainings: %v", mockOutputPublisher.lastTraining)
		if len(responseTrainings) != 3 {
			t.Fatalf("Expected 3 trainings, got %d", len(responseTrainings))
		}
	})

}
