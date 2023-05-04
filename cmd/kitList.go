/*
 */
package cmd

import (
	"encoding/json"
	"log"
	"net/http"

	"codeberg.org/haulproject/haul/api"
	"codeberg.org/haulproject/haul/cli"
	"codeberg.org/haulproject/haul/types"
	"github.com/spf13/cobra"
)

// kitListCmd represents the kitList command
var kitListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "Prints values of all kits",
	Args:    cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		client := cli.New()

		kits_bytes, err := api.Call(http.MethodGet, "/v1/kit")
		if err != nil {
			log.Fatal(err)
		}

		output, err := rootCmd.PersistentFlags().GetString("output")
		if err != nil {
			log.Fatal(err)
		}

		client.OutputStyle = output

		var kits types.KitsWithID

		if err := json.Unmarshal(kits_bytes, &kits.KitsWithID); err != nil {
			log.Fatal(err)
		}

		err = client.OutputObject(&kits)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	kitCmd.AddCommand(kitListCmd)
}
