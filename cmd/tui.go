package cmd

import (
	"fmt"

	tea "charm.land/bubbletea/v2"
	"github.com/spf13/cobra"
)

var tuiCmd = &cobra.Command{
	Use:   "tui",
	Short: "Launch interactive TUI",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("TUI is not yet implemented.")
		fmt.Println("Coming in Layer 4.")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(tuiCmd)
}

// newTUIProgram creates a new Bubble Tea program.
// Will be used when TUI is implemented in Layer 4.
func newTUIProgram() *tea.Program {
	// TODO: Implement TUI model in Layer 4
	return tea.NewProgram(nil)
}
