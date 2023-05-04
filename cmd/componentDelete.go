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

// componentDeleteCmd represents the componentDelete command
var componentDeleteCmd = &cobra.Command{
	Use:     "delete OBJECT_ID...",
	Aliases: []string{"rm", "remove", "del"},
	Short:   "Deletes components identified one or more by OBJECT_ID",
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		for _, arg := range args {
			result, err := api.Call(http.MethodDelete, fmt.Sprintf("/v1/component/%s", arg))
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(string(result))
		}
	},
}

func init() {
	componentCmd.AddCommand(componentDeleteCmd)
}
