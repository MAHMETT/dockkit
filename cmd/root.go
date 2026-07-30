package cmd

import (
	"fmt"
	"os"

	"github.com/MAHMETT/dockkit/internal/config"
	"github.com/spf13/cobra"
)

var (
	version   = "dev"
	cfgFile   string
	verbose   bool
	noColor   bool
	cfg       *config.Config
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

// GetConfig returns the loaded config (available after initConfig).
func GetConfig() *config.Config {
	return cfg
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default: ~/.config/dockkit/config.yaml)")
	rootCmd.PersistentFlags().BoolVar(&verbose, "verbose", false, "enable debug logging")
	rootCmd.PersistentFlags().BoolVar(&noColor, "no-color", false, "disable colored output")
}

func initConfig() {
	var err error
	if cfgFile != "" {
		cfg, err = config.LoadFromFile(cfgFile)
	} else {
		cfg, err = config.Load()
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: failed to load config: %v\n", err)
		cfg = config.DefaultConfig()
	}
}
