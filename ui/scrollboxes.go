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
	"bufio"
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/jroimartin/gocui"
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

	var lines []*string

	lines = make([]*string, 0)

	if err := scrollbox.splitText(text, lines); err != nil {
		return err
	}

	scrollbox.rows = make([]*string, len(lines))

	for index, line := range lines {
		scrollbox.rows[index] = line
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

		if _, err := fmt.Print(row); err != nil {
			return err
		}
	}

	return nil
}

func (scrollbox *TextScrollbox) splitText(text string, rows []*string) (err error) {

	var (
		reader  *strings.Reader
		scanner *bufio.Scanner
		line    string
	)

	reader = strings.NewReader(text)
	scanner = bufio.NewScanner(reader)

	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {

		line = scanner.Text()

		if err = scrollbox.splitLine(line, rows); err != nil {
			return
		}
	}

	return
}

func (scrollbox *TextScrollbox) splitLine(line string, rows []*string) (err error) {

	var (
		reader  *strings.Reader
		scanner *bufio.Scanner
		word    string
		row     string
		w, rc   int
	)

	w = scrollbox.Width()

	reader = strings.NewReader(line)
	scanner = bufio.NewScanner(reader)

	scanner.Split(bufio.ScanWords)

	for scanner.Scan() {

		word = scanner.Text()
		rc = utf8.RuneCountInString(word) + utf8.RuneCountInString(row)

		if rc > w {

			var s string

			s = row
			rows = append(rows, &s)

			row = ""
		}

		if len(row) == 0 {
			row = word
		} else {
			row = strings.Join([]string{row, word}, " ")
		}
	}

	return
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
