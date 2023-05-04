/*
 */
package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"codeberg.org/haulproject/haul/api"
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
		var components []types.Component

		for _, arg := range args {
			var component types.Component

			err := json.Unmarshal([]byte(arg), &component)
			if err != nil {
				log.Println(err)
				continue
			}

			components = append(components, component)

		}

		if len(components) == 0 {
			os.Exit(1)
		}

		for _, component := range components {
			if component.Name == "" {
				log.Fatal("component.Name cannot be empty")
			}
		}

		for _, component := range components {

			currentComponent, err := json.Marshal(component)
			if err != nil {
				log.Fatal(err)
			}

			result, err := api.CallWithData(http.MethodPost, "/v1/component", currentComponent)
			if err != nil {
				log.Fatal(err)
			}

			fmt.Println(result)
		}
	},
}

func init() {
	componentCmd.AddCommand(componentCreateCmd)
}
