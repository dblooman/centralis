package manager

import (
	"context"
	"fmt"
)

func (rm *ResourceManager) DeleteResource(ctx context.Context, resourceType string, id string) error {
	provider, exists := rm.providers[resourceType]
	if !exists {
		return fmt.Errorf("no provider registered for resource type %s", resourceType)
	}
	err := provider.Delete(id)
	if err != nil {
		return err
	}
	err = rm.storage.Delete(ctx, resourceType, id)
	if err != nil {
		return err
	}
	return nil
}
