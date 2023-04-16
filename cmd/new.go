package cmd

import (
	"github.com/shirokurostone/hosts-selector/lib"
	"github.com/spf13/cobra"
)

func newNewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "new",
		Short: "Create a new hosts file",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return ExecuteNewCmd(config)
		},
	}
	return cmd
}

func ExecuteNewCmd(config *lib.Config) error {

	part := lib.Hosts{
		Name:        "name",
		Description: "description",
		Content:     "",
		Url:         "",
		Enabled:     false,
	}

	result, err := EditHostsFile(part)
	if err != nil {
		return err
	}
	config.Hosts = append(config.Hosts, result)

	return lib.SaveConfig(configFilePath, config)
}
