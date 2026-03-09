package services

import (
	"encoding/json"
	"net/http"

	"github.com/jbl1108/goRunner/usecases/datamodel"
)

type KeyValueRepository struct {
	requestUrl string
}

func NewKeyValueRepository(requestUrl string) *KeyValueRepository {
	return &KeyValueRepository{
		requestUrl: requestUrl,
	}
}

func (r *KeyValueRepository) GetAllTrainings() ([]datamodel.Training, error) {
	response, err := http.Get(r.requestUrl)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	var trainings []datamodel.Training
	err = json.NewDecoder(response.Body).Decode(&trainings)
	if err != nil {
		return nil, err
	}

	return trainings, nil
}
