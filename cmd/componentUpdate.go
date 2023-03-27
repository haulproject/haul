/*
Copyright © 2023 The Haul Authors
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"log"

	"codeberg.org/haulproject/haul/api"
	"github.com/spf13/cobra"
)

var (
	id, update string
)

// componentUpdateCmd represents the componentUpdate command
var componentUpdateCmd = &cobra.Command{
	Use:     "update",
	Aliases: []string{"set"},
	Short:   "Update a component in the database",
	Long: `Update a component in the database, identified by an ObjectID, with updated fields in JSON format.

Any fields not specified will be unaffected by the update.

To empty a field, provide the zero value for the field. Note that "name" cannot be made empty.`,
	Example: `Update component identified by ObjectID 64212ede8e7046c7a1e88557, to replace all tags with "status=broken".

    $ haul component update --id '64212ede8e7046c7a1e88557' --update '{ "tags": [ "status=broken" ] }'`,
	Args: cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {

		//var component bson.E
		var component map[string]interface{}

		err := json.Unmarshal([]byte(update), &component)
		if err != nil {
			log.Fatal(err)
		}

		//fmt.Printf("%#v\n", component)

		// Empty name check

		/*
			var componentObject types.Component

			err = json.Unmarshal([]byte(update), &componentObject)
			if err != nil {
				log.Fatal(err)
			}
		*/
		/* component.Name can actually be empty

		if componentObject.Name == "" {
			log.Fatal("Component.Name cannot be empty")
		}
		*/

		//fmt.Printf("%#v\n", component.Value)
		//fmt.Printf("%#v\n", component)
		//fmt.Printf("%#v\n", componentObject)

		//TODO validate keys
		/*
			validated, err := types.ValidateFields(bson.D{bson.E{Value: component}}, types.Component{})
			if err != nil {
				log.Fatal(err)
			}

			fmt.Printf("%#v\n", validated.Map())
		*/

		currentComponent, err := json.Marshal(component)
		if err != nil {
			log.Fatal(err)
		}

		result, err := api.CallWithData(api.PUT, fmt.Sprintf("/v1/component/%s", id), currentComponent)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("%s", result)
	},
}

func init() {
	componentCmd.AddCommand(componentUpdateCmd)

	componentUpdateCmd.Flags().StringVar(&id, "id", "", "ObjectID to update")
	componentUpdateCmd.MarkFlagRequired("id")

	componentUpdateCmd.Flags().StringVar(&update, "update", "", "Data to use in the update, in JSON format")
	componentUpdateCmd.MarkFlagRequired("update")
}
