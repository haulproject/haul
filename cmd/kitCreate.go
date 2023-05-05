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

// kitCreateCmd represents the kitCreate command
var kitCreateCmd = &cobra.Command{
	Use:     "create KIT...",
	Aliases: []string{"add"},
	Short:   "Create kits in the database",
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var kits types.Kits

		for _, arg := range args {
			var kit types.Kit

			err := json.Unmarshal([]byte(arg), &kit)
			if err != nil {
				log.Fatalf(`Bad argument: %s

%s`, arg, err)
			}

			kits.Kits = append(kits.Kits, kit)

		}

		if len(kits.Kits) == 0 {
			os.Exit(1)
		}

		kits_bytes, err := json.Marshal(kits.Kits)
		if err != nil {
			log.Fatal("json.Marshal:", err)
		}

		result, err := api.CallWithDataB(http.MethodPost, "/v1/kit", kits_bytes)
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
			log.Fatal("Error unmarshalling POST /v1/kit:", err, "\n")
		}

		err = client.OutputObject(result_object)
		if err != nil {
			log.Fatal("Error outputting object:", err)
		}
	},
}

func init() {
	kitCmd.AddCommand(kitCreateCmd)
}
