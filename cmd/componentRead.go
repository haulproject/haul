/*
 */
package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"codeberg.org/haulproject/haul/api"
	"codeberg.org/haulproject/haul/cli"
	"codeberg.org/haulproject/haul/types"
	"github.com/spf13/cobra"
)

// componentReadCmd represents the componentRead command
var componentReadCmd = &cobra.Command{
	Use:     "read OBJECT_ID",
	Aliases: []string{"get"},
	Short:   "Prints values of component identified by OBJECT_ID",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := cli.New()

		component_bytes, err := api.Call(http.MethodGet, fmt.Sprintf("/v1/component/%s", args[0]))
		if err != nil {
			log.Fatal(err)
		}

		output, err := rootCmd.PersistentFlags().GetString("output")
		if err != nil {
			log.Fatal(err)
		}

		client.OutputStyle = output

		var component types.ComponentWithID

		err = json.Unmarshal(component_bytes, &component)
		if err != nil {
			log.Fatal(err)
		}

		err = client.OutputObject(&component)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	componentCmd.AddCommand(componentReadCmd)
}
