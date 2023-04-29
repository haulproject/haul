/*
 */
package cmd

import (
	"github.com/spf13/cobra"
)

// kitCmd represents the kit command
var kitCmd = &cobra.Command{
	Use:     "kit",
	Aliases: []string{"k"},
	Short:   "Kits are groups of components and assemblies",
}

func init() {
	rootCmd.AddCommand(kitCmd)
}
