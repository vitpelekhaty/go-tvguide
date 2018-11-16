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
	"log"

	// append sqlite3 support for database/sql
	_ "github.com/mattn/go-sqlite3"
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
	db *sql.DB
}

// GuideItem contains info about tv programme
type GuideItem struct {
}

// Guide content
type Guide struct {
	db *sql.DB
}

var (
	db *sql.DB
	pl *Playlist
	g  *Guide
)

func init() {

	db, err := sql.Open("sqlite3", ":memory:")

	if err != nil {
		log.Fatal(err)
	}

	err = initDatabaseStructure(db)

	if err != nil {
		log.Fatal(err)
	}

	pl = &Playlist{db: db}
	g = &Guide{db: db}
}

// CurrentPlaylist returns playlist object
func CurrentPlaylist() *Playlist {
	return pl
}

// CurrentGuide returns guide object
func CurrentGuide() *Guide {
	return g
}

const sqlCreatePlaylistTable = `CREATE TABLE playlist(
    id text,
    channels_group text,
    channel text,
    source text	
)
`
const sqlCreateChannelsTable = `CREATE TABLE channels(
	channel_id text,
	display_name_lang text,
	display_name text
)
`

const sqlCreateChannelURLTable = `CREATE TABLE channels_urls(
	channel_id text,
	url text
)
`

func initDatabaseStructure(db *sql.DB) (err error) {

	if err = createTablePlaylist(db); err != nil {
		return
	}

	if err = createTableChannels(db); err != nil {
		return
	}

	if err = createTableChannelsURL(db); err != nil {
		return
	}

	return
}

func createTable(cmd string, db *sql.DB) error {

	stmt, err := db.Prepare(cmd)

	if err != nil {
		return err
	}

	_, err = stmt.Exec()

	if err != nil {
		return err
	}

	return nil
}

func createTablePlaylist(db *sql.DB) error {
	return createTable(sqlCreatePlaylistTable, db)
}

func createTableChannels(db *sql.DB) error {
	return createTable(sqlCreateChannelsTable, db)
}

func createTableChannelsURL(db *sql.DB) error {
	return createTable(sqlCreateChannelURLTable, db)
}

const sqlSelectGroups = `SELECT pl.channels_group FROM playlist AS pl
GROUP BY pl.channels_group
ORDER BY rowid
`

// Groups returns existing groups in the playlist
func (p *Playlist) Groups() []string {

	var group string

	g := make([]string, 0)

	stmt, err := p.db.Prepare(sqlSelectGroups)

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

const sqlSelectGroupCount = `SELECT COUNT(*) AS cc FROM (
	SELECT pl.channels_group FROM playlist AS pl
	GROUP BY pl.channels_group
) AS items
`

// GroupCount returns number of groups in the playlist
func (p *Playlist) GroupCount() int {

	var count int

	stmt, err := p.db.Prepare(sqlSelectGroupCount)

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

const sqlSelectChannels = `SELECT pl.id, pl.channels_group, pl.channel, pl.source
FROM playlist AS pl 
WHERE pl.channels_group = ?
ORDER BY rowid
`

// Channels returns names of channels for the specified group
func (p *Playlist) Channels(group string) []*PlaylistItem {

	var (
		id       string
		url      string
		channel  string
		gchannel string
	)

	items := make([]*PlaylistItem, 0)

	stmt, err := p.db.Prepare(sqlSelectChannels)

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

const sqlInsertPlaylistItem = `INSERT INTO playlist (id, channels_group, channel, source)
	VALUES(?, ?, ?, ?)
`

// AppendItem appends the item of the playlist into the collection
func (p *Playlist) AppendItem(item *PlaylistItem) (err error) {

	if item == nil {
		return errors.New("Playlist.AppendItem: cannot append an empty item")
	}

	tx, err := p.db.Begin()

	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}

		tx.Commit()
	}()

	if err != nil {
		return
	}

	stmt, err := tx.Prepare(sqlInsertPlaylistItem)

	if err != nil {
		return
	}

	_, err = stmt.Exec(&item.ID, &item.GroupTitle, &item.Name, &item.URL)

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

// AppendChannel appends info about channel into collection
func (g *Guide) AppendChannel(c *XMLTVChannel) error {

	if c == nil {
		return errors.New("Guide.AppendChannel: cannot append an empty channel")
	}

	cid := c.ID

	dn := make([]*XMLTVChannelDisplayName, len(c.DisplayName))

	for index, d := range c.DisplayName {
		dn[index] = &d
	}

	if err := g.appendChannelDisplayNames(cid, dn, g.db); err != nil {
		return err
	}

	urls := make([]*XMLTVChannelURL, len(c.URL))

	for index, url := range c.URL {
		urls[index] = &url
	}

	if err := g.appendChannelURL(cid, urls, g.db); err != nil {
		return err
	}

	return nil
}

const sqlAppendChannelDisplayName = `INSERT INTO channels(channel_id, display_name_lang, display_name) VALUES(?, ?, ?)`

func (g *Guide) appendChannelDisplayNames(cid string, d []*XMLTVChannelDisplayName, db *sql.DB) (err error) {

	tx, err := db.Begin()

	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}

		tx.Commit()
	}()

	if err != nil {
		return
	}

	stmt, err := tx.Prepare(sqlAppendChannelDisplayName)

	if err != nil {
		return
	}

	for _, dn := range d {
		if _, err := stmt.Exec(&cid, &dn.Lang, &dn.Value); err != nil {
			return err
		}
	}

	return
}

const sqlAppendChannelURL = `INSERT INTO channels_urls(channel_id, url) VALUES(?, ?)`

func (g *Guide) appendChannelURL(cid string, urls []*XMLTVChannelURL, db *sql.DB) (err error) {

	tx, err := db.Begin()

	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}

		tx.Commit()
	}()

	if err != nil {
		return
	}

	stmt, err := tx.Prepare(sqlAppendChannelURL)

	if err != nil {
		return
	}

	for _, url := range urls {
		if _, err = stmt.Exec(&cid, &url.Value); err != nil {
			return
		}
	}

	return
}

// Clear clean tv guide collection
func (g *Guide) Clear() (err error) {

	tables := [2]string{"channels", "channels_urls"}

	tx, err := g.db.Begin()

	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}

		tx.Commit()
	}()

	if err != nil {
		return
	}

	stmt, err := tx.Prepare(`DELETE FROM ?`)

	if err != nil {
		return
	}

	for _, table := range tables {
		if _, err = stmt.Exec(&table); err != nil {
			return
		}
	}

	return
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
