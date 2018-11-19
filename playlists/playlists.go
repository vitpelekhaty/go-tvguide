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

const (
	cmdCreateTablePlaylist    = `CREATE TABLE playlist(id TEXT, channels_group TEXT, channel TEXT, source TEXT)`
	cmdCreateIndexPlaylistCID = `CREATE INDEX ix_playlist_channel_id ON playlist(id)`

	cmdCreateTableChannels          = `CREATE TABLE channels(cid INTEGER, channel_id TEXT)`
	cmdCreateIndexChannelsCID       = `CREATE INDEX ix_channels_cid ON channels(cid)`
	cmdCreateIndexChannelsChannelID = `CREATE INDEX ix_channels_channel_id ON channels(channel_id)`

	cmdCreateTableChannelDisplayNames    = `CREATE TABLE channel_display_names(cid INTEGER, lang TEXT, display_name TEXT)`
	cmdCreateIndexChannelDisplayNamesCID = `CREATE INDEX ix_channel_display_names_cid ON channel_display_names(cid)`

	cmdCreateChannelURLTable    = `CREATE TABLE channels_urls(cid INTEGER, url TEXT)`
	cmdCreateIndexChannelURLCID = `CREATE INDEX ix_channels_urls_cid ON channels_urls(cid)`

	cmdCreateTableProgramme = `CREATE TABLE programme (
	pid INTEGER,
	channel_id TEXT,
	start TEXT,
	stop TEXT,
	pdc_start TEXT,
	vps_start TEXT,
	show_view TEXT,
	video_plus TEXT,
	clump_idx TEXT		
	)`

	cmdCreateIndexProgrammePID       = `CREATE INDEX ix_programme_pid ON programme(pid)`
	cmdCreateIndexProgrammeChannelID = `CREATE INDEX ix_programme_channel_id ON programme(channel_id)`

	cmdCreateTableProgrammeTitles    = `CREATE TABLE programme_titles(pid INTEGER, lang TEXT, title TEXT)`
	cmdCreateIndexProgrammeTitlesPID = `CREATE INDEX ix_programme_titles_pid ON programme_titles(pid)`

	cmdCreateTableProgrammeSubTitle    = `CREATE TABLE programme_sub_titles(pid INTEGER, lang TEXT, sub_title TEXT)`
	cmdCreateIndexProgrammeSubTitlePID = `CREATE INDEX ix_programme_sub_titles_pid ON programme_sub_titles(pid)`

	cmdCreateTableProgrammeDesc    = `CREATE TABLE programme_desc(pid INTEGER, lang TEXT, desc TEXT)`
	cmdCreateIndexProgrammeDescPID = `CREATE INDEX ix_programme_desc_pid ON programme_desc(pid)`

	cmdCreateTableProgrammeDates    = `CREATE TABLE programme_dates(pid INTEGER, date TEXT)`
	cmdCreateIndexProgrammeDatesPID = `CREATE INDEX ix_programme_dates_pid ON programme_dates(pid)`

	cmdCreateTableProgrammeCategories    = `CREATE TABLE programme_categories(pid INTEGER, lang TEXT, category TEXT)`
	cmdCreateIndexProgrammeCategoriesPID = `CREATE INDEX ix_programme_categories_pid ON programme_categories(pid)`

	cmdCreateTableProgrammeKeywords    = `CREATE TABLE programme_keywords(pid INTEGER, lang TEXT, keyword TEXT)`
	cmdCreateIndexProgrammeKeywordsPID = `CREATE INDEX ix_programme_keywords_pid ON programme_keywords(pid)`

	cmdCreateTableProgrammeLanguage = `CREATE TABLE programme_languages(pid INTEGER, lang TEXT, language TEXT)`

	cmdCreateTableProgrammeOriginalLanguage = `CREATE TABLE programme_original_languages(pid INTEGER, lang TEXT, language TEXT)`

	cmdCreateTableProgrammeCountries = `CREATE TABLE programme_countries(pid INTEGER, lang TEXT, country TEXT)`

	cmdCreateTableProgrammeDirectors = `CREATE TABLE programme_directors(pid INTEGER, director TEXT)`

	cmdCreateTableProgrammeWriters = `CREATE TABLE programme_writers(pid INTEGER, writer TEXT)`

	cmdCreateTableProgrammeAdapters = `CREATE TABLE programme_adapters(pid INTEGER, adapter TEXT)`

	cmdCreateTableProgrammeProducers = `CREATE TABLE programme_producers(pid INTEGER, producer TEXT)`

	cmdCreateTableProgrammeComposers = `CREATE TABLE programme_composers(pid INTEGER, composer TEXT)`

	cmdCreateTableProgrammeEditors = `CREATE TABLE programme_editors(pid INTEGER, editor TEXT)`

	cmdCreateTableProgrammePresenters = `CREATE TABLE programme_presenter(pid INTEGER, presenter TEXT)`

	cmdCreateTableProgrammeCommentators = `CREATE TABLE programme_commentators(pid INTEGER, commentator TEXT)`

	cmdCreateTableProgrammeGuests = `CREATE TABLE programme_guests(pid INTEGER, guest TEXT)`

	cmdCreateTableProgrammeActors = `CREATE TABLE programme_actors(pid INTEGER, actor TEXT, role TEXT)`

	cmdCreateTableProgrammeLength = `CREATE TABLE programme_length(pid INTEGER, value TEXT, units TEXT)`

	cmdCreateTableProgrammeIcon = `CREATE TABLE programme_icon(pid INTEGER, src TEXT, width TEXT, height TEXT)`

	cmdCreateTableProgrammeEpisodeNum = `CREATE TABLE programme_episode_num(pid INTEGER, system TEXT, episode_num TEXT)`

	cmdCreateTableProgrammeVideo = `CREATE TABLE programme_video(pid INTEGER, present TEXT, colour TEXT, aspect TEXT, quality TEXT)`

	cmdCreateTableProgrammeAudio = `CREATE TABLE programme_audio(pid INTEGER, present TEXT, stereo TEXT)`

	cmdCreateTableProgrammePreviouslyShown = `CREATE TABLE programme_previously_shown(pid INTEGER, start TEXT, channel TEXT)`

	cmdCreateTableProgrammePremiere = `CREATE TABLE programme_premiere(pid INTEGER, lang TEXT, premiere TEXT)`

	cmdCreateTableProgrammeLastChance = `CREATE TABLE programme_last_chance(pid INTEGER, lang TEXT, last_chance TEXT)`
)

