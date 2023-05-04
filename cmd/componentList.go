/*
 */
package cmd

import (
	"encoding/json"
	"log"
	"net/http"

	"codeberg.org/haulproject/haul/api"
	"codeberg.org/haulproject/haul/cli"
	"codeberg.org/haulproject/haul/types"
	"github.com/spf13/cobra"
)

// componentListCmd represents the componentList command
var componentListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "Prints values of all components",
	Args:    cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		client := cli.New()

		components_bytes, err := api.Call(http.MethodGet, "/v1/component")
		if err != nil {
			log.Fatal(err)
		}

		output, err := rootCmd.PersistentFlags().GetString("output")
		if err != nil {
			log.Fatal(err)
		}

		client.OutputStyle = output

		var components types.ComponentsWithID

		if err := json.Unmarshal(components_bytes, &components.ComponentsWithID); err != nil {
			log.Fatal(err)
		}

		err = client.OutputObject(&components)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	componentCmd.AddCommand(componentListCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// readCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// readCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
