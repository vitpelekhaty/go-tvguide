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
	"unicode/utf8"

	"github.com/jroimartin/gocui"
)

// Page holdes info about a list based view
type Page struct {
	offset, limit int
}

// OnGetTextEvent - an event that fires if you want to get the item text
type OnGetTextEvent func(view *gocui.View, item interface{}) string

// VirtualListBox ui control
type VirtualListBox struct {
	*gocui.View

	title     string
	items     []interface{}
	pages     []Page
	pageIndex int
	ordered   bool

	OnGetText OnGetTextEvent
}

// CreateVirtualListBox is a constructor for VirtualListBox
func CreateVirtualListBox(view *gocui.View, ordered bool) *VirtualListBox {

	listbox := &VirtualListBox{}

	listbox.View = view
	listbox.SelBgColor = gocui.ColorBlack
	listbox.SelFgColor = gocui.ColorWhite | gocui.AttrBold
	listbox.Autoscroll = true
	listbox.ordered = ordered

	return listbox
}

// Count indicates the number of items in the list box
func (listbox *VirtualListBox) Count() int {
	return len(listbox.items)
}

// Empty indicates whether a list has items or not
func (listbox *VirtualListBox) Empty() bool {
	return (listbox.Count() == 0)
}

// Focused indicates whether the list box has input focus
func (listbox *VirtualListBox) Focused() bool {
	return listbox.Highlight
}

// SetFocus gives the input focus to the list box
func (listbox *VirtualListBox) SetFocus(gui *gocui.Gui) error {

	listbox.Highlight = true
	_, err := gui.SetCurrentView(listbox.Name())

	return err
}

// UnFocus removes the focus from the list box
func (listbox *VirtualListBox) UnFocus() {
	listbox.Highlight = false
}

// ResetCursor moves the cursor at the beginning of the view
func (listbox *VirtualListBox) ResetCursor() {
	listbox.SetCursor(0, 0)
}

// Reset resets items of the list box and clears the view
func (listbox *VirtualListBox) Reset() {

	listbox.items = make([]interface{}, 0)
	listbox.pages = []Page{}

	listbox.Clear()
	listbox.ResetCursor()

}

// pageCount returns number of pages in view
func (listbox *VirtualListBox) pageCount() int {
	return len(listbox.pages)
}

// page returns the current page
func (listbox *VirtualListBox) page() Page {
	return listbox.pages[listbox.pageIndex]
}

// currPageIndex returns index of current page
func (listbox *VirtualListBox) currPageIndex() int {

	if listbox.Empty() {
		return 0
	}

	return listbox.pageIndex + 1
}

// SetTitle sets the title of the view and display paging information of the listbox
func (listbox *VirtualListBox) SetTitle(title string) {

	listbox.title = title

	if listbox.pageCount() > 1 {
		listbox.Title = fmt.Sprintf(" %d/%d - %s ", listbox.currPageIndex(), listbox.pageCount(), listbox.title)
	} else {
		listbox.Title = fmt.Sprintf(" %s ", listbox.title)
	}

}

// SetItems sets items of the list box and redraw the view
func (listbox *VirtualListBox) SetItems(items []interface{}) error {

	listbox.items = items

	listbox.RecalcPages()
	return listbox.Draw()

}

// AddItem appends a given item to the existing list box
func (listbox *VirtualListBox) AddItem(item interface{}) error {

	listbox.items = append(listbox.items, item)

	listbox.RecalcPages()
	return listbox.Draw()

}

// cursorY returns the current Y-position of the cursor
func (listbox *VirtualListBox) cursorY() int {

	_, y := listbox.Cursor()
	return y

}

// UpdateCurrentItem sets the current item's new value
func (listbox *VirtualListBox) UpdateCurrentItem(item string) {

	page := listbox.page()
	data := listbox.items[page.offset : page.offset+page.limit]

	data[listbox.cursorY()] = item

}

// MoveDown moves the cursor to the line below on the next page
func (listbox *VirtualListBox) MoveDown() error {

	if listbox.Empty() {
		return nil
	}

	y := listbox.cursorY() + 1

	if listbox.isBottomOfPage() {
		y = 0

		if listbox.hasMultiplePages() {
			listbox.drawPage(listbox.nextPageIndex())
		}
	}

	return listbox.SetCursor(0, y)
}

// MoveUp moves the cursor to the line above on the previous page
func (listbox *VirtualListBox) MoveUp() error {

	if listbox.Empty() {
		return nil
	}

	y := listbox.cursorY() - 1

	if listbox.isTopOfPage() {
		y = listbox.pages[listbox.prevPageIndex()].limit - 1

		if listbox.hasMultiplePages() {
			listbox.drawPage(listbox.prevPageIndex())
		}
	}

	return listbox.SetCursor(0, y)
}

