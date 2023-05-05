/*
 */
package cmd

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"codeberg.org/haulproject/haul/api"
	"codeberg.org/haulproject/haul/cli"
	"codeberg.org/haulproject/haul/types"
	"github.com/spf13/cobra"
)

// componentCreateCmd represents the componentCreate command
var componentCreateCmd = &cobra.Command{
	Use:     "create COMPONENT...",
	Aliases: []string{"add"},
	Short:   "Create components in the database",
	Long: `Create a component in the database, using the COMPONENT defined in args in JSON format.

	Multiple components can be created by giving splitting them as individual arguments.

The "name" field must be non-blank, but its value can be any string. Examples of "name" include a description of the component, or something like a serial number, mac address, or other identifier. It currently does not need to be unique in the database.

The "tags" field is non-mandatory. It can however be used to convey more detailed information about the component.`,

	Example: `Create new 8gb RAM stick, with the "name" field used for a simple description

    $ haul component create '{ "name": "Generic 8gb RAM", "tags": [ "manufacturer=generic", "type=ram", "size=8gb" ] }'

Create a new set of speakers without any tags

    $ haul component create '{ "name": "Speakers" }'`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var components types.Components

		for _, arg := range args {
			var component types.Component

			err := json.Unmarshal([]byte(arg), &component)
			if err != nil {
				// I believe it should crash if one of the args is bad
				log.Fatal("json.Unmarshal:", err)
			}

			components.Components = append(components.Components, component)

		}

		if len(components.Components) == 0 {
			os.Exit(1)
		}

		components_bytes, err := json.Marshal(components.Components)
		if err != nil {
			log.Fatal("json.Marshal:", err)
		}

		result, err := api.CallWithDataB(http.MethodPost, "/v1/component", components_bytes)
		if err != nil {
			log.Fatal("api.CallWithDataB:", err)
		}

		// Using cli object
		client := cli.New()

		output, err := rootCmd.PersistentFlags().GetString("output")
		if err != nil {
			log.Fatal(err)
		}

		client.OutputStyle = output

		var result_object types.InsertResult

		err = json.Unmarshal(result, &result_object)
		if err != nil {
			log.Fatal("Error unmarshalling POST /v1/component:", err, "\n")
		}

		err = client.OutputObject(result_object)
		if err != nil {
			log.Fatal("Error outputting object:", err)
		}
	},
}

func init() {
	componentCmd.AddCommand(componentCreateCmd)
}
