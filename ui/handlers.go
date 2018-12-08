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
	"time"

	"github.com/jroimartin/gocui"

	pl "../playlists"
)

func quit(ui *gocui.Gui, view *gocui.View) error {
	return gocui.ErrQuit
}

func listDown(ui *gocui.Gui, view *gocui.View) error {

	switch view.Name() {
	case captionOfGroupsView:

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

	case captionOfChannelsView:

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

	case captionOfGuideView:
		if err := guide.MoveDown(); err != nil {
			return err
		}
	}

	return nil
}

func listUp(ui *gocui.Gui, view *gocui.View) error {

	switch view.Name() {
	case captionOfGroupsView:

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

	case captionOfChannelsView:

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

	case captionOfGuideView:
		if err := guide.MoveUp(); err != nil {
			return err
		}
	}

	return nil
}

func listPageUp(ui *gocui.Gui, view *gocui.View) error {

	switch view.Name() {
	case captionOfGroupsView:

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

	case captionOfChannelsView:

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

	case captionOfGuideView:
		if err := guide.MovePageUp(); err != nil {
			return err
		}
	}

	return nil
}

func listPageDown(ui *gocui.Gui, view *gocui.View) error {

	switch view.Name() {
	case captionOfGroupsView:

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

	case captionOfChannelsView:

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

	case captionOfGuideView:
		if err := guide.MovePageDown(); err != nil {
			return err
		}
	}

	return nil
}

func switchView(ui *gocui.Gui, view *gocui.View) error {
	switch view.Name() {
	case captionOfGroupsView:

		channels.SetFocus(ui)
		groups.UnFocus()

	case captionOfChannelsView:

		guide.SetFocus(ui)
		channels.UnFocus()

	case captionOfGuideView:

		groups.SetFocus(ui)
		guide.UnFocus()

	}

	return nil
}

func loadGroups(p *pl.Playlist) error {

	groups.SetTitle(captionOfGroupsView)

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

	channels.SetTitle(captionOfChannelsView)

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

	guide.SetTitle(captionOfGuideView)

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
