# Centralis

Centralis is a Golang library that allows developers to create, manage, and track external resources such as AWS services, third-party database providers, auth providers, and temporary resources. It provides a central place to manage these resources with a simple CRUD interface and handles dependencies between them.

## Features

- **CRUD operations**: Create, read, update, and delete resources.
- **Storage interface**: Store resource metadata in different backends (e.g., AWS S3).
- **Dependency management**: Use a directed acyclic graph (DAG) to resolve and manage dependencies.
- **Provider interface**: Easily extend the library with new providers to handle different types of resources.
- **Synchronous and asynchronous execution**: Execute plans synchronously with rollback on failure or asynchronously.

## Installation

```sh
go get github.com/dblooman/centralis
```

## Usage

Here’s an example of how to use Centralis to manage AWS SNS topics.

```go
package main

import (
	"context"
	"log"

	"github.com/dblooman/centralis/manager"
	"github.com/dblooman/centralis/resource"
	"github.com/dblooman/centralis/centralis"
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
```

## Explanation

    1.	Storage Setup: Initialize the S3 storage backend.
    2.	Resource Manager: Create a resource manager instance.
    3.	SNS Provider: Create an SNS provider and register it with the resource manager.
    4.	Define Resources: Define the resources you want to manage, specifying any dependencies.
    5.	Dependency Resolver: Plan and execute the resource creation using the dependency resolver.
