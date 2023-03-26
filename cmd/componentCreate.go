/*
Copyright © 2023 The Haul Authors
*/
package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"codeberg.org/haulproject/haul/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// componentCreateCmd represents the componentCreate command
var componentCreateCmd = &cobra.Command{
	Use:   "create [COMPONENT]",
	Short: "Create a component in the database",
	Long: `Create a component in the database, using the COMPONENT defined in args in JSON format.

The "name" field must be non-blank, but its value can be any string. Examples of "name" include a description of the component, or something like a serial number, mac address, or other identifier. It currently does not need to be unique in the database.

The "tags" field is non-mandatory. It can however be used to convey more detailed information about the component.`,

	Example: `Create new 8gb RAM stick, with the "name" field used for a simple description

    $ haul component create '{ "name": "Generic 8gb RAM", "tags": [ "manufacturer=generic", "type=ram", "size=8gb" ] }'

Create a new set of speakers without any tags

    $ haul component create '{ "name": "Speakers" }'`,
	Args: cobra.ExactArgs(1),
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

		_, err := json.Marshal(components)
		if err != nil {
			log.Fatal(err)
		}

		for _, component := range components {
			if component.Name == "" {
				log.Fatal("component.Name cannot be empty")
			}
		}

		if err != nil {
			log.Fatal(err)
		}

		currentComponent, err := json.Marshal(components[0])
		if err != nil {
			log.Fatal(err)
		}

		endpoint := fmt.Sprintf("%s://%s:%d",
			viper.GetString("api.protocol"),
			viper.GetString("api.host"),
			viper.GetInt("api.port"),
		)
		request := fmt.Sprintf("%s%s", endpoint, "/v1/component")

		resp, err := http.Post(request, "application/json",
			bytes.NewBuffer(currentComponent))

		if err != nil {
			log.Fatal(err)
		}

		var res map[string]interface{}

		json.NewDecoder(resp.Body).Decode(&res)

		fmt.Println(res["message"])
	},
}

func init() {
	componentCmd.AddCommand(componentCreateCmd)
}
