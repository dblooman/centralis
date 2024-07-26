package main

import (
	"context"
	"log"

	"github.com/dblooman/centralis/centralis"
	"github.com/dblooman/centralis/manager"
	"github.com/dblooman/centralis/resource"
	"github.com/dblooman/centralis/sns"
)

func main() {
	ctx := context.Background()
	storage, err := manager.NewS3Storage(ctx, "my-bucket", "us-west-2")
	if err != nil {
		log.Fatalf("Failed to create S3 storage: %v", err)
	}

	rm := manager.NewResourceManager(storage)

	awsSNSProvider, err := sns.NewSNSProvider(ctx, "us-west-2")
	if err != nil {
		log.Fatalf("Failed to create SNS provider: %v", err)
	}

	rm.RegisterProvider(ctx, "aws_sns_topic", awsSNSProvider)

	// Define resources
	resources := []resource.Resource{
		{
			ID:   "resource1",
			Type: "aws_sns_topic",
			Args: map[string]interface{}{
				"name":         "example-topic",
				"dependencies": []string{},
			},
		},
		{
			ID:   "resource2",
			Type: "aws_sns_topic",
			Args: map[string]interface{}{
				"name":         "example-topic-2",
				"dependencies": []string{"resource1"},
			},
		},
	}

	dr := centralis.NewDependencyResolver(rm)
	plan, err := dr.Plan(resources)
	if err != nil {
		log.Fatalf("Failed to plan: %v", err)
	}

	err = dr.Execute(ctx, plan, false)
	if err != nil {
		log.Fatalf("Failed to execute: %v", err)
	}

	log.Println("Resources created successfully")
}
