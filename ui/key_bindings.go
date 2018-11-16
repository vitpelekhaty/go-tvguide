package ui

import (
	"github.com/jroimartin/gocui"
)

func setKeyBindings(ui *gocui.Gui) error {

	err := ui.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit)

	if err != nil {
		return err
	}

	err = ui.SetKeybinding("", gocui.KeyArrowDown, gocui.ModNone, listDown)

	if err != nil {
		return err
	}

	err = ui.SetKeybinding("", gocui.KeyArrowUp, gocui.ModNone, listUp)

	if err != nil {
		return err
	}

	err = ui.SetKeybinding("", gocui.KeyPgup, gocui.ModNone, listPageUp)

	if err != nil {
		return err
	}

	err = ui.SetKeybinding("", gocui.KeyPgdn, gocui.ModNone, listPageDown)

	if err != nil {
		return err
	}

	err = ui.SetKeybinding("", gocui.KeyTab, gocui.ModNone, switchView)

	if err != nil {
		return err
	}

	return nil
}
