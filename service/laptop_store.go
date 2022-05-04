package service

import (
	"errors"
	"fmt"
	"gRPC-tutori/pb"
	"github.com/jinzhu/copier"
	"sync"
)

var ErrAlreadyExists = errors.New("record already exists")

// LaptopStore is a interface to store laptop
type LaptopStore interface {
	Save(laptop *pb.Laptop) error
}

// InMemoryLaptopStore store laptops in memory
type InMemoryLaptopStore struct {
	mutex sync.RWMutex
	data  map[string]*pb.Laptop
}

// NewInMemoryLaptopStore create a new InMemoryLaptopStore example
func NewInMemoryLaptopStore() *InMemoryLaptopStore {
	return &InMemoryLaptopStore{
		data: make(map[string]*pb.Laptop),
	}
}

// Save save laptop in memory
func (store *InMemoryLaptopStore) Save(laptop *pb.Laptop) error {
	store.mutex.Lock()
	defer store.mutex.Unlock()

	// if the laptop is already in map,return error
	if store.data[laptop.Id] != nil {
		return ErrAlreadyExists
	}

	// copy the laptop to store in the map
	other, err := deepCopy(laptop)
	if err != nil {
		return nil
	}

	store.data[other.Id] = other
	return nil
}

// deepCopy return a copy of laptop
func deepCopy(laptop *pb.Laptop) (*pb.Laptop, error) {
	other := &pb.Laptop{}

	err := copier.Copy(other, laptop)
	if err != nil {
		return nil, fmt.Errorf("cannot copy laptop data:%w", err)
	}

	return other, nil
}