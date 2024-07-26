package provider

type CRUD interface {
	Create(args map[string]interface{}) (string, error) // Returns resource ID or ARN
	Read(id string) (map[string]interface{}, error)     // Returns resource details
	Update(id string, args map[string]interface{}) error
	Delete(id string) error
}

type Provider interface {
	CRUD
	GetOutputs(id string) (map[string]interface{}, error) // Outputs that can be used by other providers
}
