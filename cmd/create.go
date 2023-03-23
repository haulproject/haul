/*
Copyright © 2023 The Haul Authors
*/
package cmd

import (
	"github.com/spf13/cobra"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create an entry in the database",
}

func init() {
	rootCmd.AddCommand(createCmd)
}
