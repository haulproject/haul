package cmd

import (
	"fmt"
	"log"

	"codeberg.org/haulproject/haul/api"
	"github.com/spf13/cobra"
)

// apiCmd represents the api command
var apiCmd = &cobra.Command{
	Use:   "api",
	Short: "Obtain the api routes on the haul instance",
	Run: func(cmd *cobra.Command, args []string) {
		result, err := api.Call(api.GET, "/v1")
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(string(result))
	},
}

func init() {
	rootCmd.AddCommand(apiCmd)
}
