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
	"github.com/spf13/cobra"
)

// componentTargetCmd represents the componentTarget command
var componentTargetCmd = &cobra.Command{
	Use:   "target ID",
	Short: "Access and edit target for a component",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var (
			clear bool
			set   string
		)

		//TODO verify errors
		clear, _ = cmd.Flags().GetBool("clear")
		set, _ = cmd.Flags().GetString("set")

		if clear {
			// Clear target
			log.Println("Clearing target")

			result, err := api.Call(http.MethodDelete, fmt.Sprintf("/v1/component/%s/target", args[0]))
			if err != nil {
				log.Fatalf("api.CallWithData: %s\n", err)
			}

			log.Println(string(result))

			os.Exit(0)
		}

		if set != "" {
			// Add some target
			data, err := json.Marshal(set)
			if err != nil {
				log.Fatalf("json.Marshal: %s", err)
			}

			result, err := api.CallWithData(http.MethodPost, fmt.Sprintf("/v1/component/%s/target", args[0]), data)
			if err != nil {
				log.Fatalf("api.Call: %s\n", err)
			}

			log.Println(result)
			os.Exit(0)
		}

		// Show target
		result, err := api.Call(http.MethodGet, fmt.Sprintf("/v1/component/%s/target", args[0]))
		if err != nil {
			log.Fatalf("api.Call: %s\n", err)
		}

		var target string

		err = json.Unmarshal(result, &target)
		if err != nil {
			log.Fatalf("json.Unmarshal: %s\nresult: %s", err, string(result))
		}

		fmt.Println(target)
	},
}

func init() {
	componentCmd.AddCommand(componentTargetCmd)

	componentTargetCmd.Flags().String("set", "", "Set this object's target object")

	componentTargetCmd.Flags().Bool("clear", false, "If set, will clear target for this object")

	componentTargetCmd.MarkFlagsMutuallyExclusive("clear", "set")
}
