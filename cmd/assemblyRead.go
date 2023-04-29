/*
 */
package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"codeberg.org/haulproject/haul/api"
	"codeberg.org/haulproject/haul/cli"
	"codeberg.org/haulproject/haul/types"
	"github.com/spf13/cobra"
)

// assemblyReadCmd represents the assemblyRead command
var assemblyReadCmd = &cobra.Command{
	Use:     "read OBJECT_ID",
	Aliases: []string{"get"},
	Short:   "Prints values of assembly identified by OBJECT_ID",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := cli.New()

		assembly_bytes, err := api.Call(http.MethodGet, fmt.Sprintf("/v1/assembly/%s", args[0]))
		if err != nil {
			log.Fatal(err)
		}

		output, err := rootCmd.PersistentFlags().GetString("output")
		if err != nil {
			log.Fatal(err)
		}

		client.OutputStyle = output

		var assembly types.AssemblyWithID

		err = json.Unmarshal(assembly_bytes, &assembly)
		if err != nil {
			log.Fatal(err)
		}

		err = client.OutputObject(&assembly)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	assemblyCmd.AddCommand(assemblyReadCmd)
}
