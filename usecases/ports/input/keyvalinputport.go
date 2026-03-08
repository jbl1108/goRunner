package input

import "github.com/jbl1108/github.com/jbl1108/goSecret/usecases/datamodel"

type KeyValInputPort interface {
	/*
		Retrieves the value associated with the provided key.
		Returns the value as a byte slice and an error if the operation fails.
	*/
	GetKey(key string) ([]byte, error)
	/*	Sets the value for the provided key.
		Takes a key as a string and a value as a byte slice.
		Returns an error if the operation fails.
	*/
	SetKey(message datamodel.Message) error
}
