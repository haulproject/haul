package cmd

import (
	"fmt"
	"log"

	"codeberg.org/haulproject/haul/graph"
	"github.com/spf13/cobra"
)

// graphCmd represents the graph command
var graphCmd = &cobra.Command{
	Use:   "graph",
	Short: "Produce a graphviz graph of objects",
	Run: func(cmd *cobra.Command, args []string) {
		buf, err := graph.GetGraph()
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(buf.String())
	},
}

func init() {
	rootCmd.AddCommand(graphCmd)
}
