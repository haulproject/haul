package graph

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"codeberg.org/haulproject/haul/api"
	"codeberg.org/haulproject/haul/types"
	"github.com/goccy/go-graphviz"
	"github.com/goccy/go-graphviz/cgraph"
)

func GetGraph() (*bytes.Buffer, error) {
	var (
		components types.ComponentsWithID
		assemblies types.AssembliesWithID
		kits       types.KitsWithID
	)

	// By default, show all objects in the graph

	components_bytes, err := api.Call(http.MethodGet, "/v1/component")
	if err != nil {
		return nil, err
	}

	if err = json.Unmarshal(components_bytes, &components.ComponentsWithID); err != nil {
		return nil, err
	}

	assemblies_bytes, err := api.Call(http.MethodGet, "/v1/assembly")
	if err != nil {
		return nil, err
	}

	if err = json.Unmarshal(assemblies_bytes, &assemblies.AssembliesWithID); err != nil {
		return nil, err
	}

	kits_bytes, err := api.Call(http.MethodGet, "/v1/kit")
	if err != nil {
		return nil, err
	}

	if err = json.Unmarshal(kits_bytes, &kits.KitsWithID); err != nil {
		return nil, err
	}

	//TODO api calls
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

		c.SetLabel(fmt.Sprintf(`%s
		---
		%s`, string(label_id), component.Name))
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

		a.SetLabel(fmt.Sprintf(`%s
		---
		%s`, string(label_id), assembly.Name))
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

		k.SetLabel(fmt.Sprintf(`%s
		---
		%s`, string(label_id), kit.Name))
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
	if err := g.Render(graph, "dot", &buf); err != nil {
		return nil, err
	}
	return &buf, nil
}
