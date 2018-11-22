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
	"log"

	// append sqlite3 support for database/sql
	_ "github.com/mattn/go-sqlite3"
)

const (
	cmdCreateTablePlaylist    = `CREATE TABLE playlist(id TEXT, channels_group TEXT, channel TEXT, source TEXT)`
	cmdCreateIndexPlaylistCID = `CREATE INDEX ix_playlist_channel_id ON playlist(id)`

	cmdCreateTableChannels          = `CREATE TABLE channels(cid INTEGER, channel_id TEXT)`
	cmdCreateIndexChannelsCID       = `CREATE INDEX ix_channels_cid ON channels(cid)`
	cmdCreateIndexChannelsChannelID = `CREATE INDEX ix_channels_channel_id ON channels(channel_id)`

	cmdCreateTableChannelDisplayNames    = `CREATE TABLE channel_display_names(cid INTEGER, lang TEXT, display_name TEXT)`
	cmdCreateIndexChannelDisplayNamesCID = `CREATE INDEX ix_channel_display_names_cid ON channel_display_names(cid)`

	cmdCreateChannelURLTable    = `CREATE TABLE channel_urls(cid INTEGER, url TEXT)`
	cmdCreateIndexChannelURLCID = `CREATE INDEX ix_channel_urls_cid ON channel_urls(cid)`

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

	cmdCreateTableProgrammeLanguage    = `CREATE TABLE programme_languages(pid INTEGER, lang TEXT, language TEXT)`
	cmdCreateIndexProgrammeLanguagePID = `CREATE INDEX ix_programme_languages_pid ON programme_languages(pid)`

	cmdCreateTableProgrammeOriginalLanguage    = `CREATE TABLE programme_original_languages(pid INTEGER, lang TEXT, language TEXT)`
	cmdCreateIndexProgrammeOriginalLanguagePID = `CREATE INDEX ix_programme_original_languages_pid ON programme_original_languages(pid)`

	cmdCreateTableProgrammeCountries    = `CREATE TABLE programme_countries(pid INTEGER, lang TEXT, country TEXT)`
	cmdCreateIndexProgrammeCountriesPID = `CREATE INDEX ix_programme_countries_pid ON programme_countries(pid)`

	cmdCreateTableProgrammeDirectors    = `CREATE TABLE programme_directors(pid INTEGER, director TEXT)`
	cmdCreateIndexProgrammeDirectorsPID = `CREATE INDEX ix_programme_directors_pid ON programme_directors(pid)`

	cmdCreateTableProgrammeWriters    = `CREATE TABLE programme_writers(pid INTEGER, writer TEXT)`
	cmdCreateIndexProgrammeWritersPID = `CREATE INDEX ix_programme_writers_pid ON programme_writers(pid)`

	cmdCreateTableProgrammeAdapters    = `CREATE TABLE programme_adapters(pid INTEGER, adapter TEXT)`
	cmdCreateIndexProgrammeAdaptersPID = `CREATE INDEX ix_programme_adapters_pid ON programme_adapters(pid)`

	cmdCreateTableProgrammeProducers    = `CREATE TABLE programme_producers(pid INTEGER, producer TEXT)`
	cmdCreateIndexProgrammeProducersPID = `CREATE INDEX ix_programme_producers_pid ON programme_producers(pid)`

	cmdCreateTableProgrammeComposers    = `CREATE TABLE programme_composers(pid INTEGER, composer TEXT)`
	cmdCreateIndexProgrammeComposersPID = `CREATE INDEX ix_programme_composers_pid ON programme_composers(pid)`

	cmdCreateTableProgrammeEditors    = `CREATE TABLE programme_editors(pid INTEGER, editor TEXT)`
	cmdCreateIndexProgrammeEditorsPID = `CREATE INDEX ix_programme_editors_pid ON programme_editors(pid)`

	cmdCreateTableProgrammePresenters    = `CREATE TABLE programme_presenters(pid INTEGER, presenter TEXT)`
	cmdCreateIndexProgrammePresentersPID = `CREATE INDEX ix_programme_presenters_pid ON programme_presenters(pid)`

	cmdCreateTableProgrammeCommentators    = `CREATE TABLE programme_commentators(pid INTEGER, commentator TEXT)`
	cmdCreateIndexProgrammeCommentatorsPID = `CREATE INDEX ix_programme_commentators_pid ON programme_commentators(pid)`

	cmdCreateTableProgrammeGuests    = `CREATE TABLE programme_guests(pid INTEGER, guest TEXT)`
	cmdCreateIndexProgrammeGuestsPID = `CREATE INDEX ix_programme_guests_pid ON programme_guests(pid)`

	cmdCreateTableProgrammeActors    = `CREATE TABLE programme_actors(pid INTEGER, actor TEXT, role TEXT)`
	cmdCreateIndexProgrammeActorsPID = `CREATE INDEX ix_programme_actors_pid ON programme_actors(pid)`

	cmdCreateTableProgrammeLength    = `CREATE TABLE programme_length(pid INTEGER, value TEXT, units TEXT)`
	cmdCreateIndexProgrammeLengthPID = `CREATE INDEX ix_programme_length_pid ON programme_length(pid)`

	cmdCreateTableProgrammeIcon    = `CREATE TABLE programme_icon(pid INTEGER, src TEXT, width TEXT, height TEXT)`
	cmdCreateIndexProgrammeIconPID = `CREATE INDEX ix_programme_icon_pid ON programme_icon(pid)`

	cmdCreateTableProgrammeEpisodeNum    = `CREATE TABLE programme_episode_num(pid INTEGER, system TEXT, episode_num TEXT)`
	cmdCreateIndexProgrammeEpisodeNumPID = `CREATE INDEX ix_programme_episode_num_pid ON programme_episode_num(pid)`

	cmdCreateTableProgrammeVideo    = `CREATE TABLE programme_video(pid INTEGER, present TEXT, colour TEXT, aspect TEXT, quality TEXT)`
	cmdCreateIndexProgrammeVideoPID = `CREATE INDEX ix_programme_video_pid ON programme_video(pid)`

	cmdCreateTableProgrammeAudio    = `CREATE TABLE programme_audio(pid INTEGER, present TEXT, stereo TEXT)`
	cmdCreateIndexProgrammeAudioPID = `CREATE INDEX ix_programme_audio_pid ON programme_audio(pid)`

	cmdCreateTableProgrammePreviouslyShown    = `CREATE TABLE programme_previously_shown(pid INTEGER, start TEXT, channel TEXT)`
	cmdCreateIndexProgrammePreviouslyShownPID = `CREATE INDEX ix_programme_previously_shown_pid ON programme_previously_shown(pid)`

	cmdCreateTableProgrammePremiere    = `CREATE TABLE programme_premiere(pid INTEGER, lang TEXT, premiere TEXT)`
	cmdCreateIndexProgrammePremierePID = `CREATE INDEX ix_programme_premiere_pid ON programme_premiere(pid)`

	cmdCreateTableProgrammeLastChance    = `CREATE TABLE programme_last_chance(pid INTEGER, lang TEXT, last_chance TEXT)`
	cmdCreateIndexProgrammeLastChancePID = `CREATE INDEX ix_programme_last_chance_pid ON programme_last_chance(pid)`
)

