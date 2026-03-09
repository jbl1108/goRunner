package datamodel

import "fmt"

type Training struct {
	Uid       string `json:"uid"`
	Week      int    `json:"week"`
	Dayofweek int    `json:"dayofweek"`
	Activity  string `json:"activity"`
	Done      bool   `json:"done"`
}

type TrainingList struct {
	trainings map[string]Training
}

func NewTrainingList() *TrainingList {
	return &TrainingList{
		trainings: make(map[string]Training),
	}
}

func (t *TrainingList) AddTraining(training Training) {
	t.trainings[training.Uid] = training
}

func (t *TrainingList) GetAllTrainings() []Training {
	var trainings []Training
	for _, training := range t.trainings {
		trainings = append(trainings, training)
	}
	return trainings
}

func (t *TrainingList) GetTrainingByUid(uid string) (Training, error) {
	if training, exists := t.trainings[uid]; exists {
		return training, nil
	}
	return Training{}, fmt.Errorf("training with uid %s not found", uid)
}

func (t *TrainingList) UpdateTraining(updatedTraining Training) error {
	if _, exists := t.trainings[updatedTraining.Uid]; exists {
		t.trainings[updatedTraining.Uid] = updatedTraining
		return nil
	}
	return fmt.Errorf("training with uid %s not found", updatedTraining.Uid)
}

func (t *TrainingList) DeleteTraining(uid string) error {
	if _, exists := t.trainings[uid]; exists {
		delete(t.trainings, uid)
		return nil
	}
	return fmt.Errorf("training with uid %s not found", uid)
}
