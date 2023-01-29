package cmd

import (
	"github.com/shirokurostone/hosts-selector/lib"
	"github.com/spf13/cobra"
	"io"
	"net/http"
)

func newFetchCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "fetch",
		Short: "Fetch remote hosts files",
		RunE: func(cmd *cobra.Command, args []string) error {
			return ExecuteFetchCmd(config)
		},
	}
	return cmd
}

func ExecuteFetchCmd(config *lib.Config) error {
	for i := range config.Hosts {
		if config.Hosts[i].Url != "" {
			resp, err := http.Get(config.Hosts[i].Url)
			if err != nil {
				return err
			}
			defer resp.Body.Close()

			b, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			config.Hosts[i].Content = string(b)
		}
	}
	return lib.SaveConfig(configFilePath, config)
}
