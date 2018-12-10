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
	"errors"
	"fmt"
	"time"

	"github.com/jroimartin/gocui"
	"github.com/logrusorgru/aurora"

	pl "../playlists"
)

func quit(ui *gocui.Gui, view *gocui.View) error {
	return gocui.ErrQuit
}

func listDown(ui *gocui.Gui, view *gocui.View) error {

	switch view.Name() {
	case viewGroups:

		if err := groups.MoveDown(); err != nil {
			return err
		}

		index := groups.ItemIndex()
		group := groups.Item(index)

		if gtext, ok := group.(string); ok {
			ui.Update(func(g *gocui.Gui) error {

				if err := loadChannels(playlist, gtext); err != nil {
					return err
				}

				index := channels.ItemIndex()
				ch := channels.Item(index)

				if pi, ok := ch.(*pl.PlaylistItem); ok {

					t := CurrentTime()

					if err := loadChannelGuide(tvg, pi.ID, lang, t); err != nil {
						return err
					}
				}

				return nil
			})
		}

	case viewChannels:

		if err := channels.MoveDown(); err != nil {
			return err
		}

		index := channels.ItemIndex()
		ch := channels.Item(index)

		if pi, ok := ch.(*pl.PlaylistItem); ok {

			t := CurrentTime()

			if err := loadChannelGuide(tvg, pi.ID, lang, t); err != nil {
				return err
			}
		}

	case viewGuide:
		if err := guide.MoveDown(); err != nil {
			return err
		}
	}

	return nil
}

func listUp(ui *gocui.Gui, view *gocui.View) error {

	switch view.Name() {
	case viewGroups:

		if err := groups.MoveUp(); err != nil {
			return err
		}

		index := groups.ItemIndex()
		group := groups.Item(index)

		if gtext, ok := group.(string); ok {
			ui.Update(func(g *gocui.Gui) error {

				if err := loadChannels(playlist, gtext); err != nil {
					return err
				}

				index := channels.ItemIndex()
				ch := channels.Item(index)

				if pi, ok := ch.(*pl.PlaylistItem); ok {

					t := CurrentTime()

					if err := loadChannelGuide(tvg, pi.ID, lang, t); err != nil {
						return err
					}
				}

				return nil
			})
		}

	case viewChannels:

		if err := channels.MoveUp(); err != nil {
			return err
		}

		index := channels.ItemIndex()
		ch := channels.Item(index)

		if pi, ok := ch.(*pl.PlaylistItem); ok {

			t := CurrentTime()

			if err := loadChannelGuide(tvg, pi.ID, lang, t); err != nil {
				return err
			}
		}

	case viewGuide:
		if err := guide.MoveUp(); err != nil {
			return err
		}
	}

	return nil
}

func listPageUp(ui *gocui.Gui, view *gocui.View) error {

	switch view.Name() {
	case viewGroups:

		if err := groups.MovePageUp(); err != nil {
			return err
		}

		index := groups.ItemIndex()
		group := groups.Item(index)

		if gtext, ok := group.(string); ok {
			ui.Update(func(g *gocui.Gui) error {

				if err := loadChannels(playlist, gtext); err != nil {
					return err
				}

				index := channels.ItemIndex()
				ch := channels.Item(index)

				if pi, ok := ch.(*pl.PlaylistItem); ok {

					t := CurrentTime()

					if err := loadChannelGuide(tvg, pi.ID, lang, t); err != nil {
						return err
					}
				}

				return nil
			})
		}

	case viewChannels:

		if err := channels.MovePageUp(); err != nil {
			return err
		}

		index := channels.ItemIndex()
		ch := channels.Item(index)

		if pi, ok := ch.(*pl.PlaylistItem); ok {

			t := CurrentTime()

			if err := loadChannelGuide(tvg, pi.ID, lang, t); err != nil {
				return err
			}
		}

	case viewGuide:
		if err := guide.MovePageUp(); err != nil {
			return err
		}
	}

	return nil
}

func listPageDown(ui *gocui.Gui, view *gocui.View) error {

	switch view.Name() {
	case viewGroups:

		if err := groups.MovePageDown(); err != nil {
			return err
		}

		index := groups.ItemIndex()
		group := groups.Item(index)

		if gtext, ok := group.(string); ok {
			ui.Update(func(g *gocui.Gui) error {

				if err := loadChannels(playlist, gtext); err != nil {
					return err
				}

				index := channels.ItemIndex()
				ch := channels.Item(index)

				if pi, ok := ch.(*pl.PlaylistItem); ok {

					t := CurrentTime()

					if err := loadChannelGuide(tvg, pi.ID, lang, t); err != nil {
						return err
					}
				}

				return nil
			})
		}

	case viewChannels:

		if err := channels.MovePageDown(); err != nil {
			return err
		}

		index := channels.ItemIndex()
		ch := channels.Item(index)

		if pi, ok := ch.(*pl.PlaylistItem); ok {

			t := CurrentTime()

			if err := loadChannelGuide(tvg, pi.ID, lang, t); err != nil {
				return err
			}
		}

	case viewGuide:
		if err := guide.MovePageDown(); err != nil {
			return err
		}
	}

	return nil
}

func switchView(ui *gocui.Gui, view *gocui.View) error {

	switch view.Name() {
	case viewGroups:

		channels.SetFocus(ui)
		groups.UnFocus()

	case viewChannels:

		guide.SetFocus(ui)
		channels.UnFocus()

	case viewGuide:

		groups.SetFocus(ui)
		guide.UnFocus()

	}

	return nil
}

func loadGroups(p *pl.Playlist) error {

	groups.SetTitle(titleGroups)

	if p == nil {
		return errors.New("Failed to load playlist")
	}

	g := p.Groups()
	data := make([]interface{}, len(g))

	for i, gr := range g {
		data[i] = gr
	}

	return groups.SetItems(data)
}

