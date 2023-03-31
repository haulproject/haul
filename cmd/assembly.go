/*
Copyright © 2023 The Haul Authors
*/
package cmd

import (
	"github.com/spf13/cobra"
)

// assemblyCmd represents the assembly command
var assemblyCmd = &cobra.Command{
	Use:     "assembly",
	Aliases: []string{"a"},
	Short:   "Assemblies are things assembled from components",
}

func init() {
	rootCmd.AddCommand(assemblyCmd)
}
