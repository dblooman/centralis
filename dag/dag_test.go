package dag

import (
	"testing"

	"github.com/dblooman/centralis/resource"
)

func TestNewDAG(t *testing.T) {
	dag := NewDAG()

	if dag == nil {
		t.Error("NewDAG should return a non-nil DAG")
	}

	if dag == nil || dag.Nodes == nil {
		t.Error("NewDAG should initialize the Nodes map")
	}

	if dag.ResolvedNodes == nil {
		t.Error("NewDAG should initialize the ResolvedNodes slice")
	}
}

func TestAddNode(t *testing.T) {
	dag := NewDAG()
	resource := resource.Resource{ID: "resource1"}

	dag.AddNode(resource)

	if len(dag.Nodes) != 1 {
		t.Error("AddNode should add the node to the Nodes map")
	}

	if len(dag.ResolvedNodes) != 0 {
		t.Error("AddNode should not add the node to the ResolvedNodes slice")
	}
}

func TestAddDependency(t *testing.T) {
	dag := NewDAG()
	resource1 := resource.Resource{ID: "resource1"}
	resource2 := resource.Resource{ID: "resource2"}

	dag.AddNode(resource1)
	dag.AddNode(resource2)

	err := dag.AddDependency("resource1", "resource2")
	if err != nil {
		t.Errorf("AddDependency returned an error: %v", err)
	}

	if len(dag.Nodes["resource1"].Dependencies) != 1 {
		t.Error("AddDependency should add the dependency to the node's Dependencies slice")
	}

	if len(dag.ResolvedNodes) != 0 {
		t.Error("AddDependency should not add the node to the ResolvedNodes slice")
	}
}

func TestResolve(t *testing.T) {
	dag := NewDAG()
	resource1 := resource.Resource{ID: "resource1"}
	resource2 := resource.Resource{ID: "resource2"}

	dag.AddNode(resource1)
	dag.AddNode(resource2)
	dag.AddDependency("resource1", "resource2")

	resolvedNodes, err := dag.Resolve()
	if err != nil {
		t.Errorf("Resolve returned an error: %v", err)
	}

	if len(resolvedNodes) != 2 {
		t.Error("Resolve should return all nodes in the correct order")
	}

	if len(dag.ResolvedNodes) != 2 {
		t.Error("Resolve should populate the ResolvedNodes slice")
	}
}
