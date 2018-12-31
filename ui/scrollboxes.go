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
	"fmt"

	"github.com/jroimartin/gocui"

	strutils "../strutils"
)

// TextScrollbox - scrollbox control
type TextScrollbox struct {
	*gocui.View

	topLineIndex int

	text string
	rows []*string
}

// CreateTextScrollbox is a constructor of TextScrollbox
func CreateTextScrollbox(view *gocui.View) *TextScrollbox {

	scrollbox := &TextScrollbox{}

	scrollbox.View = view
	scrollbox.Autoscroll = true

	return scrollbox
}

// RowCount indicates the number of rows in the scrollbox
func (scrollbox *TextScrollbox) RowCount() int {

	return len(scrollbox.rows)
}

// Empty indicates whether a scrollbox has rows or not
func (scrollbox *TextScrollbox) Empty() bool {

	return len(scrollbox.rows) == 0
}

// SetText sets text of the scrollbox and redraw the view
func (scrollbox *TextScrollbox) SetText(text string) error {

	var w int

	w = scrollbox.Width()
	sl := strutils.StringList{RightMargin: w}

	err := sl.SetText(text)

	if err != nil {
		return err
	}

	scrollbox.rows = make([]*string, sl.Count())

	for index := 0; index < len(scrollbox.rows); index++ {

		var row string

		row, err = sl.Item(index)

		if err != nil {
			return err
		}

		scrollbox.rows[index] = &row
	}

	return scrollbox.Draw()
}

// Draw draws text in the scrollbox
func (scrollbox *TextScrollbox) Draw() error {

	if scrollbox.Empty() {
		scrollbox.Clear()
		return nil
	}

	return scrollbox.DrawText(0)
}

// DrawText displays text rows from specified top line
func (scrollbox *TextScrollbox) DrawText(topidx int) error {

	var h int
	h = scrollbox.Height()

	scrollbox.Clear()
	scrollbox.topLineIndex = topidx

	rows := scrollbox.rows[topidx : h-topidx]

	for _, row := range rows {

		if _, err := fmt.Fprintln(scrollbox, row); err != nil {
			return err
		}
	}

	return nil
}

// Focused indicates whether the scrollbox has input focus
func (scrollbox *TextScrollbox) Focused() bool {

	return scrollbox.Highlight
}

// SetFocus gives the input focus to the scrollbox
func (scrollbox *TextScrollbox) SetFocus(gui *gocui.Gui) error {

	scrollbox.Highlight = true
	_, err := gui.SetCurrentView(scrollbox.Name())

	return err
}

// UnFocus removes the focus from the list box
func (scrollbox *TextScrollbox) UnFocus() {

	scrollbox.Highlight = false
}

// ResetCursor moves the cursor at the beginning of the view
func (scrollbox *TextScrollbox) ResetCursor() {

	scrollbox.SetCursor(0, 0)
}

// Reset resets items of the scrollbox and clears the view
func (scrollbox *TextScrollbox) Reset() {

	scrollbox.rows = make([]*string, 0)

	scrollbox.Clear()
	scrollbox.ResetCursor()
}

// Height returns height of the scrollbox
func (scrollbox *TextScrollbox) Height() int {

	_, h := scrollbox.Size()
	return h - 1
}

// Width returns width of the scrollbox
func (scrollbox *TextScrollbox) Width() int {

	w, _ := scrollbox.Size()
	return w - 1
}
