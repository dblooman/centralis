package resource

import (
	"time"

	"github.com/dblooman/centralis/provider"
)

type Resource struct {
	ID       string
	Type     string
	Args     map[string]interface{}
	Outputs  map[string]interface{}
	Provider provider.Provider
}

// ResourceData represents the schema for data stored in the backend
type ResourceData struct {
	ID           string                 `json:"id"`
	Type         string                 `json:"type"`
	Label        string                 `json:"label"`
	ResourceID   string                 `json:"resource_id"`
	Name         string                 `json:"name"`
	CreatedAt    time.Time              `json:"created_at"`
	UpdatedAt    time.Time              `json:"updated_at"`
	CustomFields map[string]interface{} `json:"custom_fields,omitempty"`
}
