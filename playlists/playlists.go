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

var dbname = getPlaylistDatabaseName()

func init() {

	db, err := sql.Open("sqlite3", dbname)

	if err != nil {
		log.Fatal(err)
	}

	err = createDatabaseStructure(db)

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

	stmt, err := tx.Prepare(cmdInsertPlaylistItem)

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
func (g *Guide) AppendChannel(c *XMLTVChannel) (err error) {

	if c == nil {
		return errors.New("Guide.AppendChannel: cannot append an empty channel")
	}

	var cid int64

	tx, err := g.db.Begin()

	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}

		tx.Commit()
	}()

	cs, err := tx.Prepare(cmdAppendGuideChannel)

	if err != nil {
		return
	}

	res, err := cs.Exec(&c.ID)

	if err != nil {
		return
	}

	cid, err = res.LastInsertId()

	if err != nil {
		return
	}

	us, err := tx.Prepare(cmdUpdateGuideChannelID)

	if err != nil {
		return
	}

	_, err = us.Exec(&cid, &cid)

	if err != nil {
		return
	}

	dn := make([]*XMLTVChannelDisplayName, len(c.DisplayName))

	for index, d := range c.DisplayName {
		dn[index] = &d
	}

	if err := g.appendChannelDisplayNames(cid, dn, tx); err != nil {
		return err
	}

	urls := make([]*XMLTVChannelURL, len(c.URL))

	for index, url := range c.URL {
		urls[index] = &url
	}

	if err := g.appendChannelURL(cid, urls, tx); err != nil {
		return err
	}

	return nil
}

func (g *Guide) appendChannelDisplayNames(cid int64, d []*XMLTVChannelDisplayName, tx *sql.Tx) (err error) {

	stmt, err := tx.Prepare(cmdAppendChannelDisplayName)

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

func (g *Guide) appendChannelURL(cid int64, urls []*XMLTVChannelURL, tx *sql.Tx) (err error) {

	stmt, err := tx.Prepare(cmdAppendChannelURL)

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

	commands := [2]string{cmdDeleteFromChannels, cmdDeleteFromChannelsURL}

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

	for _, command := range commands {
		stmt, err := tx.Prepare(command)

		if err != nil {
			return err
		}

		if _, err = stmt.Exec(&command); err != nil {
			return err
		}
	}

	return
}

// AppendProgramme appends info about programme into collection
func (g *Guide) AppendProgramme(p *XMLTVProgramme) (err error) {

	if p == nil {
		return errors.New("Guide.AppendProgramme: cannot append an empty programme")
	}

	tx, err := g.db.Begin()

	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}

		tx.Commit()
	}()

	_, err = g.appendProgramme(p, tx)

	if err != nil {
		return
	}

	return
}

func (g *Guide) appendProgramme(p *XMLTVProgramme, tx *sql.Tx) (int64, error) {

	var pid int64 = -1

	cs, err := tx.Prepare(cmdAppendGuideProgramme)

	if err != nil {
		return pid, err
	}

	res, err := cs.Exec(&p.Channel, &p.Start, &p.Stop, &p.PDCStart, &p.VPSStart,
		&p.ShowView, &p.VideoPlus, &p.ClumpIdx)

	if err != nil {
		return pid, err
	}

	pid, err = res.LastInsertId()

	if err != nil {
		return pid, err
	}

	us, err := tx.Prepare(cmdUpdateGuideProgrammePID)

	if err != nil {
		return pid, err
	}

	_, err = us.Exec(&pid, &pid)

	return pid, err
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
