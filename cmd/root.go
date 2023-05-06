/*
Package cmd implements cobra commands for the application.
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "haul",
	Short: "Inventory management system for patchwork components and assets.",
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.haul.yaml)")

	// api.protocol
	rootCmd.PersistentFlags().String("api-protocol", "http", "Remote api protocol (http/https) (config: 'api.protocol')")
	viper.BindPFlag("api.protocol", rootCmd.PersistentFlags().Lookup("api-protocol"))

	// api.host
	rootCmd.PersistentFlags().String("api-host", "localhost", "Remote api host (ip or dns) (config: 'api.host')")
	viper.BindPFlag("api.host", rootCmd.PersistentFlags().Lookup("api-host"))

	// api.port
	rootCmd.PersistentFlags().Int("api-port", 1315, "Remote api port (config: 'api.port')")
	viper.BindPFlag("api.port", rootCmd.PersistentFlags().Lookup("api-port"))

	// api.key
	rootCmd.PersistentFlags().String("api-key", "", "Remote api key (config: 'api.key')")
	viper.BindPFlag("api.key", rootCmd.PersistentFlags().Lookup("api-key"))

	rootCmd.PersistentFlags().StringP("output", "o", "tabby", "Output style { tabby | json | json_pretty }")
	viper.BindPFlag("cli.output", rootCmd.PersistentFlags().Lookup("output"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".haul" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".haul")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
