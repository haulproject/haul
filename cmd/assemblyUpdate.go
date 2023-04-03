/*
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"log"

	"codeberg.org/haulproject/haul/api"
	"github.com/spf13/cobra"
)

// assemblyUpdateCmd represents the assemblyUpdate command
var assemblyUpdateCmd = &cobra.Command{
	Use:     "update",
	Aliases: []string{"set"},
	Short:   "Update an assembly in the database",
	Long: `Update an assembly in the database, identified by an ObjectID, with updated fields in JSON format.

Any fields not specified will be unaffected by the update.

To empty a field, provide the zero value for the field. Note that "name" cannot be made empty.`,
	Example: `Update assembly identified by ObjectID 64212ede8e7046c7a1e88557, to replace name with "Database server 01".

    $ haul assembly update --id '64212ede8e7046c7a1e88557' --update '{ "name": "Database server 01" }'`,
	Args: cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {

		var assembly map[string]interface{}

		id, err := cmd.Flags().GetString("id")
		if err != nil {
			log.Fatal(err)
		}

		update, err := cmd.Flags().GetString("update")
		if err != nil {
			log.Fatal(err)
		}

		err = json.Unmarshal([]byte(update), &assembly)
		if err != nil {
			log.Fatal(err)
		}

		currentAssembly, err := json.Marshal(assembly)
		if err != nil {
			log.Fatal(err)
		}

		result, err := api.CallWithData(api.PUT, fmt.Sprintf("/v1/assembly/%s", id), currentAssembly)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("%s", result)
	},
}

func init() {
	assemblyCmd.AddCommand(assemblyUpdateCmd)

	assemblyUpdateCmd.Flags().String("id", "", "ObjectID to update")
	assemblyUpdateCmd.MarkFlagRequired("id")

	assemblyUpdateCmd.Flags().String("update", "", "Data to use in the update, in JSON format")
	assemblyUpdateCmd.MarkFlagRequired("update")
}
