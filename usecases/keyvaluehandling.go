package usecases

import (
	"github.com/jbl1108/github.com/jbl1108/goSecret/usecases/datamodel"
	"github.com/jbl1108/github.com/jbl1108/goSecret/usecases/ports/output"
)

type KeyValueHandling struct {
	storage output.KeyValueStorage
}

func NewKeyValueHandling(storage output.KeyValueStorage) *KeyValueHandling {
	return &KeyValueHandling{storage: storage}
}

func (uc *KeyValueHandling) SetKey(message datamodel.Message) error {
	err := uc.storage.Open()
	if err != nil {
		return err
	}
	defer uc.storage.Close()
	key := message.Topic + ":" + message.Data.Key
	return uc.storage.Set(key, []byte(message.Data.Value))
}
func (uc *KeyValueHandling) GetKey(key string) ([]byte, error) {
	err := uc.storage.Open()
	if err != nil {
		return nil, err
	}
	defer uc.storage.Close()

	return uc.storage.Get(key)
}
