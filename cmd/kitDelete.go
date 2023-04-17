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
	Use:     "delete OBJECT_ID",
	Aliases: []string{"rm", "remove", "del"},
	Short:   "Deletes kit identified by OBJECT_ID",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		result, err := api.Call(http.MethodDelete, fmt.Sprintf("/v1/kit/%s", args[0]))
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(result))
	},
}

func init() {
	kitCmd.AddCommand(kitDeleteCmd)
}
