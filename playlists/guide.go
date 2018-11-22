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
)

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
	stmtAppendProgrammeTitle     *sql.Stmt
}

var g *Guide

// CurrentGuide returns guide object
func CurrentGuide() *Guide {

	if g == nil {
		g = &Guide{db: db}
	}

	return g
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

	g.stmtAppendProgrammeTitle, err = tx.Prepare(cmdAppendProgrammeTitle)

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

	var pid int64

	pid, err = g.appendProgrammeRecord(p)

	if err != nil {
		return
	}

	if len(p.Title) > 0 {

		titles := make([]*XMLTVProgrammeTitle, len(p.Title))

		for idx, t := range p.Title {
			titles[idx] = &t
		}

		err = g.appendProgrammeTitle(pid, titles)

		if err != nil {
			return
		}
	}

	return
}

func (g *Guide) appendProgrammeRecord(p *XMLTVProgramme) (int64, error) {

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

func (g *Guide) appendProgrammeTitle(pid int64, titles []*XMLTVProgrammeTitle) (err error) {

	for _, t := range titles {
		_, err = g.stmtAppendProgrammeTitle.Exec(&pid, &t.Lang, &t.Value)
	}

	return nil
}
