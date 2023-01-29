package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var version = "0.0.1"
var commit = ""
var date = ""
var builtBy = ""

func newVersionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Print version",
		RunE: func(cmd *cobra.Command, args []string) error {
			ExecuteVersionCmd()
			return nil
		},
	}
	return cmd
}

func ExecuteVersionCmd() {
	fmt.Printf("v%s commit=%s date=%s builtBy=%s\n", version, commit, date, builtBy)
}
