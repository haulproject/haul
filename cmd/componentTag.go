/*
 */
package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"codeberg.org/haulproject/haul/api"
	"github.com/spf13/cobra"
)

var (
	clear       bool
	add, remove []string
)

// componentTagCmd represents the componentTag command
var componentTagCmd = &cobra.Command{
	Use:     "tags ID",
	Aliases: []string{"t", "tag"},
	Short:   "Access and edit tags for a component",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if clear {
			// Clear tags
			log.Println("Clearing tags")

			result, err := api.Call(api.DELETE, fmt.Sprintf("/v1/component/%s/tags", args[0]))
			if err != nil {
				log.Fatalf("api.CallWithData: %s\n", err)
			}

			log.Println(string(result))

			os.Exit(0)
		}

		if len(remove) > 0 {
			// Remove some tags
			data, err := json.Marshal(remove)
			if err != nil {
				log.Fatalf("json.Marshal: %s", err)
			}

			result, err := api.CallWithData(api.POST, fmt.Sprintf("/v1/component/%s/tags/remove", args[0]), data)
			if err != nil {
				log.Fatalf("api.Call: %s\n", err)
			}

			log.Println(result)
			os.Exit(0)
		}

		if len(add) > 0 {
			// Add some tags
			data, err := json.Marshal(add)
			if err != nil {
				log.Fatalf("json.Marshal: %s", err)
			}

			result, err := api.CallWithData(api.POST, fmt.Sprintf("/v1/component/%s/tags/add", args[0]), data)
			if err != nil {
				log.Fatalf("api.Call: %s\n", err)
			}

			log.Println(result)
			os.Exit(0)
		}

		// Show tags
		result, err := api.Call(api.GET, fmt.Sprintf("/v1/component/%s/tags", args[0]))
		if err != nil {
			log.Fatalf("api.Call: %s\n", err)
		}

		var tags []string
		err = json.Unmarshal(result, &tags)
		if err != nil {
			log.Fatalf("json.Unmarshal: %s\n", err)
		}

		// JSON print
		//fmt.Println(tags)

		//fmt.Println(result)

		// Pretty print
		for _, tag := range tags {
			fmt.Println(tag)
		}

	},
}

func init() {
	componentCmd.AddCommand(componentTagCmd)

	componentTagCmd.Flags().BoolVar(&clear, "clear", false, "If set, will delete all tags in this object")
	componentTagCmd.Flags().StringSliceVar(&add, "add", nil, "List of tags to add")
	componentTagCmd.Flags().StringSliceVar(&remove, "remove", nil, "List of tags to remove")

	componentTagCmd.MarkFlagsMutuallyExclusive("clear", "add", "remove")
}
