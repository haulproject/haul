/*
Copyright © 2023 The Haul Authors
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

// kitCreateCmd represents the kitCreate command
var kitCreateCmd = &cobra.Command{
	Use:     "create KIT",
	Aliases: []string{"add"},
	Short:   "Create a kit in the database",
	Args:    cobra.ExactArgs(1),
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

		currentKit, err := json.Marshal(kits[0])
		if err != nil {
			log.Fatal(err)
		}

		result, err := api.CallWithData(api.POST, "/v1/kit", currentKit)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(result)
	},
}

func init() {
	kitCmd.AddCommand(kitCreateCmd)
}
