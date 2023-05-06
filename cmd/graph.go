package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// graphCmd represents the graph command
var graphCmd = &cobra.Command{
	Use:   "graph",
	Short: "Produce a graphviz graph of objects",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("graph called")
		//TODO graphviz from api calls
	},
}

func init() {
	rootCmd.AddCommand(graphCmd)
}
