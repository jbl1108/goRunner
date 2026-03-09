package input

import "github.com/jbl1108/goRunner/usecases/datamodel"

type TrainingInputPort interface {

	/* AddTraining adds a training to the system */
	AddTraining(training datamodel.Training) error

	/* GetTraining retrieves a training by its UID */
	GetTraining(uid string) (datamodel.Training, error)

	/* GetAllTrainings retrieves all trainings in the system */
	GetAllTrainings() ([]datamodel.Training, error)

	/* UpdateTraining updates an existing training */
	UpdateTraining(training datamodel.Training) error

	/* DeleteTraining removes a training by its UID */
	DeleteTraining(uid string) error
}
