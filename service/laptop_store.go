package service

import (
	"context"
	"errors"
	"fmt"
	"gRPC-tutori/pb"
	"github.com/jinzhu/copier"
	"log"
	"sync"
)

var ErrAlreadyExists = errors.New("record already exists")

// LaptopStore is a interface to laptopStore laptop
type LaptopStore interface {
	Save(laptop *pb.Laptop) error
	Find(id string) (*pb.Laptop, error)
	Search(ctx context.Context, filter *pb.Filter, found func(laptop *pb.Laptop) error) error
}

// InMemoryLaptopStore laptopStore laptops in memory
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

	// copy the laptop to laptopStore in the map
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

// Find get the copy of laptop with specified id
func (store *InMemoryLaptopStore) Find(id string) (*pb.Laptop, error) {
	store.mutex.Lock()
	defer store.mutex.Unlock()

	laptop := store.data[id]
	if laptop == nil {
		return nil, nil
	}
	return deepCopy(laptop)
}

// Search find all laptops fit the filter
func (store *InMemoryLaptopStore) Search(ctx context.Context, filter *pb.Filter, found func(laptop *pb.Laptop) error) error {
	store.mutex.RLock()
	defer store.mutex.RUnlock()

	for _, laptop := range store.data {
		if ctx.Err() == context.Canceled || ctx.Err() == context.DeadlineExceeded {
			log.Print("context is canceled")
			return nil
		}

		if isQualified(filter, laptop) {
			other, err := deepCopy(laptop)
			if err != nil {
				return err
			}

			err = found(other)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// isQualified return if the laptop fit the filter
func isQualified(filter *pb.Filter, laptop *pb.Laptop) bool {
	if laptop.GetPriceUsd() > filter.GetMaxPriceUsd() {
		return false
	}
	if laptop.GetCpu().GetNumberCores() < filter.MinCpuCores {
		return false
	}
	if laptop.GetCpu().GetMinGhz() < filter.GetMinCpuGhz() {
		return false
	}
	if toBit(laptop.GetRam()) < toBit(filter.GetMinRam()) {
		return false
	}
	return true
}

// toBit return the memory size of bit
func toBit(memory *pb.Memory) uint64 {
	value := memory.GetValue()

	switch memory.GetUnit() {
	case pb.Memory_BIT:
		return value
	case pb.Memory_BYTE:
		return value << 3
	case pb.Memory_KILOBYTE:
		return value << 13
	case pb.Memory_MEGABYTE:
		return value << 23
	case pb.Memory_GIGABYTE:
		return value << 33
	case pb.Memory_TERABYTE:
		return value << 43
	default:
		return 0
	}
}
