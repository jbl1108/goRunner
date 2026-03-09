package output

import "github.com/jbl1108/goRunner/usecases/datamodel"

type TrainingSynchronize interface {
	/* Used for synchronizing trainings from storage */
	GetAllTrainings() ([]datamodel.Training, error)
}
