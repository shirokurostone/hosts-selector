package cmd

import (
	"fmt"
	"github.com/shirokurostone/hosts-selector/lib"
	"github.com/spf13/cobra"
)

func newListCmd() *cobra.Command {
	var enabled, disabled bool

	cmd := &cobra.Command{
		Use:   "list",
		Short: "show hosts files list",
		RunE: func(cmd *cobra.Command, args []string) error {

			if !enabled && !disabled {
				enabled = true
				disabled = true
			}

			return ExecuteListCmd(config, enabled, disabled)
		},
	}

	cmd.Flags().BoolVar(&enabled, "enabled", false, "")
	cmd.Flags().BoolVar(&disabled, "disabled", false, "")

	return cmd
}

func ExecuteListCmd(config *lib.Config, showEnabled bool, showDisabled bool) error {
	for _, h := range config.Hosts {
		if showEnabled && h.Enabled || showDisabled && !h.Enabled {
			fmt.Println(h.Name)
		}
	}
	return nil
}
