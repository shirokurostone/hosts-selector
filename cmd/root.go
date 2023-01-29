package cmd

import (
	"fmt"
	"github.com/shirokurostone/hosts-selector/lib"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
)

func newRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "",
		Short: "",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if configFilePath == "" {
				home, err := os.UserHomeDir()
				if err != nil {
					return nil
				}
				configFilePath = filepath.Join(home, ".hosts-selector")
			}

			var err error
			config, err = lib.LoadConfig(configFilePath)
			return err
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return ExecuteSelectCmd(config)
		},
	}

	cmd.PersistentFlags().StringVar(&configFilePath, "config", "", "")

	return cmd
}

var configFilePath string
var config *lib.Config

func Execute() {

	rootCmd := newRootCmd()
	rootCmd.AddCommand(newEditCmd())
	rootCmd.AddCommand(newFetchCmd())
	rootCmd.AddCommand(newNewCmd())
	rootCmd.AddCommand(newSelectCmd())
	rootCmd.AddCommand(newVersionCmd())

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