const (
	cmdSelectGroups = `SELECT pl.channels_group FROM playlist AS pl
	GROUP BY pl.channels_group
	ORDER BY rowid
	`

	cmdSelectGroupCount = `SELECT COUNT(*) AS cc FROM (
		SELECT pl.channels_group FROM playlist AS pl
		GROUP BY pl.channels_group
	) AS items
	`

	cmdSelectChannels = `SELECT pl.id, pl.channels_group, pl.channel, pl.source
	FROM playlist AS pl 
	WHERE pl.channels_group = ?
	ORDER BY rowid
	`
)

const (
	cmdInsertPlaylistItem = `INSERT INTO playlist (id, channels_group, channel, source) VALUES(?, ?, ?, ?)`

	cmdAppendChannelDisplayName = `INSERT INTO channel_display_names(cid, lang, display_name) VALUES(?, ?, ?)`

	cmdAppendChannelURL = `INSERT INTO channel_urls(cid, url) VALUES(?, ?)`

	cmdAppendGuideChannel   = `INSERT INTO channels(channel_id) VALUES(?)`
	cmdUpdateGuideChannelID = `UPDATE channels SET cid = ? WHERE rowid = ?`

	cmdAppendGuideProgramme = `INSERT INTO programme(channel_id, start, stop, pdc_start,
	vps_start, show_view, video_plus, clump_idx) VALUES(?, ?, ?, ?, ?, ?, ?, ?)`
	cmdUpdateGuideProgrammePID = `UPDATE programme SET pid = ? where rowid = ?`

	cmdAppendProgrammeTitle = `INSERT INTO programme_titles(pid, lang, title) VALUES(?, ?, ?)`
)

