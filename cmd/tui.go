package cmd

import (
	"fmt"

	tea "charm.land/bubbletea/v2"
	"github.com/spf13/cobra"

	tui "github.com/MAHMETT/dockkit/internal/tui"
)

var tuiCmd = &cobra.Command{
	Use:   "tui",
	Short: "Launch interactive TUI",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Starting dockkit TUI...")
		fmt.Println("(Press ? for help, Ctrl+C to quit)")
		fmt.Println()

		// Create and run TUI
		p := tea.NewProgram(tui.NewModel())
		if _, err := p.Run(); err != nil {
			return fmt.Errorf("TUI error: %w", err)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(tuiCmd)
}
