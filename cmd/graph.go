package cmd

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
	"github.com/spf13/cobra"
)

// graphCmd represents the graph command
var graphCmd = &cobra.Command{
	Use:   "graph",
	Short: "Produce a graphviz graph of objects",
	Run: func(cmd *cobra.Command, args []string) {

		// No outputstyle needed for something that outputs graph data
		/*
			output, err := rootCmd.PersistentFlags().GetString("output")
			if err != nil {
				log.Fatal(err)
			}

			client := cli.New()
			client.OutputStyle = output
		*/

		var (
			components types.ComponentsWithID
			assemblies types.AssembliesWithID
			kits       types.KitsWithID
		)

		// By default, show all objects in the graph

		components_bytes, err := api.Call(http.MethodGet, "/v1/component")
		if err != nil {
			log.Fatal(err)
		}

		if err = json.Unmarshal(components_bytes, &components.ComponentsWithID); err != nil {
			log.Fatal(err)
		}

		assemblies_bytes, err := api.Call(http.MethodGet, "/v1/assembly")
		if err != nil {
			log.Fatal(err)
		}

		if err = json.Unmarshal(assemblies_bytes, &assemblies.AssembliesWithID); err != nil {
			log.Fatal(err)
		}

		kits_bytes, err := api.Call(http.MethodGet, "/v1/kit")
		if err != nil {
			log.Fatal(err)
		}

		if err = json.Unmarshal(kits_bytes, &kits.KitsWithID); err != nil {
			log.Fatal(err)
		}

		//TODO api calls
		g := graphviz.New()

		graph, err := g.Graph()
		if err != nil {
			log.Fatal(err)
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
				log.Fatal(err)
			}

			label_id, err := component.ID.MarshalText()
			if err != nil {
				log.Fatal(err)
			}

			c.SetLabel(fmt.Sprint(component.Name, "\n---\n", string(label_id)))
		}

		for _, assembly := range assemblies.AssembliesWithID {
			a, err := graph.CreateNode(assembly.ID.String())
			if err != nil {
				log.Fatal(err)
			}

			a.SetShape(cgraph.BoxShape)

			label_id, err := assembly.ID.MarshalText()
			if err != nil {
				log.Fatal(err)
			}

			a.SetLabel(fmt.Sprint(assembly.Name, "\n---\n", string(label_id)))
		}

		for _, kit := range kits.KitsWithID {
			k, err := graph.CreateNode(kit.ID.String())
			if err != nil {
				log.Fatal(err)
			}

			k.SetShape(cgraph.DiamondShape)

			label_id, err := kit.ID.MarshalText()
			if err != nil {
				log.Fatal(err)
			}

			k.SetLabel(fmt.Sprint(kit.Name, "\n---\n", string(label_id)))
		}

		// edges

		for _, component := range components.ComponentsWithID {
			if !component.Target.IsZero() {
				self, err := graph.Node(component.ID.String())
				if err != nil {
					log.Fatal(err)
				}

				target, err := graph.Node(component.Target.String())
				if err != nil {
					log.Fatal(err)
				}

				e, err := graph.CreateEdge("", self, target)
				if err != nil {
					log.Fatal(err)
				}

				_ = e
			}
		}

		for _, assembly := range assemblies.AssembliesWithID {
			if !assembly.Target.IsZero() {
				self, err := graph.Node(assembly.ID.String())
				if err != nil {
					log.Fatal(err)
				}

				target, err := graph.Node(assembly.Target.String())
				if err != nil {
					log.Fatal(err)
				}

				e, err := graph.CreateEdge("", self, target)
				if err != nil {
					log.Fatal(err)
				}

				_ = e
			}
		}

		var buf bytes.Buffer
		if err := g.Render(graph, "dot", &buf); err != nil {
			log.Fatal(err)
		}

		fmt.Println(buf.String())
	},
}

func init() {
	rootCmd.AddCommand(graphCmd)
}
