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
	"database/sql"
	"errors"
	"fmt"
)

// PlaylistItem contains info about tv channel (URL, name, etc)
type PlaylistItem struct {
	Name       string
	GroupTitle string
	URL        string
	ID         string
}

// Playlist content
type Playlist struct {
	pdb
	db                     *sql.DB
	tx                     *sql.Tx
	stmtInsertPlaylistItem *sql.Stmt
}

var (
	pl *Playlist
)

// CurrentPlaylist returns playlist object
func CurrentPlaylist() *Playlist {

	if pl == nil {
		pl = &Playlist{db: db}
	}

	return pl
}

// Read reads content of the playlist
func (p *Playlist) Read(data []byte, parser IPlaylistParser) (err error) {

	tx, err := p.db.Begin()

	if err != nil {
		return
	}

	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}

		tx.Commit()
	}()

	p.tx = tx
	p.stmtInsertPlaylistItem, err = tx.Prepare(cmdInsertPlaylistItem)

	if err != nil {
		return
	}

	callback := func(item *PlaylistItem) error {
		return p.appendItem(item)
	}

	err = parser.AsyncParse(data, callback)

	return
}

// Groups returns existing groups in the playlist
func (p *Playlist) Groups() []string {

	var group string

	g := make([]string, 0)

	stmt, err := p.db.Prepare(cmdSelectGroups)

	if err != nil {
		return g
	}

	defer stmt.Close()

	rows, err := stmt.Query()

	if err != nil {
		return g
	}

	for rows.Next() {
		err := rows.Scan(&group)

		if err == nil {
			g = append(g, group)
		}

	}

	err = rows.Err()

	if err != nil {
		return make([]string, 0)
	}

	return g
}

// GroupCount returns number of groups in the playlist
func (p *Playlist) GroupCount() int {

	var count int

	stmt, err := p.db.Prepare(cmdSelectGroupCount)

	if err != nil {
		return 0
	}

	defer stmt.Close()

	err = stmt.QueryRow().Scan(&count)

	if err != nil {
		return 0
	}

	return count
}

// Channels returns names of channels for the specified group
func (p *Playlist) Channels(group string) []*PlaylistItem {

	var (
		id       string
		url      string
		channel  string
		gchannel string
	)

	items := make([]*PlaylistItem, 0)

	stmt, err := p.db.Prepare(cmdSelectChannels)

	if err != nil {
		return items
	}

	defer stmt.Close()

	rows, err := stmt.Query(group)

	if err != nil {
		return items
	}

	for rows.Next() {
		err := rows.Scan(&id, &gchannel, &channel, &url)

		if err == nil {
			item := &PlaylistItem{ID: id, GroupTitle: gchannel, Name: channel, URL: url}
			items = append(items, item)
		}
	}

	err = rows.Err()

	if err != nil {
		return make([]*PlaylistItem, 0)
	}

	return items
}

func (p *Playlist) appendItem(item *PlaylistItem) (err error) {

	if item == nil {
		return errors.New("Playlist.AppendItem: cannot append an empty item")
	}

	_, err = p.stmtInsertPlaylistItem.Exec(&item.ID, &item.GroupTitle, &item.Name, &item.URL)

	if err != nil {
		return
	}

	return nil
}

// Group returns group name with specified index
func (p *Playlist) Group(index int) (string, error) {

	g := p.Groups()

	if index >= 0 && index < len(g) {
		return g[index], nil
	}

	return "", fmt.Errorf("Index (%d) out of bounds", index)
}

// PlaylistParser return the parser fo appropriate playlist format
func PlaylistParser(data []byte) IPlaylistParser {

	if isM3U(data) {
		return &M3UPlaylistParser{}
	}

	return nil
}

func contains(l []string, s string) bool {

	if len(l) == 0 {
		return false
	}

	for _, item := range l {
		if item == s {
			return true
		}
	}

	return false
}
