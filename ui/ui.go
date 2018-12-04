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
	"fmt"
	"strings"

	"github.com/jroimartin/gocui"

	pl "../playlists"
)

const (
	captionOfGroupsView   = "Groups"
	captionOfChannelsView = "Channels"
	captionOfGuideView    = "Guide"
	captionUndefined      = "<undefined>"
)

var (
	groups   *VirtualListBox
	channels *VirtualListBox
	guide    *VirtualListBox
	playlist *pl.Playlist
	tvg      *pl.Guide
	lang     string
)

// NewPlaylistViewer returns the iptv playlist viewer
func NewPlaylistViewer(p *pl.Playlist, g *pl.Guide) (*gocui.Gui, error) {

	gui, err := gocui.NewGui(gocui.OutputNormal)

	if err != nil {
		return gui, err
	}

	playlist = p
	tvg = g

	lang = tvg.DefaultProgrammeLanguage()

	group, err := playlist.Group(0)

	if err != nil {
		return gui, err
	}

	gui.SetManagerFunc(layout)
	setKeyBindings(gui)

	gw, h := sizeOfGroupsView(gui)
	view, err := gui.SetView(captionOfGroupsView, 0, 0, gw, h-1)

	if err != nil && err != gocui.ErrUnknownView {
		return gui, err
	}

	groups = CreateVirtualListBox(view, true)
	groups.OnGetText = getGroupText

	groups.SetFocus(gui)

	gui.Update(func(g *gocui.Gui) error {
		if err := loadGroups(playlist); err != nil {
			return err
		}

		return nil
	})

	cw, h := sizeOfChannelsView(gui)
	view, err = gui.SetView(captionOfChannelsView, gw+1, 0, gw+cw, h-1)

	if err != nil && err != gocui.ErrUnknownView {
		return gui, err
	}

	channels = CreateVirtualListBox(view, true)
	channels.OnGetText = getChannelText

	gui.Update(func(g *gocui.Gui) error {

		if err := loadChannels(playlist, group); err != nil {
			return err
		}
		return nil
	})

	w, h := gui.Size()
	view, err = gui.SetView(captionOfGuideView, gw+cw+1, 0, w-1, h-1)

	if err != nil && err != gocui.ErrUnknownView {
		return gui, err
	}

	guide = CreateVirtualListBox(view, true)
	guide.SetTitle(captionOfGuideView)

	return gui, nil
}

func getGroupText(view *gocui.View, item interface{}) string {

	if text, ok := item.(string); ok {
		if strings.Trim(text, " ") == "" {
			return captionUndefined
		}

		return text
	}

	return fmt.Sprintf("%v", item)
}

func getChannelText(view *gocui.View, item interface{}) string {

	if pitem, ok := item.(*pl.PlaylistItem); ok {
		text := pitem.Name
		return text
	}

	return fmt.Sprintf("%v", item)
}

func layout(ui *gocui.Gui) error {

	gw, h := sizeOfGroupsView(ui)
	_, err := ui.SetView(captionOfGroupsView, 0, 0, gw, h-1)

	if err != nil {
		return err
	}

	cw, h := sizeOfChannelsView(ui)
	_, err = ui.SetView(captionOfChannelsView, gw+1, 0, gw+cw, h-1)

	if err != nil {
		return err
	}

	w, h := ui.Size()
	_, err = ui.SetView(captionOfGuideView, gw+cw+1, 0, w-1, h-1)

	if err != nil {
		return err
	}

	return nil
}

func sizeOfGroupsView(ui *gocui.Gui) (int, int) {

	w, h := ui.Size()
	return w / 5, h
}

func sizeOfChannelsView(ui *gocui.Gui) (int, int) {

	w, h := ui.Size()
	return w / 5, h
}
