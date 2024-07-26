package manager

import (
	"context"
	"fmt"
	"time"

	"github.com/dblooman/centralis/provider"
	"github.com/dblooman/centralis/resource"
	"github.com/google/uuid"
)

type Storage interface {
	Save(ctx context.Context, data resource.ResourceData) error
	Load(ctx context.Context, resourceType, id string) (resource.ResourceData, error)
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

func (rm *ResourceManager) CreateResource(ctx context.Context, resourceType string, args map[string]interface{}, customFields map[string]interface{}) (string, error) {
	provider, exists := rm.providers[resourceType]
	if !exists {
		return "", fmt.Errorf("no provider registered for resource type %s", resourceType)
	}

	// Generate a unique ID for the resource
	id := uuid.New().String()

	// Create the resource using the provider
	resourceID, err := provider.Create(args)
	if err != nil {
		return "", err
	}

	// Construct ResourceData with appropriate checks
	resourceData := resource.ResourceData{
		ID:           id,
		Type:         resourceType,
		Label:        args["label"].(string),
		ResourceID:   resourceID,
		Name:         args["name"].(string),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		CustomFields: customFields,
	}

	// Save the resource data to storage
	err = rm.storage.Save(ctx, resourceData)
	if err != nil {
		return "", err
	}

	return resourceID, nil
}

func (rm *ResourceManager) ReadResource(ctx context.Context, resourceType, id string) (resource.ResourceData, error) {
	resourceData, err := rm.storage.Load(ctx, resourceType, id)
	if err != nil {
		return resource.ResourceData{}, err
	}
	return resourceData, nil
}

func (rm *ResourceManager) UpdateResource(ctx context.Context, resourceType, id string, args map[string]interface{}, customFields map[string]interface{}) error {
	provider, exists := rm.providers[resourceType]
	if !exists {
		return fmt.Errorf("no provider registered for resource type %s", resourceType)
	}

	// Retrieve the current resource data
	resourceData, err := rm.storage.Load(ctx, resourceType, id)
	if err != nil {
		return err
	}

	// Update the resource using the provider
	err = provider.Update(resourceData.ResourceID, args)
	if err != nil {
		return err
	}

	// Update the resource data
	resourceData.Label = args["label"].(string)
	resourceData.Name = args["name"].(string)
	resourceData.UpdatedAt = time.Now()
	resourceData.CustomFields = customFields

	// Save the updated resource data to storage
	err = rm.storage.Save(ctx, resourceData)
	if err != nil {
		return err
	}

	return nil
}

func (rm *ResourceManager) DeleteResource(ctx context.Context, resourceType, id string) error {
	provider, exists := rm.providers[resourceType]
	if !exists {
		return fmt.Errorf("no provider registered for resource type %s", resourceType)
	}

	// Retrieve the current resource data
	resourceData, err := rm.storage.Load(ctx, resourceType, id)
	if err != nil {
		return err
	}

	// Delete the resource using the provider
	err = provider.Delete(resourceData.ResourceID)
	if err != nil {
		return err
	}

	// Delete the resource data from storage
	err = rm.storage.Delete(ctx, resourceType, id)
	if err != nil {
		return err
	}

	return nil
}
