// IPTV guide viewer
//
// Copyright 2018 Vitaly Pelekhaty
// Based on https://github.com/antavelos/terminews/blob/master/ui.go
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

// TextScrollbox - scrollbox control
type TextScrollbox struct {
	*gocui.View

	rows []string
}

// CreateTextScrollbox is a constructor of TextScrollbox
func CreateTextScrollbox(view *gocui.View) *TextScrollbox {

	scrollbox := &TextScrollbox{}

	scrollbox.View = view
	scrollbox.Autoscroll = true

	return scrollbox
}

func (scrollbox *TextScrollbox) RowCount() int {

	return len(scrollbox.rows)
}

func (scrollbox *TextScrollbox) Empty() bool {

	return len(scrollbox.rows) == 0
}

func (scrollbox *TextScrollbox) SetText(text string) error {

	return nil
}

func (scrollbox *TextScrollbox) Focused() bool {

	return scrollbox.Highlight
}

func (scrollbox *TextScrollbox) SetFocus(gui *gocui.Gui) error {

	scrollbox.Highlight = true
	_, err := gui.SetCurrentView(scrollbox.Name())

	return err
}

func (scrollbox *TextScrollbox) UnFocus() {

	scrollbox.Highlight = false
}

func (scrollbox *TextScrollbox) ResetCursor() {

	scrollbox.SetCursor(0, 0)
}

func (scrollbox *TextScrollbox) Reset() {

	scrollbox.rows = make([]string, 0)

	scrollbox.Clear()
	scrollbox.ResetCursor()
}

func (scrollbox *TextScrollbox) Height() int {

	_, h := scrollbox.Size()
	return h - 1
}

func (scrollbox *TextScrollbox) Width() int {

	w, _ := scrollbox.Size()
	return w - 1
}
