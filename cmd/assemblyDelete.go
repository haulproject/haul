/*
*/
package cmd

import (
	"fmt"
	"log"

	"codeberg.org/haulproject/haul/api"
	"github.com/spf13/cobra"
)

// assemblyDeleteCmd represents the assemblyDelete command
var assemblyDeleteCmd = &cobra.Command{
	Use:     "delete OBJECT_ID",
	Aliases: []string{"rm", "remove", "del"},
	Short:   "Deletes assembly identified by OBJECT_ID",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		result, err := api.Call(api.DELETE, fmt.Sprintf("/v1/assembly/%s", args[0]))
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(result))
	},
}

func init() {
	assemblyCmd.AddCommand(assemblyDeleteCmd)
}
