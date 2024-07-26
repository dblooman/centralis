package resource

import "github.com/dblooman/centralis/provider"

type Resource struct {
	ID       string
	Type     string
	Args     map[string]interface{}
	Outputs  map[string]interface{}
	Provider provider.Provider
}
