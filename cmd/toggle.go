package cmd

import (
	"bytes"
	"fmt"
	"github.com/shirokurostone/hosts-selector/lib"
	"github.com/spf13/cobra"
	"os"
)

type UpdateType int

const (
	toggle UpdateType = iota
	enable
	disable
)

func newToggleCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "toggle",
		Short: "toggle hosts files",
		RunE: func(cmd *cobra.Command, args []string) error {
			return ExecuteToggleCmd(config, toggle, args)
		},
	}
	return cmd
}

func newEnableCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "enable",
		Short: "enable hosts files",
		RunE: func(cmd *cobra.Command, args []string) error {
			return ExecuteToggleCmd(config, enable, args)
		},
	}
	return cmd
}

func newDisableCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "disable",
		Short: "disable hosts files",
		RunE: func(cmd *cobra.Command, args []string) error {
			return ExecuteToggleCmd(config, disable, args)
		},
	}
	return cmd
}

func ExecuteToggleCmd(config *lib.Config, t UpdateType, names []string) error {

	affected := []*lib.Hosts{}
	hostsNotFound := false
	for _, n := range names {
		hosts := config.Hosts.SearchByName(n)
		if hosts == nil {
			fmt.Fprintf(os.Stderr, "hostsfile not found : %s\n", n)
			hostsNotFound = true
			continue
		}

		switch t {
		case toggle:
			hosts.Enabled = !hosts.Enabled
		case enable:
			hosts.Enabled = true
		case disable:
			hosts.Enabled = false
		}

		affected = append(affected, hosts)
	}

	if hostsNotFound {
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
