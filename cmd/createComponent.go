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

// createComponentCmd represents the createComponent command
var createComponentCmd = &cobra.Command{
	Use:   "component",
	Short: "Components are things that can be assembled",
	Long: `Components are things that can be assembled to create servers, 
workstations, etc.

Examples: 

  - a RAM stick
  - a CPU
  - a set of speakers
  - a monitor
  - a keyboard
  - ...`,
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
	createCmd.AddCommand(createComponentCmd)
}
