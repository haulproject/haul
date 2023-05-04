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

// kitCreateCmd represents the kitCreate command
var kitCreateCmd = &cobra.Command{
	Use:     "create KIT...",
	Aliases: []string{"add"},
	Short:   "Create kits in the database",
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var kits []types.Kit

		for _, arg := range args {
			var kit types.Kit

			err := json.Unmarshal([]byte(arg), &kit)
			if err != nil {
				log.Println(err)
				continue
			}

			kits = append(kits, kit)

		}

		if len(kits) == 0 {
			os.Exit(1)
		}

		for _, kit := range kits {
			if kit.Name == "" {
				log.Fatal("kit.Name cannot be empty")
			}
		}

		for _, kit := range kits {

			currentKit, err := json.Marshal(kit)
			if err != nil {
				log.Fatal(err)
			}

			result, err := api.CallWithData(http.MethodPost, "/v1/kit", currentKit)
			if err != nil {
				log.Fatal(err)
			}

			fmt.Println(result)
		}
	},
}

func init() {
	kitCmd.AddCommand(kitCreateCmd)
}
