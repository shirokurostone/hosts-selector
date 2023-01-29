package cmd

import (
	"github.com/rivo/tview"
	"github.com/shirokurostone/hosts-selector/lib"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
)

func newEditCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "edit",
		Short: "Edit a hosts file",
		Args:  cobra.RangeArgs(0, 1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := ""
			if len(args) != 0 {
				name = args[0]
			}
			return ExecuteEditCmd(config, name)
		},
	}
	return cmd
}

func ExecuteEditCmd(config *lib.Config, name string) error {
	var err error
	index := -1

	if name == "" {
		index, err = SelectHostsFile(config)
		if err != nil {
			return err
		}
	} else {
		index = config.SearchHostsFileName(name)
	}

	result, err := EditHostsFile(config.Hosts[index])
	if err != nil {
		return err
	}
	config.Hosts[index] = result

	return lib.SaveConfig(configFilePath, config)
}

func EditHostsFile(part lib.HostsFile) (lib.HostsFile, error) {
	f, err := os.CreateTemp("", "hosts")
	if err != nil {
		return lib.HostsFile{}, err
	}
	defer os.Remove(f.Name())

	err = os.WriteFile(f.Name(), []byte(lib.Marshal(part)), 0600)
	if err != nil {
		return lib.HostsFile{}, err
	}

	c := exec.Command("vim", f.Name())
	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	if err = c.Run(); err != nil {
		return lib.HostsFile{}, err
	}

	b, err := os.ReadFile(f.Name())
	if err != nil {
		return lib.HostsFile{}, err
	}

	buffer := string(b)
	return lib.Unmarshal(buffer)
}

func SelectHostsFile(config *lib.Config) (int, error) {
	app := tview.NewApplication()
	list := tview.NewList()

	for _, p := range config.Hosts {
		mainText, secondaryText := getListItemText(p)
		list.AddItem(mainText, secondaryText, 0, nil)
	}

	textView := tview.NewTextView()
	if len(config.Hosts) != 0 {
		textView.SetText(tview.Escape(config.Hosts[0].Content))
	}
	selected := -1

	list.SetSelectedFunc(func(index int, mainText string, secondaryText string, shortcut rune) {
		selected = index
		app.Stop()
	})

	list.SetChangedFunc(func(index int, mainText string, secondaryText string, shortcut rune) {
		textView.SetText(tview.Escape(config.Hosts[index].Content))
	})

	err := app.SetRoot(list, true).SetFocus(list).Run()
	if err != nil {
		return 0, err
	}

	return selected, err
}
