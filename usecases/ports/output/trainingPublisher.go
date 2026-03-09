package output

import "github.com/jbl1108/goRunner/usecases/datamodel"

type TrainingPublisher interface {
	/* Called when a training is added, updated or deleted */

	TrainingUpdated(training datamodel.Training) error

	TrainingAdded(training datamodel.Training) error

	TrainingDeleted(uid string) error

	Connect()

	Disconnect()
}
