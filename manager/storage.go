package manager

import (
	"context"
	"fmt"

	"github.com/dblooman/centralis/provider"
)

type Storage interface {
	Save(ctx context.Context, resourceType, id string, data map[string]interface{}) error
	Load(ctx context.Context, resourceType, id string) (map[string]interface{}, error)
	Delete(ctx context.Context, resourceType, id string) error
	List(ctx context.Context, resourceType string) ([]string, error)
}

type ResourceManager struct {
	storage   Storage
	providers map[string]provider.Provider
}

func NewResourceManager(storage Storage) *ResourceManager {
	return &ResourceManager{
		storage:   storage,
		providers: make(map[string]provider.Provider),
	}
}

func (rm *ResourceManager) RegisterProvider(ctx context.Context, resourceType string, provider provider.Provider) {
	rm.providers[resourceType] = provider
}

func (rm *ResourceManager) CreateResource(ctx context.Context, resourceType string, args map[string]interface{}) (string, error) {
	provider, exists := rm.providers[resourceType]
	if !exists {
		return "", fmt.Errorf("no provider registered for resource type %s", resourceType)
	}

	id, err := provider.Create(args)
	if err != nil {
		return "", err
	}

	resourceData, err := provider.Read(id)
	if err != nil {
		return "", err
	}

	err = rm.storage.Save(ctx, resourceType, id, resourceData)
	if err != nil {
		return "", err
	}

	return id, nil
}
