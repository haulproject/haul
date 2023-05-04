/*
 */
package cmd

import (
	"fmt"
	"log"
	"net/http"

	"codeberg.org/haulproject/haul/api"
	"github.com/spf13/cobra"
)

// kitDeleteCmd represents the kitDelete command
var kitDeleteCmd = &cobra.Command{
	Use:     "delete OBJECT_ID...",
	Aliases: []string{"rm", "remove", "del"},
	Short:   "Deletes kits identified by one or more OBJECT_ID",
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		for _, arg := range args {
			result, err := api.Call(http.MethodDelete, fmt.Sprintf("/v1/kit/%s", arg))
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(string(result))
		}
	},
}

func init() {
	kitCmd.AddCommand(kitDeleteCmd)
}
