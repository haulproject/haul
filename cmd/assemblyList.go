/*
Copyright © 2023 The Haul Authors
*/
package cmd

import (
	"fmt"
	"log"

	"codeberg.org/haulproject/haul/api"
	"github.com/spf13/cobra"
)

// assemblyListCmd represents the assemblyList command
var assemblyListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "Prints values of all assemblies",
	Args:    cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		result, err := api.Call(api.GET, "/v1/assembly")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(result))
	},
}

func init() {
	assemblyCmd.AddCommand(assemblyListCmd)
}
