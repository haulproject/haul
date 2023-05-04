/*
 */
package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"codeberg.org/haulproject/haul/api"
	"github.com/spf13/cobra"
)

// componentUpdateCmd represents the componentUpdate command
var componentUpdateCmd = &cobra.Command{
	Use:     "update OBJECT_ID",
	Aliases: []string{"u", "set", "s"},
	Short:   "Update a component in the database",
	Long: `Update a component in the database, identified by an ObjectID, with updated fields in JSON format.

Any fields not specified will be unaffected by the update.

To empty a field, provide the zero value for the field. Note that "name" cannot be made empty.`,
	Example: `Update component identified by ObjectID 64212ede8e7046c7a1e88557, to set status to "broken"

    $ haul component update 64212ede8e7046c7a1e88557 --data '{ "status":"broken" }'`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		var component map[string]interface{}

		id := args[0]

		update, err := cmd.Flags().GetString("data")
		if err != nil {
			log.Fatal(err)
		}

		err = json.Unmarshal([]byte(update), &component)
		if err != nil {
			log.Fatal(err)
		}

		currentComponent, err := json.Marshal(component)
		if err != nil {
			log.Fatal(err)
		}

		result, err := api.CallWithData(http.MethodPut, fmt.Sprintf("/v1/component/%s", id), currentComponent)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("%s", result)
	},
}

func init() {
	componentCmd.AddCommand(componentUpdateCmd)

	componentUpdateCmd.Flags().String("data", "", "Data to use in the update, in JSON format")
	componentUpdateCmd.MarkFlagRequired("data")
}
