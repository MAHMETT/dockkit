package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	version   = "dev"
	cfgFile   string
	verbose   bool
	noColor   bool
)

var rootCmd = &cobra.Command{
	Use:   "dockkit",
	Short: "Docker development infrastructure manager",
	Long:  `dockkit is a CLI TUI tool for managing Docker-based development infrastructure services.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func SetVersion(v string) {
	version = v
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default: ~/.config/dockkit/config.yaml)")
	rootCmd.PersistentFlags().BoolVar(&verbose, "verbose", false, "enable debug logging")
	rootCmd.PersistentFlags().BoolVar(&noColor, "no-color", false, "disable colored output")
}

func initConfig() {
	// TODO: Initialize config with Viper
}
