package model

import ()

type StorageDriver interface {
	AddStats() error

	// Close will clear the state of the storage driver. The elements
	// stored in the underlying storage may or may not be deleted depending
	// on the implementation of the storage driver.
	Close() error
}

func New(name string) (StorageDriver, error) {
	return nil, nil
}
