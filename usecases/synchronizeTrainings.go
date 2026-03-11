package usecases

import (
	"github.com/jbl1108/goRunner/usecases/datamodel"
	"github.com/jbl1108/goRunner/usecases/ports/output"
)

type SynchronizeTrainingsUseCase struct {
	trainingSynchronize output.TrainingSynchronize
}

func NewSynchronizeTrainingsUseCase(trainingSynchronize output.TrainingSynchronize, trainingDatamodel datamodel.TrainingList) *SynchronizeTrainingsUseCase {
	return &SynchronizeTrainingsUseCase{
		trainingSynchronize: trainingSynchronize,
	}
}

func (u *SynchronizeTrainingsUseCase) SynchronizeTrainings() (datamodel.TrainingList, error) {
	trainings, err := u.trainingSynchronize.GetAllTrainings()
	if err != nil {
		return datamodel.TrainingList{}, err
	}
	trainingDatamodel := datamodel.NewTrainingList()
	for _, training := range trainings {
		trainingDatamodel.AddTraining(training)
	}
	return *trainingDatamodel, nil
}
