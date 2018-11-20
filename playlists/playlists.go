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
	db                     *sql.DB
	tx                     *sql.Tx
	stmtInsertPlaylistItem *sql.Stmt
}

// GuideItem contains info about tv programme
type GuideItem struct {
}

// Guide content
type Guide struct {
	db                           *sql.DB
	tx                           *sql.Tx
	stmtAppendGuideChannel       *sql.Stmt
	stmtUpdateGuideChannelID     *sql.Stmt
	stmtAppendChannelDisplayName *sql.Stmt
	stmtAppendChannelURL         *sql.Stmt
	stmtAppendGuideProgramme     *sql.Stmt
	stmtUpdateGuideProgrammePID  *sql.Stmt
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

// Read reads content of the tv guide
func (g *Guide) Read(data []byte, parser *XMLTVParser) (err error) {

	onHead := parser.OnHead
	onChannel := parser.OnChannel
	onProgramme := parser.OnProgramme

	defer func() {
		parser.OnHead = onHead
		parser.OnChannel = onChannel
		parser.OnProgramme = onProgramme
	}()

	tx, err := g.db.Begin()

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

	g.tx = tx

	g.stmtAppendGuideChannel, err = tx.Prepare(cmdAppendGuideChannel)

	if err != nil {
		return
	}

	g.stmtUpdateGuideChannelID, err = tx.Prepare(cmdUpdateGuideChannelID)

	if err != nil {
		return
	}

	g.stmtAppendChannelDisplayName, err = tx.Prepare(cmdAppendChannelDisplayName)

	if err != nil {
		return
	}

	g.stmtAppendChannelURL, err = tx.Prepare(cmdAppendChannelURL)

	if err != nil {
		return
	}

	g.stmtAppendGuideProgramme, err = tx.Prepare(cmdAppendGuideProgramme)

	if err != nil {
		return
	}

	g.stmtUpdateGuideProgrammePID, err = tx.Prepare(cmdUpdateGuideProgrammePID)

	if err != nil {
		return
	}

	parser.OnChannel = func(ch *XMLTVChannel) error {
		return g.appendChannel(ch)
	}

	parser.OnProgramme = func(p *XMLTVProgramme) error {
		return g.appendProgramme(p)
	}

	err = parser.Parse(data)

	return
}

func (g *Guide) appendChannel(c *XMLTVChannel) (err error) {

	if c == nil {
		return errors.New("Guide.AppendChannel: cannot append an empty channel")
	}

	var cid int64

	cs := g.stmtAppendGuideChannel
	us := g.stmtUpdateGuideChannelID

	res, err := cs.Exec(&c.ID)

	if err != nil {
		return
	}

	cid, err = res.LastInsertId()

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

	if err := g.appendChannelDisplayNames(cid, dn); err != nil {
		return err
	}

	urls := make([]*XMLTVChannelURL, len(c.URL))

	for index, url := range c.URL {
		urls[index] = &url
	}

	if err := g.appendChannelURL(cid, urls); err != nil {
		return err
	}

	return nil
}

func (g *Guide) appendChannelDisplayNames(cid int64, d []*XMLTVChannelDisplayName) (err error) {

	for _, dn := range d {
		if _, err := g.stmtAppendChannelDisplayName.Exec(&cid, &dn.Lang, &dn.Value); err != nil {
			return err
		}
	}

	return
}

func (g *Guide) appendChannelURL(cid int64, urls []*XMLTVChannelURL) (err error) {

	for _, url := range urls {
		if _, err = g.stmtAppendChannelURL.Exec(&cid, &url.Value); err != nil {
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

func (g *Guide) appendProgramme(p *XMLTVProgramme) (err error) {

	_, err = g.appendProgrammeRecord(p, g.tx)

	if err != nil {
		return
	}

	return
}

func (g *Guide) appendProgrammeRecord(p *XMLTVProgramme, tx *sql.Tx) (int64, error) {

	var pid int64 = -1

	cs := g.stmtAppendGuideProgramme
	us := g.stmtUpdateGuideProgrammePID

	res, err := cs.Exec(&p.Channel, &p.Start, &p.Stop, &p.PDCStart, &p.VPSStart,
		&p.ShowView, &p.VideoPlus, &p.ClumpIdx)

	if err != nil {
		return pid, err
	}

	pid, err = res.LastInsertId()

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
