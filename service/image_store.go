package service

import (
	"bytes"
	"fmt"
	"github.com/google/uuid"
	"os"
	"sync"
)

// ImageStore is an interface that could laptopStore image
type ImageStore interface {
	Save(laptopID string, imageType string, imageData bytes.Buffer) (string, error)
}

// DiskImageStore save the image to the disk
type DiskImageStore struct {
	mutex       sync.RWMutex
	imageFolder string
	images      map[string]*ImageInfo
}

// ImageInfo contains the basic info of image
type ImageInfo struct {
	LaptopID string
	Type     string
	Path     string
}

// NewDiskImageStore creates a new DiskImageStore
func NewDiskImageStore(imageFolder string) *DiskImageStore {
	return &DiskImageStore{
		imageFolder: imageFolder,
		images:      make(map[string]*ImageInfo),
	}
}

// Save save the image data to disk
func (store *DiskImageStore) Save(laptopID string, imageType string, imageData bytes.Buffer) (string, error) {
	imageID, err := uuid.NewRandom()
	if err != nil {
		return "", fmt.Errorf("cannot generate image id:%w", err)
	}

	imagePath := fmt.Sprintf("%s/%s%s", store.imageFolder, imageID, imageType)

	file, err := os.Create(imagePath)
	if err != nil {
		return "", fmt.Errorf("cannot create image file:%w", err)
	}

	_, err = imageData.WriteTo(file)
	if err != nil {
		return "", fmt.Errorf("cannot write image to file:%w", err)
	}

	store.mutex.Lock()
	defer store.mutex.Unlock()

	store.images[imageID.String()] = &ImageInfo{
		laptopID,
		imageType,
		imagePath,
	}

	return imageID.String(), nil
}
