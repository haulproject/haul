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

// healthcheckCmd represents the healthcheck command
var healthcheckCmd = &cobra.Command{
	Use:     "healthcheck",
	Aliases: []string{"ping"},
	Short:   "Print server's healthcheck to test connection",
	Run: func(cmd *cobra.Command, args []string) {
		result, err := api.Call(api.GET, "/v1/healthcheck")
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(string(result))
	},
}

func init() {
	rootCmd.AddCommand(healthcheckCmd)
}
