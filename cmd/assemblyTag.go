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

// assemblyTagCmd represents the assemblyTag command
var assemblyTagCmd = &cobra.Command{
	Use:     "tags ID",
	Aliases: []string{"t", "tag"},
	Short:   "Access and edit tags for a assembly",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var (
			clear       bool
			add, remove []string
		)

		//TODO verify errors
		clear, _ = cmd.Flags().GetBool("clear")
		add, _ = cmd.Flags().GetStringSlice("add")
		remove, _ = cmd.Flags().GetStringSlice("remove")

		if clear {
			// Clear tags
			log.Println("Clearing tags")

			result, err := api.Call(http.MethodDelete, fmt.Sprintf("/v1/assembly/%s/tags", args[0]))
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

			result, err := api.CallWithData(http.MethodPost, fmt.Sprintf("/v1/assembly/%s/tags/remove", args[0]), data)
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

			result, err := api.CallWithData(http.MethodPost, fmt.Sprintf("/v1/assembly/%s/tags/add", args[0]), data)
			if err != nil {
				log.Fatalf("api.Call: %s\n", err)
			}

			log.Println(result)
			os.Exit(0)
		}

		// Show tags
		result, err := api.Call(http.MethodGet, fmt.Sprintf("/v1/assembly/%s/tags", args[0]))
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
	assemblyCmd.AddCommand(assemblyTagCmd)

	assemblyTagCmd.Flags().Bool("clear", false, "If set, will delete all tags in this object")
	assemblyTagCmd.Flags().StringSlice("add", nil, "List of tags to add")
	assemblyTagCmd.Flags().StringSlice("remove", nil, "List of tags to remove")

	assemblyTagCmd.MarkFlagsMutuallyExclusive("clear", "add", "remove")
}
