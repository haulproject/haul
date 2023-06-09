package graph

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"

	"codeberg.org/haulproject/haul/types"
	"github.com/goccy/go-graphviz"
	"github.com/goccy/go-graphviz/cgraph"
)

func GetGraph(format graphviz.Format, components types.ComponentsWithID, assemblies types.AssembliesWithID, kits types.KitsWithID) (*bytes.Buffer, error) {
	g := graphviz.New()

	graph, err := g.Graph()
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := graph.Close(); err != nil {
			log.Fatal(err)
		}
		g.Close()
	}()

	graph.SetRankDir(cgraph.LRRank)

	graph.SetLabel("haul graph")

	// nodes

	for _, component := range components.ComponentsWithID {
		c, err := graph.CreateNode(component.ID.String())
		if err != nil {
			return nil, err
		}

		label_id, err := component.ID.MarshalText()
		if err != nil {
			return nil, err
		}

		status, err := json.Marshal(component.Status)
		if err != nil {
			return nil, err
		}

		tags, err := json.Marshal(component.Tags)
		if err != nil {
			return nil, err
		}

		c.SetLabel(fmt.Sprintf(`%s
		---
		Name: %s
		Status: %s
		Tags: %s`, string(label_id), component.Name, string(status), string(tags)))
	}

	for _, assembly := range assemblies.AssembliesWithID {
		a, err := graph.CreateNode(assembly.ID.String())
		if err != nil {
			return nil, err
		}

		a.SetShape(cgraph.BoxShape)

		label_id, err := assembly.ID.MarshalText()
		if err != nil {
			return nil, err
		}

		status, err := json.Marshal(assembly.Status)
		if err != nil {
			return nil, err
		}

		tags, err := json.Marshal(assembly.Tags)
		if err != nil {
			return nil, err
		}

		a.SetLabel(fmt.Sprintf(`%s
		---
		Name: %s
		Status: %s
		Tags: %s`, string(label_id), assembly.Name, string(status), string(tags)))
	}

	for _, kit := range kits.KitsWithID {
		k, err := graph.CreateNode(kit.ID.String())
		if err != nil {
			return nil, err
		}

		k.SetShape(cgraph.DiamondShape)

		label_id, err := kit.ID.MarshalText()
		if err != nil {
			return nil, err
		}

		status, err := json.Marshal(kit.Status)
		if err != nil {
			return nil, err
		}

		tags, err := json.Marshal(kit.Tags)
		if err != nil {
			return nil, err
		}

		k.SetLabel(fmt.Sprintf(`%s
		---
		Name: %s
		Status: %s
		Tags: %s`, string(label_id), kit.Name, string(status), string(tags)))
	}

	// edges

	for _, component := range components.ComponentsWithID {
		if !component.Target.IsZero() {
			self, err := graph.Node(component.ID.String())
			if err != nil {
				return nil, err
			}

			target, err := graph.Node(component.Target.String())
			if err != nil {
				return nil, err
			}

			e, err := graph.CreateEdge("", self, target)
			if err != nil {
				return nil, err
			}

			_ = e
		}
	}

	for _, assembly := range assemblies.AssembliesWithID {
		if !assembly.Target.IsZero() {
			self, err := graph.Node(assembly.ID.String())
			if err != nil {
				return nil, err
			}

			target, err := graph.Node(assembly.Target.String())
			if err != nil {
				return nil, err
			}

			e, err := graph.CreateEdge("", self, target)
			if err != nil {
				return nil, err
			}

			_ = e
		}
	}

	var buf bytes.Buffer
	if err := g.Render(graph, format, &buf); err != nil {
		return nil, err
	}

	return &buf, nil
}
