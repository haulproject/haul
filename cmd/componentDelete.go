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

// componentDeleteCmd represents the componentDelete command
var componentDeleteCmd = &cobra.Command{
	Use:     "delete OBJECT_ID",
	Aliases: []string{"rm", "remove"},
	Short:   "Deletes component identified by OBJECT_ID",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		result, err := api.Call(api.DELETE, fmt.Sprintf("/v1/component/%s", args[0]))
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(result))
	},
}

func init() {
	componentCmd.AddCommand(componentDeleteCmd)
}
