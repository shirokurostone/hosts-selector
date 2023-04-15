package cmd

import (
	"bytes"
	"fmt"
	"github.com/shirokurostone/hosts-selector/lib"
	"github.com/spf13/cobra"
	"os"
)

func newToggleCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "toggle",
		Short: "toggle hosts files",
		RunE: func(cmd *cobra.Command, args []string) error {
			return ExecuteToggleCmd(config, args)
		},
	}
	return cmd
}

func ExecuteToggleCmd(config *lib.Config, names []string) error {

	affected := []*lib.HostsFile{}
	hostsfileNotFound := false
	for _, n := range names {
		match := false
		for i, _ := range config.Hosts {
			if n != config.Hosts[i].Name {
				continue
			}
			match = true
			config.Hosts[i].Enabled = !config.Hosts[i].Enabled
			affected = append(affected, &config.Hosts[i])
		}
		if !match {
			fmt.Fprintf(os.Stderr, "hostsfile not found : %s\n", n)
			hostsfileNotFound = true
		}
	}

	if hostsfileNotFound {
		return nil
	}

	if err := lib.SaveConfig(configFilePath, config); err != nil {
		return err
	}
	bs, err := os.ReadFile(config.HostsFilePath)
	if err != nil {
		return err
	}
	content := string(bs)
	buffer := &bytes.Buffer{}
	if err := lib.ReplaceHostsFile(content, buffer, config.Hosts); err != nil {
		return err
	}
	if err := os.WriteFile(config.HostsFilePath, buffer.Bytes(), 0666); err != nil {
		return err
	}

	for _, h := range affected {
		if h.Enabled {
			fmt.Printf("Enable %s\n", h.Name)
		} else {
			fmt.Printf("Disable %s\n", h.Name)
		}
	}

	return nil
}
