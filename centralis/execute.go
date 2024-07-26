package centralis

import (
	"context"
	"sync"

	"github.com/dblooman/centralis/dag"
	"github.com/dblooman/centralis/manager"
	"github.com/dblooman/centralis/resource"
)

type DependencyResolver struct {
	resourceManager *manager.ResourceManager
}

func NewDependencyResolver(rm *manager.ResourceManager) *DependencyResolver {
	return &DependencyResolver{
		resourceManager: rm,
	}
}

func (dr *DependencyResolver) Plan(resources []resource.Resource) (*dag.DAG, error) {
	dag := dag.NewDAG()

	for _, resource := range resources {
		dag.AddNode(resource)
	}

	for _, resource := range resources {
		for _, depID := range resource.Args["dependencies"].([]string) {
			if err := dag.AddDependency(resource.ID, depID); err != nil {
				return nil, err
			}
		}
	}

	_, err := dag.Resolve()
	if err != nil {
		return nil, err
	}

	return dag, nil
}

func (dr *DependencyResolver) Execute(ctx context.Context, plan *dag.DAG, async bool, customFields map[string]interface{}) error {
	if async {
		return dr.executeAsync(ctx, plan, customFields)
	}
	return dr.executeSync(ctx, plan, customFields)
}

func (dr *DependencyResolver) executeSync(ctx context.Context, plan *dag.DAG, customFields map[string]interface{}) error {
	createdResources := make([]string, 0)

	for _, node := range plan.ResolvedNodes {
		_, err := dr.resourceManager.CreateResource(ctx, node.Resource.Type, node.Resource.Args, customFields)
		if err != nil {
			// Rollback created resources
			for _, id := range createdResources {
				resourceType := plan.Nodes[id].Resource.Type
				dr.resourceManager.DeleteResource(ctx, resourceType, id)
			}
			return err
		}
		createdResources = append(createdResources, node.Resource.ID)
	}

	return nil
}

func (dr *DependencyResolver) executeAsync(ctx context.Context, plan *dag.DAG, customFields map[string]interface{}) error {
	var wg sync.WaitGroup
	errorChan := make(chan error, len(plan.ResolvedNodes))

	for _, node := range plan.ResolvedNodes {
		wg.Add(1)
		go func(node *dag.DAGNode) {
			defer wg.Done()
			_, err := dr.resourceManager.CreateResource(ctx, node.Resource.Type, node.Resource.Args, customFields)
			if err != nil {
				errorChan <- err
			}
		}(node)
	}

	wg.Wait()
	close(errorChan)

	if len(errorChan) > 0 {
		return <-errorChan
	}

	return nil
}
