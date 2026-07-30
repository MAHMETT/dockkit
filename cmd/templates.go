package cmd

import (
	"fmt"

	"github.com/MAHMETT/dockkit/internal/templates"
	"github.com/spf13/cobra"
)

var templatesCmd = &cobra.Command{
	Use:   "templates",
	Short: "List available service templates",
	RunE: func(cmd *cobra.Command, args []string) error {
		names, err := templates.ListBuiltin()
		if err != nil {
			return fmt.Errorf("listing templates: %w", err)
		}

		fmt.Println("Available templates:")
		for _, name := range names {
			tmpl, err := templates.LoadBuiltin(name)
			if err != nil {
				fmt.Printf("  %s (error: %v)\n", name, err)
				continue
			}
			versions := make([]string, len(tmpl.Versions))
			for i, v := range tmpl.Versions {
				versions[i] = v.Key
			}
			fmt.Printf("  %s %s — %s\n", tmpl.Icon, tmpl.Name, tmpl.Description)
			fmt.Printf("    versions: %v\n", versions)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(templatesCmd)
}
