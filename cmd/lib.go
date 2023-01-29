package cmd

import (
	"github.com/rivo/tview"
	"github.com/shirokurostone/hosts-selector/lib"
)

func getListItemText(hosts lib.HostsFile) (string, string) {
	if hosts.Enabled {
		return tview.Escape("[X] " + hosts.Name), tview.Escape("    " + hosts.Description)
	} else {
		return "[gray]" + tview.Escape("[ ] "+hosts.Name) + "[-]", "[gray]    " + tview.Escape(hosts.Description) + "[-]"
	}
}
