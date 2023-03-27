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

// componentReadCmd represents the componentRead command
var componentReadCmd = &cobra.Command{
	Use:     "read OBJECT_ID",
	Aliases: []string{"get"},
	Short:   "Prints values of component identified by OBJECT_ID",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		result, err := api.Call(api.GET, fmt.Sprintf("/v1/component/%s", args[0]))
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(result))
	},
}

func init() {
	componentCmd.AddCommand(componentReadCmd)
}