// MovePageDown displays the next page
func (listbox *VirtualListBox) MovePageDown() error {

	if listbox.Empty() {
		return nil
	}

	listbox.drawPage(listbox.nextPageIndex())

	return listbox.SetCursor(0, 0)
}

// MovePageUp displays the previous page
func (listbox *VirtualListBox) MovePageUp() error {

	if listbox.Empty() {
		return nil
	}

	listbox.drawPage(listbox.prevPageIndex())

	return listbox.SetCursor(0, 0)
}

// hasMultiplePages determines whether there is more than one page to be displayed
func (listbox *VirtualListBox) hasMultiplePages() bool {
	return listbox.pageCount() > 1
}

// nextPageIndex return the index of the next page to be displayed
func (listbox *VirtualListBox) nextPageIndex() int {
	return (listbox.pageIndex + 1) % listbox.pageCount()
}

// prevPageIndex returns the index of the previous page to be displayed
func (listbox *VirtualListBox) prevPageIndex() int {

	index := (listbox.pageIndex - 1) % listbox.pageCount()

	if listbox.pageIndex == 0 {
		index = listbox.pageCount() - 1
	}

	return index
}

// isBottomOfPage determines whether the cursor is at the bottom of the current page
func (listbox *VirtualListBox) isBottomOfPage() bool {
	return listbox.cursorY() == (listbox.page().limit - 1)
}

// isTopOfPage determines whether the cursor is at the top of the current page
func (listbox *VirtualListBox) isTopOfPage() bool {
	return listbox.cursorY() == 0
}

// Width returnes the current width of the view
func (listbox *VirtualListBox) Width() int {

	w, _ := listbox.Size()
	return w - 1

}

// Height returns the current height of the view
func (listbox *VirtualListBox) Height() int {

	_, h := listbox.Size()
	return h - 1

}

// RecalcPages recalculates the pages data based on the current length of the list box
// and the current height of the view
func (listbox *VirtualListBox) RecalcPages() {

	listbox.pages = []Page{}

	for offset := 0; offset < listbox.Count(); offset += listbox.Height() {

		limit := listbox.Height()

		if offset+limit > listbox.Count() {
			limit = listbox.Count() % listbox.Height()
		}

		listbox.pages = append(listbox.pages, Page{offset, limit})

	}
}

// Draw calculates the pages and draws the first one
func (listbox *VirtualListBox) Draw() error {

	if listbox.Empty() {
		listbox.Clear()
		return nil
	}

	return listbox.drawPage(0)
}

// DrawCurrentPage calculates the pages and draws the current page
func (listbox *VirtualListBox) DrawCurrentPage() error {

	if listbox.Empty() {
		listbox.Clear()
		return nil
	}

	return listbox.drawPage(listbox.pageIndex)
}

// ItemIndex returns index of the selected item
func (listbox *VirtualListBox) ItemIndex() int {

	if len(listbox.items) == 0 {
		return -1
	}

	page := listbox.pages[listbox.pageIndex]
	return page.offset + listbox.cursorY()
}

// Item returns item with index
func (listbox *VirtualListBox) Item(index int) interface{} {
	return listbox.items[index]
}

// drawPage resets the pageIndex and displays the requested page
func (listbox *VirtualListBox) drawPage(index int) error {

	listbox.Clear()
	listbox.pageIndex = index

	page := listbox.pages[listbox.pageIndex]

	for i := page.offset; i < page.offset+page.limit; i++ {
		if _, err := fmt.Fprintln(listbox.View, listbox.itemCaption(i)); err != nil {
			return err
		}
	}

	listbox.SetTitle(listbox.title)
	listbox.SetCursor(0, 0)

	return nil
}

// itemCaption returns caption of the list box item with index index
func (listbox *VirtualListBox) itemCaption(index int) string {

	item := listbox.items[index]
	text := listbox.doGetText(item)

	if listbox.ordered {
		return fmt.Sprintf(" %2d. %s", index+1, strPadRight(text, 0x0020 /* space */, listbox.Width()-3))
	}

	return fmt.Sprintf(" %s", strPadRight(text, 0x20 /* space */, listbox.Width()-3))
}

func (listbox *VirtualListBox) doGetText(item interface{}) string {

	if listbox.OnGetText != nil {
		return listbox.OnGetText(listbox.View, item)
	}

	return fmt.Sprintf("%v", item)
}

func strPadRight(str string, r rune, count int) string {

	s := str

	if len(s) < count {
		for i := utf8.RuneCountInString(s); i <= count; i++ {
			s += string(r)
		}
	}

	return s
}
