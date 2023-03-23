/*
Copyright © 2023 The Haul Authors
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"codeberg.org/haulproject/haul/types"
	"github.com/spf13/cobra"
)

// componentCreateCmd represents the componentCreate command
var componentCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a component in the database",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			log.Fatal("Must specify at least 1 component to create")
		}

		var components []types.Component

		for _, arg := range args {
			var component types.Component

			err := json.Unmarshal([]byte(arg), &component)
			if err != nil {
				log.Println(err)
				continue
			}

			components = append(components, component)

		}
		result, err := json.Marshal(components)
		if err != nil {
			log.Fatal(err)
		}

		for _, component := range components {
			if component.Name == "" {
				os.Stderr.WriteString("component.Name cannot be empty")
			}
		}
		fmt.Printf("%s\n", string(result))
	},
}

func init() {
	componentCmd.AddCommand(componentCreateCmd)
}
