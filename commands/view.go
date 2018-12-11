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

package commands

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/jroimartin/gocui"
	"github.com/spf13/cobra"

	pl "../playlists"
	ui "../ui"
)

var cmdView = &cobra.Command{
	Use:   "view",
	Short: "Viewing TV guide",
	Long:  "Viewing TV guide for the specified playlist",

	RunE: func(cmd *cobra.Command, args []string) error {

		path := PlaylistPath
		loader := pl.Loader(path)

		data, err := loadPlaylistOrGuide(loader, path)

		if err != nil {
			return err
		}

		parser := pl.PlaylistParser(data)

		if parser == nil {
			return errors.New("Playlist view: unknown playlist format")
		}

		playlist := pl.CurrentPlaylist()

		err = playlist.Read(data, parser)

		if err != nil {
			return err
		}

		gpath := parser.Guide()
		gloader := pl.Loader(gpath)

		data, err = loadPlaylistOrGuide(gloader, gpath)

		if err != nil {
			return err
		}

		guide := pl.CurrentGuide()
		gparser := &pl.XMLTVParser{}

		err = func(g *pl.Guide, p *pl.XMLTVParser, d []byte) error {

			st := time.Now()

			fmt.Println("TV guide reading. Please, wait...")

			defer func(t time.Time) {

				et := time.Now()
				d := et.Sub(t)

				fmt.Printf("TV Guide reading completed in %.3fs\n", d.Seconds())
			}(st)

			return guide.Read(data, gparser)

		}(guide, gparser, data)

		if err != nil {
			return err
		}

		gui, err := ui.NewPlaylistViewer(playlist, guide)

		if err != nil {
			return err
		}

		defer gui.Close()

		if err := gui.MainLoop(); err != nil && err != gocui.ErrQuit {
			return err
		}

		return nil
	},
}

func loadPlaylistOrGuide(loader pl.ILoader, path string) ([]byte, error) {

	data := make([]byte, 0)

	switch loader.(type) {
	case *pl.FileLoader:
		if floader, ok := loader.(*pl.FileLoader); ok {
			return loadFromFile(floader, path)
		}

		return data, errors.New("Playlist or guide loading: something wrong")

	case *pl.HTTPLoader:
		if nloader, ok := loader.(*pl.HTTPLoader); ok {
			return loadFromURL(nloader, path)
		}

		return data, errors.New("Playlist or guide loading: something wrong")

	case nil:
		return data, errors.New("Playlist or guide loading: invalid path")
	}

	return data, nil
}

func loadFromFile(loader *pl.FileLoader, path string) ([]byte, error) {

	fmt.Printf("Loading file %s\t...\n", path)
	return loader.Load(path)
}

func loadFromURL(loader *pl.HTTPLoader, url string) ([]byte, error) {

	comment := "Downloading " + url

	fprogress := func(complete uint64) {
		fmt.Printf("\r%s ... %s", comment, strings.Repeat(" ", 35))
		fmt.Printf("\r%s ... %s", comment, humanize.Bytes(complete))
	}

	fdone := func() {
		fmt.Print("\n")
	}

	loader.OnProgress = fprogress
	loader.OnDone = fdone

	return loader.Load(url)
}
