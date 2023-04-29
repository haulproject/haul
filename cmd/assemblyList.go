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

// assemblyListCmd represents the assemblyList command
var assemblyListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "Prints values of all assemblies",
	Args:    cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		client := cli.New()

		assemblies_bytes, err := api.Call(http.MethodGet, "/v1/assembly")
		if err != nil {
			log.Fatal(err)
		}

		output, err := rootCmd.PersistentFlags().GetString("output")
		if err != nil {
			log.Fatal(err)
		}

		client.OutputStyle = output

		var assemblies types.AssembliesWithID

		if err := json.Unmarshal(assemblies_bytes, &assemblies.AssembliesWithID); err != nil {
			log.Fatal(err)
		}

		err = client.OutputObject(&assemblies)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	assemblyCmd.AddCommand(assemblyListCmd)
}
