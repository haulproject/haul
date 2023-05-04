/*
 */
package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"codeberg.org/haulproject/haul/api"
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
		var assemblies []types.Assembly

		for _, arg := range args {
			var assembly types.Assembly

			err := json.Unmarshal([]byte(arg), &assembly)
			if err != nil {
				log.Println(err)
				continue
			}

			assemblies = append(assemblies, assembly)

		}

		if len(assemblies) == 0 {
			os.Exit(1)
		}

		for _, assembly := range assemblies {
			if assembly.Name == "" {
				log.Fatal("assembly.Name cannot be empty")
			}
		}

		for _, assembly := range assemblies {

			currentAssembly, err := json.Marshal(assembly)
			if err != nil {
				log.Fatal(err)
			}

			result, err := api.CallWithData(http.MethodPost, "/v1/assembly", currentAssembly)
			if err != nil {
				log.Fatal(err)
			}

			fmt.Println(result)
		}
	},
}

func init() {
	assemblyCmd.AddCommand(assemblyCreateCmd)
}
