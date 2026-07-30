package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all configured services",
	RunE: func(cmd *cobra.Command, args []string) error {
		if cfg == nil || len(cfg.Services) == 0 {
			fmt.Println("No services configured.")
			fmt.Println("Run 'dockkit init' to create default config, or use TUI to add services.")
			return nil
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "SERVICE\tVERSION\tPORT\tCONTAINER\tENABLED")
		fmt.Fprintln(w, "-------\t-------\t----\t---------\t-------")

		for name, svc := range cfg.Services {
			for ver, v := range svc.Versions {
				enabled := "no"
				if v.Enabled {
					enabled = "yes"
				}
				fmt.Fprintf(w, "%s\t%s\t%d\t%s\t%s\n", name, ver, v.Port, v.ContainerName, enabled)
			}
		}

		w.Flush()
		return nil
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
