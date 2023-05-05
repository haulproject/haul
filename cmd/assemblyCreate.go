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

// assemblyCreateCmd represents the assemblyCreate command
var assemblyCreateCmd = &cobra.Command{
	Use:     "create ASSEMBLY...",
	Aliases: []string{"add"},
	Short:   "Create assemblies in the database",
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var assemblies types.Assemblies

		for _, arg := range args {
			var assembly types.Assembly

			err := json.Unmarshal([]byte(arg), &assembly)
			if err != nil {
				// I believe it should crash if one of the args is bad
				log.Fatal("json.Unmarshal:", err)
			}

			assemblies.Assemblies = append(assemblies.Assemblies, assembly)

		}

		if len(assemblies.Assemblies) == 0 {
			os.Exit(1)
		}

		assemblies_bytes, err := json.Marshal(assemblies.Assemblies)
		if err != nil {
			log.Fatal("json.Marshal:", err)
		}

		result, err := api.CallWithDataB(http.MethodPost, "/v1/assembly", assemblies_bytes)
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
			log.Fatal("Error unmarshalling POST /v1/assembly:", err, "\n")
		}

		err = client.OutputObject(result_object)
		if err != nil {
			log.Fatal("Error outputting object:", err)
		}
	},
}

func init() {
	assemblyCmd.AddCommand(assemblyCreateCmd)
}
