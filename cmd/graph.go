package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"codeberg.org/haulproject/haul/api"
	"codeberg.org/haulproject/haul/graph"
	"codeberg.org/haulproject/haul/types"
	"github.com/goccy/go-graphviz"
	"github.com/spf13/cobra"
)

// graphCmd represents the graph command
var graphCmd = &cobra.Command{
	Use:   "graph",
	Short: "Produce a graphviz graph of objects",
	Example: `
Export the haul graph to a file called 'graph.svg':

  $ haul graph --format svg --file graph.svg

Export the haul graph in dot format to stdout

  $ haul graph --format dot

or, with default settings:

  $ haul graph
`,
	Run: func(cmd *cobra.Command, args []string) {
		format, err := cmd.Flags().GetString("format")
		if err != nil {
			log.Fatal("Error:", err)
		}

		if format == "" {
			io.WriteString(os.Stderr, "Error: No graph output format selected.\n\n")
			cmd.Help()
			os.Exit(1)
		}

		filepath, err := cmd.Flags().GetString("file")
		if err != nil {
			log.Fatal("Error:", err)
		}

		var (
			components types.ComponentsWithID
			assemblies types.AssembliesWithID
			kits       types.KitsWithID
		)

		// By default, show all objects in the graph

		components_bytes, err := api.Call(http.MethodGet, "/v1/component")
		if err != nil {
			log.Fatal("Error:", err)
		}

		if err = json.Unmarshal(components_bytes, &components.ComponentsWithID); err != nil {
			log.Fatal("Error:", err)
		}

		assemblies_bytes, err := api.Call(http.MethodGet, "/v1/assembly")
		if err != nil {
			log.Fatal("Error:", err)
		}

		if err = json.Unmarshal(assemblies_bytes, &assemblies.AssembliesWithID); err != nil {
			log.Fatal("Error:", err)
		}

		kits_bytes, err := api.Call(http.MethodGet, "/v1/kit")
		if err != nil {
			log.Fatal("Error:", err)
		}

		if err = json.Unmarshal(kits_bytes, &kits.KitsWithID); err != nil {
			log.Fatal("Error:", err)
		}

		buf, err := graph.GetGraph(graphviz.Format(format), components, assemblies, kits)
		if err != nil {
			log.Fatal("Error:", err)
		}

		if filepath == "" {
			fmt.Println(buf.String())
			os.Exit(0)
		}

		file, err := os.Create(filepath)
		if err != nil {
			log.Fatal("Error:", err)
		}
		defer file.Close()

		file.Write(buf.Bytes())
	},
}

func init() {
	rootCmd.AddCommand(graphCmd)

	graphCmd.Flags().String("format", "dot", "Graph output format")

	graphCmd.Flags().String("file", "", "File to output graph data to. Leave empty for stdout")
}