func loadChannels(p *pl.Playlist, group string) error {

	channels.SetTitle(titleChannels)

	if p == nil {
		return errors.New("Failed to load playist")
	}

	c := p.Channels(group)
	data := make([]interface{}, len(c))

	for i, ch := range c {
		data[i] = ch
	}

	return channels.SetItems(data)
}

func loadChannelGuide(g *pl.Guide, cid string, lang string, t time.Time) error {

	guide.SetTitle(titleGuide)

	if g == nil {
		return errors.New("Failed to load tv guide")
	}

	gg, err := g.ChannelGuide(cid, lang, t)

	if err != nil {
		return err
	}

	data := make([]interface{}, len(gg))

	for i, p := range gg {
		data[i] = p
	}

	return guide.SetItems(data)
}

func help(ui *gocui.Gui, view *gocui.View) error {

	curview = ui.CurrentView()

	if err := createHelpView(ui, titleHelp); err != nil {
		return err
	}

	return nil
}

func setTopWindowTitle(ui *gocui.Gui, view, title string) error {

	v, err := ui.View(view)

	if err != nil {
		return err
	}

	v.Title = fmt.Sprintf(" %v (Ctrl+Q to close) ", title)

	return nil
}

func createHelpView(ui *gocui.Gui, title string) error {

	w, h := ui.Size()
	v, err := ui.SetView(viewHelp, w/4, h/4, w*3/4, h*3/4)

	if err != nil && err != gocui.ErrUnknownView {
		return err
	}

	v.Wrap = true

	setTopWindowTitle(ui, viewHelp, title)

	fmt.Fprint(v, "  \n")

	fmt.Fprintf(v, " %v: %v", aurora.Bold("Tab"), "Focuses between Groups, Channels or Guide lists\n")
	fmt.Fprintf(v, " %v: %v", aurora.Bold("Ctrl+6"), "Focuses Groups list\n")
	fmt.Fprintf(v, " %v: %v", aurora.Bold("Ctrl+7"), "Focuses Channels list\n")
	fmt.Fprintf(v, " %v: %v", aurora.Bold("Ctrl+8"), "Focuses Guide list\n")
	fmt.Fprintf(v, " %v: %v", aurora.Bold("ArrowUp"), "Moves to the previous list item circularly\n")
	fmt.Fprintf(v, " %v: %v", aurora.Bold("ArrowDn"), "Moves to the next list item circularly\n")
	fmt.Fprintf(v, " %v: %v", aurora.Bold("PgUp"), "Moves to the previous list page circularly\n")
	fmt.Fprintf(v, " %v: %v", aurora.Bold("PgDn"), "Moves to the next list page circularly\n")
	fmt.Fprintf(v, " %v: %v", aurora.Bold("Enter"), "\n")
	fmt.Fprintf(v, " %v: %v", aurora.Bold("Esc"), "Closes any window displayed on top of the main windows\n")
	fmt.Fprintf(v, " %v: %v", aurora.Bold("F1"), "Opens up the Help window\n")
	fmt.Fprintf(v, " %v: %v", aurora.Bold("Ctrl+C"), "Exits the application\n")

	_, err = ui.SetCurrentView(viewHelp)

	return err
}

func destroyTopView(ui *gocui.Gui, view *gocui.View) error {

	if curview != nil {
		if _, err := ui.SetCurrentView(curview.Name()); err != nil {
			return err
		}
	} else {
		if _, err := ui.SetCurrentView(viewGroups); err != nil {
			return err
		}
	}

	switch view.Name() {
	case viewHelp:

		if err := destroyHelpView(ui); err != nil {
			return err
		}

	case viewProgramme:

		if err := destroyProgrammeView(ui); err != nil {
			return err
		}
	}

	return nil
}

func destroyHelpView(ui *gocui.Gui) error {

	ui.Cursor = false
	return ui.DeleteView(viewHelp)
}

func destroyProgrammeView(ui *gocui.Gui) error {

	ui.Cursor = false
	return ui.DeleteView(viewProgramme)
}

func focusGroupsView(ui *gocui.Gui, view *gocui.View) error {

	switch view.Name() {

	case viewChannels:

		groups.SetFocus(ui)
		channels.UnFocus()

	case viewGuide:

		groups.SetFocus(ui)
		guide.UnFocus()

	}

	return nil
}

func focusChannelsView(ui *gocui.Gui, view *gocui.View) error {

	switch view.Name() {

	case viewGroups:

		channels.SetFocus(ui)
		groups.UnFocus()

	case viewGuide:

		channels.SetFocus(ui)
		guide.UnFocus()

	}

	return nil
}

func focusGuideView(ui *gocui.Gui, view *gocui.View) error {

	switch view.Name() {

	case viewChannels:

		guide.SetFocus(ui)
		channels.UnFocus()

	case viewGroups:

		guide.SetFocus(ui)
		groups.UnFocus()

	}
	return nil
}

func onEnter(ui *gocui.Gui, view *gocui.View) error {

	switch view.Name() {
	case viewGuide:

		curview = view

		if err := createProgrammeView(ui, titleProgramme); err != nil {
			return err
		}
	}

	return nil
}

func createProgrammeView(ui *gocui.Gui, title string) error {

	w, h := ui.Size()
	v, err := ui.SetView(viewProgramme, w/6, h/6, w*5/6, h*5/6)

	if err != nil && err != gocui.ErrUnknownView {
		return err
	}

	v.Wrap = true

	setTopWindowTitle(ui, viewProgramme, title)

	_, err = ui.SetCurrentView(viewProgramme)

	return err
}
