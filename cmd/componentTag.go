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
			log.Println("Clearing tags")

			result, err := api.CallWithData(api.PUT, fmt.Sprintf("/v1/component/%s", args[0]), []byte("{\"tags\":[]}"))
			if err != nil {
				log.Fatalf("api.CallWithData: %s\n", err)
			}

			log.Println(result)

			os.Exit(0)
		}

		if len(remove) > 0 {
			update := struct {
				Tags []string `json:"tags"`
			}{
				Tags: []string{},
			}

			result, err := api.Call(api.GET, fmt.Sprintf("/v1/component/%s", args[0]))
			if err != nil {
				log.Fatalf("api.Call: %s\n", err)
			}

			var component types.Component
			err = json.Unmarshal(result, &component)
			if err != nil {
				log.Fatalf("bson.Unmarshal: %s\n", err)
			}

			for _, tag := range component.Tags {
				taggedForRemove := false
				for _, remove_tag := range remove {
					if remove_tag == tag {
						taggedForRemove = true
						break
					}
				}

				if !taggedForRemove {
					update.Tags = append(update.Tags, tag)
				}
			}

			data, err := json.Marshal(&update)
			if err != nil {
				log.Fatalf("json.Marshal: %s\n", err)
			}

			putResult, err := api.CallWithData(http.MethodPut, fmt.Sprintf("/v1/component/%s", args[0]), data)
			if err != nil {
				log.Fatalf("api.CallWithData: %s\n", err)
			}

			fmt.Println(putResult)

			os.Exit(0)
		}

		if len(add) > 0 {

			result, err := api.Call(api.GET, fmt.Sprintf("/v1/component/%s", args[0]))
			if err != nil {
				log.Fatalf("api.Call: %s\n", err)
			}

			var component types.Component
			err = json.Unmarshal(result, &component)
			if err != nil {
				log.Fatalf("bson.Unmarshal: %s\n", err)
			}

			update := struct {
				Tags []string `json:"tags"`
			}{
				Tags: component.Tags,
			}

			for _, add_tag := range add {
				present := false
				for _, tag := range component.Tags {
					if add_tag == tag {
						present = true
						break
					}
				}

				if !present {
					update.Tags = append(update.Tags, add_tag)
				}
			}

			data, err := json.Marshal(&update)
			if err != nil {
				log.Fatalf("json.Marshal: %s\n", err)
			}

			putResult, err := api.CallWithData(http.MethodPut, fmt.Sprintf("/v1/component/%s", args[0]), data)
			if err != nil {
				log.Fatalf("api.CallWithData: %s\n", err)
			}

			fmt.Println(putResult)

			os.Exit(0)
		}

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
