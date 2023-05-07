package cmd

import (
	"fmt"
	"io"
	"log"
	"os"

	"codeberg.org/haulproject/haul/graph"
	"github.com/goccy/go-graphviz"
	"github.com/spf13/cobra"
)

// graphCmd represents the graph command
var graphCmd = &cobra.Command{
	Use:   "graph",
	Short: "Produce a graphviz graph of objects",
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

		buf, err := graph.GetGraph(graphviz.Format(format))
		if err != nil {
			log.Fatal("Error:", err)
		}

		fmt.Println(buf.String())
	},
}

func init() {
	rootCmd.AddCommand(graphCmd)

	graphCmd.Flags().String("format", "dot", "Graph output format")
}
