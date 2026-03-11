package usecases

import (
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/jbl1108/goRunner/usecases/datamodel"
	"github.com/jbl1108/goRunner/usecases/ports/output"
)

type HandleTrainingUseCase struct {
	outputPublisher output.TrainingPublisher
	trainingModel   datamodel.TrainingList
}

func NewHandleTrainingUseCase(outputPublisher output.TrainingPublisher, traningDatamodel datamodel.TrainingList) *HandleTrainingUseCase {
	return &HandleTrainingUseCase{
		outputPublisher: outputPublisher,
		trainingModel:   traningDatamodel,
	}
}

func (h *HandleTrainingUseCase) GetAllTrainings() ([]datamodel.Training, error) {
	return h.trainingModel.GetAllTrainings(), nil
}

func (h *HandleTrainingUseCase) GetTraining(uid string) (datamodel.Training, error) {
	return h.trainingModel.GetTrainingByUid(uid)
}

func (h *HandleTrainingUseCase) AddTraining(training datamodel.Training) (datamodel.Training, error) {
	if h.trainingModel.Exists(training.Uid) {
		return h.UpdateTraining(training)
	}
	log.Printf("Adding training: %v", training)

	training.Uid = uuid.New().String()
	h.trainingModel.AddTraining(training)
	err := h.outputPublisher.TrainingAdded(training)

	if err != nil {
		return datamodel.Training{}, err
	}
	return training, nil
}

func (h *HandleTrainingUseCase) UpdateTraining(training datamodel.Training) (datamodel.Training, error) {
	if !h.trainingModel.Exists(training.Uid) {
		return datamodel.Training{}, fmt.Errorf("training with uid %s not found", training.Uid)
	}
	h.trainingModel.UpdateTraining(training)
	err := h.outputPublisher.TrainingUpdated(training)
	if err != nil {
		return datamodel.Training{}, err
	}
	return training, nil
}

func (h *HandleTrainingUseCase) DeleteTraining(uid string) error {
	h.trainingModel.DeleteTraining(uid)
	err := h.outputPublisher.TrainingDeleted(uid)
	if err != nil {
		return err
	}
	return nil
}
