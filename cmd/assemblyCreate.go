/*
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"codeberg.org/haulproject/haul/api"
	"codeberg.org/haulproject/haul/types"
	"github.com/spf13/cobra"
)

// assemblyCreateCmd represents the assemblyCreate command
var assemblyCreateCmd = &cobra.Command{
	Use:     "create ASSEMBLY",
	Aliases: []string{"add"},
	Short:   "Create an assembly in the database",
	Args:    cobra.ExactArgs(1),
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

		currentAssembly, err := json.Marshal(assemblies[0])
		if err != nil {
			log.Fatal(err)
		}

		result, err := api.CallWithData(api.POST, "/v1/assembly", currentAssembly)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(result)
	},
}

func init() {
	assemblyCmd.AddCommand(assemblyCreateCmd)
}
