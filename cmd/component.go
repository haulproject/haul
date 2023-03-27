/*
Copyright © 2023 The Haul Authors
*/
package cmd

import (
	"github.com/spf13/cobra"
)

// componentCmd represents the component command
var componentCmd = &cobra.Command{
	Use:     "component",
	Aliases: []string{"c"},
	Short:   "Components are things that can be assembled",
	Long: `Components are things that can be assembled to create servers, 
workstations, etc.

Examples: 

  - a RAM stick
  - a CPU
  - a set of speakers
  - a monitor
  - a keyboard
  - ...`,
}

func init() {
	rootCmd.AddCommand(componentCmd)
}