func initDatabaseStructure(db *sql.DB) (err error) {

	objects := [45]string{cmdCreateTablePlaylist, cmdCreateIndexPlaylistCID,
		cmdCreateTableChannels, cmdCreateIndexChannelsCID, cmdCreateIndexChannelsChannelID,
		cmdCreateTableChannelDisplayNames, cmdCreateIndexChannelDisplayNamesCID,
		cmdCreateChannelURLTable, cmdCreateIndexChannelURLCID,
		cmdCreateTableProgramme, cmdCreateIndexProgrammePID, cmdCreateIndexProgrammeChannelID,
		cmdCreateTableProgrammeTitles, cmdCreateIndexProgrammeTitlesPID,
		cmdCreateTableProgrammeSubTitle, cmdCreateIndexProgrammeSubTitlePID,
		cmdCreateTableProgrammeDesc, cmdCreateIndexProgrammeDescPID,
		cmdCreateTableProgrammeDates, cmdCreateIndexProgrammeDatesPID,
		cmdCreateTableProgrammeCategories, cmdCreateIndexProgrammeCategoriesPID,
		cmdCreateTableProgrammeKeywords, cmdCreateIndexProgrammeKeywordsPID,
		cmdCreateTableProgrammeLanguage,
		cmdCreateTableProgrammeOriginalLanguage,
		cmdCreateTableProgrammeCountries,
		cmdCreateTableProgrammeDirectors,
		cmdCreateTableProgrammeWriters,
		cmdCreateTableProgrammeAdapters,
		cmdCreateTableProgrammeProducers,
		cmdCreateTableProgrammeComposers,
		cmdCreateTableProgrammeEditors,
		cmdCreateTableProgrammePresenters,
		cmdCreateTableProgrammeCommentators,
		cmdCreateTableProgrammeGuests,
		cmdCreateTableProgrammeActors,
		cmdCreateTableProgrammeLength,
		cmdCreateTableProgrammeIcon,
		cmdCreateTableProgrammeEpisodeNum,
		cmdCreateTableProgrammeVideo,
		cmdCreateTableProgrammeAudio,
		cmdCreateTableProgrammePreviouslyShown,
		cmdCreateTableProgrammePremiere,
		cmdCreateTableProgrammeLastChance}

	for _, dbobj := range objects {
		err = execsql(dbobj, db)

		if err != nil {
			return
		}
	}

	return
}

func execsql(cmd string, db *sql.DB) error {

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

const cmdSelectGroups = `SELECT pl.channels_group FROM playlist AS pl
GROUP BY pl.channels_group
ORDER BY rowid
`

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

const cmdSelectGroupCount = `SELECT COUNT(*) AS cc FROM (
	SELECT pl.channels_group FROM playlist AS pl
	GROUP BY pl.channels_group
) AS items
`

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

const cmdSelectChannels = `SELECT pl.id, pl.channels_group, pl.channel, pl.source
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

const cmdInsertPlaylistItem = `INSERT INTO playlist (id, channels_group, channel, source) VALUES(?, ?, ?, ?)`

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

const (
	cmdAppendGuideChannel   = `INSERT INTO channels(channel_id) VALUES(?)`
	cmdUpdateGuideChannelID = `UPDATE channels SET cid = ? WHERE rowid = ?`
)

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

const cmdAppendChannelDisplayName = `INSERT INTO channel_display_names(cid, lang, display_name) VALUES(?, ?, ?)`

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

const cmdAppendChannelURL = `INSERT INTO channels_urls(cid, url) VALUES(?, ?)`

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
