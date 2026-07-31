package cmd

import (
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
	"github.com/spf13/cobra"

	"github.com/MAHMETT/dockkit/internal/config"
	tui "github.com/MAHMETT/dockkit/internal/tui"
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
	RunE: func(cmd *cobra.Command, args []string) error {
		// Default action: launch TUI
		fmt.Println("Starting dockkit TUI...")
		fmt.Println("(Press ? for help, Ctrl+C to quit)")
		fmt.Println()

		p := tea.NewProgram(tui.NewModel())
		if _, err := p.Run(); err != nil {
			return fmt.Errorf("TUI error: %w", err)
		}
		return nil
	},
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
