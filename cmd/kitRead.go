/*
 */
package cmd

import (
	"fmt"
	"log"
	"net/http"

	"codeberg.org/haulproject/haul/api"
	"codeberg.org/haulproject/haul/cli"
	"github.com/spf13/cobra"
)

// kitReadCmd represents the kitRead command
var kitReadCmd = &cobra.Command{
	Use:     "read OBJECT_ID",
	Aliases: []string{"get"},
	Short:   "Prints values of kit identified by OBJECT_ID",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := cli.New()

		kit_bytes, err := api.Call(http.MethodGet, fmt.Sprintf("/v1/kit/%s", args[0]))
		if err != nil {
			log.Fatal(err)
		}

		output, err := rootCmd.PersistentFlags().GetString("output")
		if err != nil {
			log.Fatal(err)
		}

		client.OutputStyle = output

		err = client.Output(kit_bytes)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	kitCmd.AddCommand(kitReadCmd)
}
