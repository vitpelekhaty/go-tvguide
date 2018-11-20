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

package playlists

import (
	"bufio"
	"errors"
	"strings"
	"unicode"
)

// M3UPlaylistParser - parser for m3u playlist format
type M3UPlaylistParser struct {
	guide string
	items []*PlaylistItem
}

// Parse parses the data of a playlist
func (parser *M3UPlaylistParser) Parse(data []byte) error {

	fitem := func(item *PlaylistItem) error {

		if item != nil {
			parser.items = append(parser.items, item)
		}

		return nil
	}

	return parser.AsyncParse(data, fitem)
}

// AsyncParse parses the data of a playlist
func (parser *M3UPlaylistParser) AsyncParse(data []byte, onItem OnPlaylistItemEvent) error {

	var (
		guide  string
		id     string
		group  string
		name   string
		source string
	)

	parser.items = make([]*PlaylistItem, 0)
	parser.guide = ""

	if len(data) == 0 {
		return errors.New("M3UPlaylistParser: the playlist is empty")
	}

	content := string(data)
	sreader := strings.NewReader(content)

	scanner := bufio.NewScanner(sreader)

	for scanner.Scan() {

		line := printable(scanner.Text())

		if len(line) > 0 {

			if strings.HasPrefix(line, "#EXTM3U") {

				guide = strings.Trim(getValue("url-tvg", line), `"`)
				parser.guide = guide

			} else {

				if strings.HasPrefix(line, "#EXTINF") {

					id = strings.Trim(getValue("tvg-name", line), `"`)
					group = strings.Trim(getValue("group-title", line), `"`)
					name = strings.Trim(getChannelName(line), `"`)

				} else {
					source = line

					item := &PlaylistItem{Name: name, GroupTitle: group, URL: source, ID: id}

					if onItem != nil {
						if err := onItem(item); err != nil {
							return err
						}
					}
				}
			}

		}
	}

	return nil
}

// Guide returns the url of the tv guide
func (parser *M3UPlaylistParser) Guide() string {
	return parser.guide
}

// Items returns items of the playlist
func (parser *M3UPlaylistParser) Items() []*PlaylistItem {
	return parser.items
}

func getValue(key, line string) string {

	if len(strings.TrimSpace(key)) == 0 {
		return ""
	}

	ignore := false

	f := func(c rune) bool {

		fbreak := (unicode.IsSpace(c) || (c == 0x002C /* comma in unicode */)) && !ignore

		if c == 0x0022 /* quote in unicode */ {
			ignore = !ignore
			fbreak = !ignore
		}

		return fbreak
	}

	fields := strings.FieldsFunc(line, f)

	for _, field := range fields {

		if strings.HasPrefix(field, key) {
			return field[len(key)+1:]
		}
	}

	return ""
}

func getChannelName(line string) string {

	if len(line) == 0 {
		return ""
	}

	fields := strings.FieldsFunc(line, func(c rune) bool { return c == 0x002C })

	return fields[len(fields)-1]
}

func isM3U(data []byte) bool {

	if len(data) == 0 {
		return false
	}

	var line string

	s := string(data)

	reader := strings.NewReader(s)
	scanner := bufio.NewScanner(reader)

	for scanner.Scan() {

		line = printable(scanner.Text())

		if len(line) > 0 {
			return strings.HasPrefix(line, "#EXTM3U")
		}
	}

	return false
}

func printable(s string) string {

	var ps string

	for _, r := range s {
		if unicode.IsPrint(r) {
			ps += string(r)
		}
	}

	return ps
}
