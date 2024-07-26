package dag

import (
	"fmt"

	"github.com/dblooman/centralis/resource"
)

type DAGNode struct {
	Resource     resource.Resource
	Dependencies []*DAGNode
}

type DAG struct {
	Nodes         map[string]*DAGNode
	ResolvedNodes []*DAGNode
}

func NewDAG() *DAG {
	return &DAG{
		Nodes: make(map[string]*DAGNode),
	}
}

func (dag *DAG) AddNode(resource resource.Resource) {
	node := &DAGNode{
		Resource: resource,
	}
	dag.Nodes[resource.ID] = node
}

func (dag *DAG) AddDependency(resourceID string, dependencyID string) error {
	resourceNode, resourceExists := dag.Nodes[resourceID]
	dependencyNode, dependencyExists := dag.Nodes[dependencyID]

	if !resourceExists || !dependencyExists {
		return fmt.Errorf("resource or dependency not found in DAG")
	}

	resourceNode.Dependencies = append(resourceNode.Dependencies, dependencyNode)
	return nil
}

func (dag *DAG) Resolve() ([]*DAGNode, error) {
	visited := make(map[string]bool)
	stack := make([]*DAGNode, 0)

	var visit func(node *DAGNode) error
	visit = func(node *DAGNode) error {
		if visited[node.Resource.ID] {
			return nil
		}
		visited[node.Resource.ID] = true

		for _, dep := range node.Dependencies {
			if err := visit(dep); err != nil {
				return err
			}
		}

		stack = append(stack, node)
		return nil
	}

	for _, node := range dag.Nodes {
		if err := visit(node); err != nil {
			return nil, err
		}
	}

	dag.ResolvedNodes = stack
	return stack, nil
}
