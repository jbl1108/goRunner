package services

import (
	"encoding/json"
	"net/http"

	"github.com/jbl1108/goRunner/usecases/datamodel"
)

type KeyValueStorageService struct {
	requestUrl string
}

func NewKeyValueRepository(requestUrl string) *KeyValueStorageService {
	return &KeyValueStorageService{
		requestUrl: requestUrl,
	}
}

func (k *KeyValueStorageService) GetAllTrainings() ([]datamodel.Training, error) {
	response, err := http.Get(k.requestUrl)
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
