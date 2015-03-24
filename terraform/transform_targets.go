package terraform

import "github.com/hashicorp/terraform/dag"

// TargetsTransformer is a GraphTransformer that, when the user specifies a
// list of resources to target, limits the graph to only those resources and
// their dependencies.
type TargetsTransformer struct {
	Targets []string
}

func (t *TargetsTransformer) Transform(g *Graph) error {
	if len(t.Targets) > 0 {
		targetedNodes, err := t.selectTargetedNodes(g)
		if err != nil {
			return err
		}

		for _, v := range g.Vertices() {
			if !targetedNodes.Include(v) {
				g.Remove(v)
			}
		}
	}
	return nil
}

func (t *TargetsTransformer) selectTargetedNodes(g *Graph) (*dag.Set, error) {
	targetedNodes := new(dag.Set)
	for _, v := range g.Vertices() {
		// We only care about Resources and their deps
		r, ok := v.(*GraphNodeConfigResource)
		if !ok {
			continue
		}

		if t.resourceIsTarget(r) {
			targetedNodes.Add(r)
			deps, err := g.DownVertices(r)
			if err != nil {
				return nil, err
			}
			for _, d := range deps.List() {
				targetedNodes.Add(d)
			}
		}
	}
	return targetedNodes, nil
}

func (t *TargetsTransformer) resourceIsTarget(r *GraphNodeConfigResource) bool {
	for _, target := range t.Targets {
		if target == r.Name() {
			return true
		}
	}
	return false
}
