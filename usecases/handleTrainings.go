package usecases

import (
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

func (h *HandleTrainingUseCase) AddTraining(training datamodel.Training) error {
	h.trainingModel.AddTraining(training)
	err := h.outputPublisher.TrainingAdded(training)
	if err != nil {
		return err
	}
	return nil
}

func (h *HandleTrainingUseCase) UpdateTraining(training datamodel.Training) error {
	h.trainingModel.UpdateTraining(training)
	err := h.outputPublisher.TrainingUpdated(training)
	if err != nil {
		return err
	}
	return nil
}

func (h *HandleTrainingUseCase) DeleteTraining(uid string) error {
	h.trainingModel.DeleteTraining(uid)
	err := h.outputPublisher.TrainingDeleted(uid)
	if err != nil {
		return err
	}
	return nil
}
