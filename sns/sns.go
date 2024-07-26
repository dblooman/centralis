package sns

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sns"
)

// SNSProvider implements the Provider interface for AWS SNS
type SNSProvider struct {
	svc *sns.Client
}

// NewSNSProvider creates a new SNSProvider instance
func NewSNSProvider(ctx context.Context, region string) (*SNSProvider, error) {
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		return nil, err
	}

	svc := sns.NewFromConfig(cfg)

	return &SNSProvider{
		svc: svc,
	}, nil
}

// Create creates a new SNS topic
func (p *SNSProvider) Create(args map[string]interface{}) (string, error) {
	input := &sns.CreateTopicInput{
		Name: aws.String(args["name"].(string)),
	}

	result, err := p.svc.CreateTopic(context.TODO(), input)
	if err != nil {
		return "", err
	}

	return *result.TopicArn, nil
}

// Read retrieves the SNS topic attributes
func (p *SNSProvider) Read(id string) (map[string]interface{}, error) {
	input := &sns.GetTopicAttributesInput{
		TopicArn: aws.String(id),
	}

	result, err := p.svc.GetTopicAttributes(context.TODO(), input)
	if err != nil {
		return nil, err
	}

	attributes := make(map[string]interface{})
	for k, v := range result.Attributes {
		attributes[k] = v
	}

	return attributes, nil
}

// Update updates the SNS topic attributes
func (p *SNSProvider) Update(id string, args map[string]interface{}) error {
	input := &sns.SetTopicAttributesInput{
		TopicArn:       aws.String(id),
		AttributeName:  aws.String(args["attribute_name"].(string)),
		AttributeValue: aws.String(args["attribute_value"].(string)),
	}

	_, err := p.svc.SetTopicAttributes(context.TODO(), input)
	return err
}

// Delete deletes the SNS topic
func (p *SNSProvider) Delete(id string) error {
	input := &sns.DeleteTopicInput{
		TopicArn: aws.String(id),
	}

	_, err := p.svc.DeleteTopic(context.TODO(), input)
	return err
}

// GetOutputs retrieves the SNS topic ARN (used as an output for other resources)
func (p *SNSProvider) GetOutputs(id string) (map[string]interface{}, error) {
	attributes, err := p.Read(id)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"TopicArn":   id,
		"Attributes": attributes,
	}, nil
}
