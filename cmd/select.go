package cmd

import (
	"bytes"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/shirokurostone/hosts-selector/lib"
	"github.com/spf13/cobra"
	"os"
)

func newSelectCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "select",
		Short: "Select hosts files",
		RunE: func(cmd *cobra.Command, args []string) error {
			return ExecuteSelectCmd(config)
		},
	}
	return cmd
}

func ExecuteSelectCmd(config *lib.Config) error {
	app := tview.NewApplication()
	list := tview.NewList()
	list.SetBorder(true)

	for _, p := range config.Hosts {
		mainText, secondaryText := getListItemText(p)
		list.AddItem(mainText, secondaryText, 0, nil)
	}

	textView := tview.NewTextView()
	textView.SetBorder(true)
	if len(config.Hosts) != 0 {
		textView.SetText(tview.Escape(config.Hosts[0].Content))
	}

	list.SetSelectedFunc(func(index int, mainText string, secondaryText string, shortcut rune) {
		config.Hosts[index].Enabled = !config.Hosts[index].Enabled
		newMainText, newSecondaryText := getListItemText(config.Hosts[index])
		list.SetItemText(index, newMainText, newSecondaryText)
	})

	list.SetChangedFunc(func(index int, mainText string, secondaryText string, shortcut rune) {
		textView.SetText(tview.Escape(config.Hosts[index].Content))
	})

	okButton := tview.NewButton("OK")
	cancelButton := tview.NewButton("Cancel")

	mainFlex := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(list, 0, 1, true).
		AddItem(textView, 0, 1, false)

	footerFlex := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(okButton, 10, 100, true).
		AddItem(cancelButton, 10, 100, false)

	appFlex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(mainFlex, 0, 100, true).
		AddItem(footerFlex, 1, 0, false)

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyTab {
			if list.HasFocus() {
				app.SetFocus(textView)
			} else if textView.HasFocus() {
				app.SetFocus(okButton)
			} else if okButton.HasFocus() {
				app.SetFocus(cancelButton)
			} else if cancelButton.HasFocus() {
				app.SetFocus(list)
			}
			return nil
		}

		return event
	})

	mainFlex.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyLeft || event.Key() == tcell.KeyRight {
			if list.HasFocus() {
				app.SetFocus(textView)
			} else {
				app.SetFocus(list)
			}
			return nil
		}
		return event
	})

	footerFlex.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyLeft || event.Key() == tcell.KeyRight {
			if okButton.HasFocus() {
				app.SetFocus(cancelButton)
			} else {
				app.SetFocus(okButton)
			}
			return nil
		}
		return event
	})

	okButton.SetSelectedFunc(func() {
		app.Stop()
	})

	err := app.SetRoot(appFlex, true).SetFocus(list).Run()
	if err != nil {
		return err
	}

	if err = lib.SaveConfig(configFilePath, config); err != nil {
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

	return nil
}
