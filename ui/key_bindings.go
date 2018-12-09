// IPTV guide viewer
//
// Copyright 2018 Vitaly Pelekhaty
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
// either express or implied. See the License for the specific
// language governing permissions and limitations under the License.

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

	err = ui.SetKeybinding("", gocui.KeyF1, gocui.ModNone, help)

	if err != nil {
		return err
	}

	err = ui.SetKeybinding("", gocui.KeyCtrlQ, gocui.ModNone, destroyTopView)

	if err != nil {
		return err
	}

	err = ui.SetKeybinding("", gocui.KeyCtrl6, gocui.ModNone, focusGroupsView)

	if err != nil {
		return err
	}

	err = ui.SetKeybinding("", gocui.KeyCtrl7, gocui.ModNone, focusChannelsView)

	if err != nil {
		return err
	}

	err = ui.SetKeybinding("", gocui.KeyCtrl8, gocui.ModNone, focusGuideView)

	if err != nil {
		return err
	}

	err = ui.SetKeybinding("", gocui.KeyEnter, gocui.ModNone, onEnter)

	if err != nil {
		return err
	}

	return nil
}