const (
	cmdDeleteFromChannels    = `DELETE FROM channels`
	cmdDeleteFromChannelsURL = `DELETE FROM channel_urls`
)

var db *sql.DB
var dbname = getPlaylistDatabaseName()

func init() {

	var err error

	db, err = sql.Open("sqlite3", dbname)

	if err != nil {
		log.Fatal(err)
	}

	err = createDatabaseStructure(db)

	if err != nil {
		log.Fatal(err)
	}

	/*
		pl = &Playlist{db: db}
		g = &Guide{db: db}
	*/
}

func createDatabaseStructure(db *sql.DB) (err error) {

	objects := [66]string{cmdCreateTablePlaylist, cmdCreateIndexPlaylistCID,
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
		cmdCreateTableProgrammeLanguage, cmdCreateIndexProgrammeLanguagePID,
		cmdCreateTableProgrammeOriginalLanguage, cmdCreateIndexProgrammeOriginalLanguagePID,
		cmdCreateTableProgrammeCountries, cmdCreateIndexProgrammeCountriesPID,
		cmdCreateTableProgrammeDirectors, cmdCreateIndexProgrammeDirectorsPID,
		cmdCreateTableProgrammeWriters, cmdCreateIndexProgrammeWritersPID,
		cmdCreateTableProgrammeAdapters, cmdCreateIndexProgrammeAdaptersPID,
		cmdCreateTableProgrammeProducers, cmdCreateIndexProgrammeProducersPID,
		cmdCreateTableProgrammeComposers, cmdCreateIndexProgrammeComposersPID,
		cmdCreateTableProgrammeEditors, cmdCreateIndexProgrammeEditorsPID,
		cmdCreateTableProgrammePresenters, cmdCreateIndexProgrammePresentersPID,
		cmdCreateTableProgrammeCommentators, cmdCreateIndexProgrammeCommentatorsPID,
		cmdCreateTableProgrammeGuests, cmdCreateIndexProgrammeGuestsPID,
		cmdCreateTableProgrammeActors, cmdCreateIndexProgrammeActorsPID,
		cmdCreateTableProgrammeLength, cmdCreateIndexProgrammeLengthPID,
		cmdCreateTableProgrammeIcon, cmdCreateIndexProgrammeIconPID,
		cmdCreateTableProgrammeEpisodeNum, cmdCreateIndexProgrammeEpisodeNumPID,
		cmdCreateTableProgrammeVideo, cmdCreateIndexProgrammeVideoPID,
		cmdCreateTableProgrammeAudio, cmdCreateIndexProgrammeAudioPID,
		cmdCreateTableProgrammePreviouslyShown, cmdCreateIndexProgrammePreviouslyShownPID,
		cmdCreateTableProgrammePremiere, cmdCreateIndexProgrammePremierePID,
		cmdCreateTableProgrammeLastChance, cmdCreateIndexProgrammeLastChancePID}

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
