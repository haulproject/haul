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

// kitUpdateCmd represents the kitUpdate command
var kitUpdateCmd = &cobra.Command{
	Use:     "update",
	Aliases: []string{"u", "set", "s"},
	Short:   "Update a kit in the database",
	Long: `Update a kit in the database, identified by an ObjectID, with updated fields in JSON format.

Any fields not specified will be unaffected by the update.

To empty a field, provide the zero value for the field. Note that "name" cannot be made empty.`,
	Example: `Update kit identified by ObjectID 64212ede8e7046c7a1e88557, to replace name with "Rack 01" and tags with "location=floor_5" (note: it can be better to use specific tag routes, but this remains possible).

    $ haul kit update 64212ede8e7046c7a1e88557 --data '{ "name": "Rack 01", "tags": [ "location=floor_5" ] }'`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		var kit map[string]interface{}

		id := args[0]

		update, err := cmd.Flags().GetString("data")
		if err != nil {
			log.Fatal(err)
		}

		err = json.Unmarshal([]byte(update), &kit)
		if err != nil {
			log.Fatal(err)
		}

		currentKit, err := json.Marshal(kit)
		if err != nil {
			log.Fatal(err)
		}

		result, err := api.CallWithData(http.MethodPut, fmt.Sprintf("/v1/kit/%s", id), currentKit)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("%s", result)
	},
}

func init() {
	kitCmd.AddCommand(kitUpdateCmd)

	kitUpdateCmd.Flags().String("data", "", "Data to use in the update, in JSON format")
	kitUpdateCmd.MarkFlagRequired("data")
}
