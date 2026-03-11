package usecases

import (
	"testing"

	"github.com/jbl1108/goRunner/usecases/datamodel"
)

// Mock TrainingPublisher for testing
type mockPublisher struct {
	added     []datamodel.Training
	updated   []datamodel.Training
	deleted   []string
	addErr    error
	updateErr error
	deleteErr error
}

func (m *mockPublisher) TrainingAdded(training datamodel.Training) error {
	if m.addErr != nil {
		return m.addErr
	}
	m.added = append(m.added, training)
	return nil
}

func (m *mockPublisher) TrainingUpdated(training datamodel.Training) error {
	if m.updateErr != nil {
		return m.updateErr
	}
	m.updated = append(m.updated, training)
	return nil
}

func (m *mockPublisher) TrainingDeleted(uid string) error {
	if m.deleteErr != nil {
		return m.deleteErr
	}
	m.deleted = append(m.deleted, uid)
	return nil
}

func (m *mockPublisher) Connect()    {}
func (m *mockPublisher) Disconnect() {}

func TestNewHandleTrainingUseCase(t *testing.T) {
	mockPub := &mockPublisher{}
	trainingList := datamodel.NewTrainingList()

	usecase := NewHandleTrainingUseCase(mockPub, *trainingList)

	if usecase.outputPublisher != mockPub {
		t.Error("outputPublisher not set correctly")
	}
}

func TestHandleTrainingUseCase_GetAllTrainings(t *testing.T) {
	mockPub := &mockPublisher{}
	trainingList := datamodel.NewTrainingList()

	training1 := datamodel.Training{Uid: "uid1", Week: 1, Dayofweek: 1, Activity: "run", Done: false}
	training2 := datamodel.Training{Uid: "uid2", Week: 2, Dayofweek: 2, Activity: "swim", Done: true}

	trainingList.AddTraining(training1)
	trainingList.AddTraining(training2)

	usecase := NewHandleTrainingUseCase(mockPub, *trainingList)

	trainings, err := usecase.GetAllTrainings()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if len(trainings) != 2 {
		t.Errorf("expected 2 trainings, got %d", len(trainings))
	}
}

func TestHandleTrainingUseCase_GetTraining(t *testing.T) {
	mockPub := &mockPublisher{}
	trainingList := datamodel.NewTrainingList()

	training := datamodel.Training{Uid: "test-uid", Week: 1, Dayofweek: 1, Activity: "run", Done: false}
	trainingList.AddTraining(training)

	usecase := NewHandleTrainingUseCase(mockPub, *trainingList)

	// Test successful retrieval
	retrieved, err := usecase.GetTraining("test-uid")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if retrieved.Uid != "test-uid" {
		t.Errorf("wrong training retrieved: got %s, want test-uid", retrieved.Uid)
	}

	// Test retrieval of non-existent training
	_, err = usecase.GetTraining("non-existent")
	if err == nil {
		t.Error("expected error for non-existent training")
	}
}

func TestHandleTrainingUseCase_AddTraining(t *testing.T) {
	mockPub := &mockPublisher{}
	trainingList := datamodel.NewTrainingList()

	usecase := NewHandleTrainingUseCase(mockPub, *trainingList)

	training := datamodel.Training{Week: 1, Dayofweek: 1, Activity: "run", Done: false}

	addedTraining, err := usecase.AddTraining(training)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Check that training was added to the model
	retrieved, err := trainingList.GetTrainingByUid(addedTraining.Uid)
	if err != nil {
		t.Errorf("training was not added to model: %v", err)
	}
	if retrieved.Uid != addedTraining.Uid {
		t.Error("wrong training added to model")
	}

	// Check that publisher was called
	if len(mockPub.added) != 1 {
		t.Errorf("expected 1 added call, got %d", len(mockPub.added))
	}
	if mockPub.added[0].Uid != addedTraining.Uid {
		t.Error("wrong training passed to publisher")
	}
}

func TestHandleTrainingUseCase_UpdateTraining(t *testing.T) {
	mockPub := &mockPublisher{}
	trainingList := datamodel.NewTrainingList()

	original := datamodel.Training{Uid: "test-uid", Week: 1, Dayofweek: 1, Activity: "run", Done: false}
	trainingList.AddTraining(original)

	usecase := NewHandleTrainingUseCase(mockPub, *trainingList)

	updated := datamodel.Training{Uid: "test-uid", Week: 1, Dayofweek: 1, Activity: "jog", Done: true}

	_, err := usecase.UpdateTraining(updated)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Check that training was updated in the model
	retrieved, err := trainingList.GetTrainingByUid("test-uid")
	if err != nil {
		t.Errorf("training was not found after update: %v", err)
	}
	if retrieved.Activity != "jog" || !retrieved.Done {
		t.Error("training was not updated correctly")
	}

	// Check that publisher was called
	if len(mockPub.updated) != 1 {
		t.Errorf("expected 1 updated call, got %d", len(mockPub.updated))
	}
}

func TestHandleTrainingUseCase_DeleteTraining(t *testing.T) {
	mockPub := &mockPublisher{}
	trainingList := datamodel.NewTrainingList()

	training := datamodel.Training{Uid: "test-uid", Week: 1, Dayofweek: 1, Activity: "run", Done: false}
	trainingList.AddTraining(training)

	usecase := NewHandleTrainingUseCase(mockPub, *trainingList)

	err := usecase.DeleteTraining("test-uid")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Check that training was deleted from the model
	_, err = trainingList.GetTrainingByUid("test-uid")
	if err == nil {
		t.Error("training should have been deleted from model")
	}

	// Check that publisher was called
	if len(mockPub.deleted) != 1 {
		t.Errorf("expected 1 deleted call, got %d", len(mockPub.deleted))
	}
	if mockPub.deleted[0] != "test-uid" {
		t.Error("wrong uid passed to publisher")
	}
}
