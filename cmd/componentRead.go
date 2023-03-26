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
	Use:   "read [OBJECTID]",
	Short: "Prints values of component identified by ObjectID in args",
	Args:  cobra.ExactArgs(1),
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

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// readCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// readCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
