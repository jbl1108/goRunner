package usecases

import (
	"github.com/jbl1108/goRunner/usecases/datamodel"
	"github.com/jbl1108/goRunner/usecases/ports/output"
)

type SynchronizeTrainingsUseCase struct {
	trainingSynchronize output.TrainingSynchronize
	trainingDatamodel   datamodel.TrainingList
}

func NewSynchronizeTrainingsUseCase(trainingSynchronize output.TrainingSynchronize, trainingDatamodel datamodel.TrainingList) *SynchronizeTrainingsUseCase {
	return &SynchronizeTrainingsUseCase{
		trainingSynchronize: trainingSynchronize,
		trainingDatamodel:   trainingDatamodel,
	}
}

func (u *SynchronizeTrainingsUseCase) SynchronizeTrainings() error {
	trainings, err := u.trainingSynchronize.GetAllTrainings()
	if err != nil {
		return err
	}
	for _, training := range trainings {
		u.trainingDatamodel.AddTraining(training)
	}
	return nil
}
