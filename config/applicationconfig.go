package config

import (
	"github.com/jbl1108/goRunner/delivery"
	"github.com/jbl1108/goRunner/services"
	"github.com/jbl1108/goRunner/usecases"
	"github.com/jbl1108/goRunner/usecases/datamodel"
	"github.com/jbl1108/goRunner/usecases/ports/output"
)

type Application struct {
	OutputPublisher             output.TrainingPublisher
	TrainingSynchronize         output.TrainingSynchronize
	trainingDatamodel           datamodel.TrainingList
	RestService                 *delivery.TrainingRestService
	handleTrainingUseCase       *usecases.HandleTrainingUseCase
	SynchronizeTrainingsUseCase *usecases.SynchronizeTrainingsUseCase
}

func NewApplication() Application {
	c := NewConfig()
	outputPublisher := delivery.NewMQTTClient(c.MQTTAddress(), c.MQTTUsername(), c.MQTTPassword(), "trainings")
	trainingDatamodel := datamodel.NewTrainingList()
	trainingSynchronize := services.NewKeyValueRepository(c.KeyValueDBURL())
	synchronizeTrainingsUseCase := usecases.NewSynchronizeTrainingsUseCase(trainingSynchronize, *trainingDatamodel)

	handleTrainingUseCase := usecases.NewHandleTrainingUseCase(outputPublisher, *trainingDatamodel)

	restService := delivery.NewTrainingRestService(c.RestAddress(), handleTrainingUseCase) // Will set usecase later to avoid circular dependency

	return Application{
		OutputPublisher:             outputPublisher,
		TrainingSynchronize:         trainingSynchronize,
		trainingDatamodel:           *trainingDatamodel,
		RestService:                 restService,
		handleTrainingUseCase:       handleTrainingUseCase,
		SynchronizeTrainingsUseCase: synchronizeTrainingsUseCase,
	}

}
