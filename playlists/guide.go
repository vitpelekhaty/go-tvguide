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
	pdb
	db   *sql.DB
	tx   *sql.Tx
	stmt map[string]*sql.Stmt
}

var g *Guide

var queries = map[string]string{
	"cmdAppendGuideChannel":              cmdAppendGuideChannel,
	"cmdUpdateGuideChannelID":            cmdUpdateGuideChannelID,
	"cmdAppendChannelDisplayName":        cmdAppendChannelDisplayName,
	"cmdAppendChannelURL":                cmdAppendChannelURL,
	"cmdAppendGuideProgramme":            cmdAppendGuideProgramme,
	"cmdUpdateGuideProgrammePID":         cmdUpdateGuideProgrammePID,
	"cmdAppendProgrammeTitle":            cmdAppendProgrammeTitle,
	"cmdAppendProgrammeSubTitle":         cmdAppendProgrammeSubTitle,
	"cmdAppendProgrammeDesc":             cmdAppendProgrammeDesc,
	"cmdAppendProgrammeDates":            cmdAppendProgrammeDates,
	"cmdAppendProgrammeCategories":       cmdAppendProgrammeCategories,
	"cmdAppendProgrammeKeywords":         cmdAppendProgrammeKeywords,
	"cmdAppendProgrammeLanguage":         cmdAppendProgrammeLanguage,
	"cmdAppendProgrammeOriginalLanguage": cmdAppendProgrammeOriginalLanguage,
	"cmdAppendProgrammeCountries":        cmdAppendProgrammeCountries,
	"cmdAppendProgrammeDirectors":        cmdAppendProgrammeDirectors,
	"cmdAppendProgrammeWriters":          cmdAppendProgrammeWriters,
	"cmdAppendProgrammeAdapters":         cmdAppendProgrammeAdapters,
	"cmdAppendProgrammeProducers":        cmdAppendProgrammeProducers,
	"cmdAppendProgrammeComposers":        cmdAppendProgrammeComposers,
	"cmdAppendProgrammeEditors":          cmdAppendProgrammeEditors,
	"cmdAppendProgrammePresenters":       cmdAppendProgrammePresenters,
	"cmdAppendProgrammeCommentators":     cmdAppendProgrammeCommentators,
	"cmdAppendProgrammeGuests":           cmdAppendProgrammeGuests,
	"cmdAppendProgrammeActors":           cmdAppendProgrammeActors,
	"cmdAppendProgrammeLength":           cmdAppendProgrammeLength,
	"cmdAppendProgrammeIcon":             cmdAppendProgrammeIcon,
	"cmdAppendProgrammeEpisodeNum":       cmdAppendProgrammeEpisodeNum,
	"cmdAppendProgrammeVideo":            cmdAppendProgrammeVideo,
	"cmdAppendProgrammeAudio":            cmdAppendProgrammeAudio,
	"cmdAppendProgrammePreviouslyShown":  cmdAppendProgrammePreviouslyShown,
	"cmdAppendProgrammePremiere":         cmdAppendProgrammePremiere,
	"cmdAppendProgrammeLastChance":       cmdAppendProgrammeLastChance,
	"cmdAppendProgrammeSubtitles":        cmdAppendProgrammeSubtitles,
	"cmdAppendProgrammeRating":           cmdAppendProgrammeRating,
	"cmdAppendProgrammeStarRating":       cmdAppendProgrammeStarRating,
	"cmdAppendProgrammeReview":           cmdAppendProgrammeReview,
	"cmdAppendProgrammeLangStat":         cmdAppendProgrammeLangStat}

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

	g.stmt = map[string]*sql.Stmt{}

	for key, command := range queries {

		var stmt *sql.Stmt

		stmt, err = tx.Prepare(command)

		if err != nil {
			return
		}

		g.stmt[key] = stmt
	}

	parser.OnChannel = func(ch *XMLTVChannel) error {
		return g.appendChannel(ch)
	}

	parser.OnProgramme = func(p *XMLTVProgramme) error {
		return g.appendProgramme(p)
	}

	if err = parser.Parse(data); err != nil {
		return
	}

	if err = g.appendProgrammeLangStat(); err != nil {
		return
	}

	err = g.analyze(g.db, g.tx)

	return
}

func (g *Guide) appendChannel(c *XMLTVChannel) (err error) {

	if c == nil {
		return errors.New("Guide.AppendChannel: cannot append an empty channel")
	}

	var cid int64

	cs := g.stmt["cmdAppendGuideChannel"]
	us := g.stmt["cmdUpdateGuideChannelID"]

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

	if len(c.DisplayName) > 0 {

		dn := make([]*XMLTVChannelDisplayName, len(c.DisplayName))

		for index, d := range c.DisplayName {
			dn[index] = &d
		}

		if err = g.appendChannelDisplayNames(cid, dn); err != nil {
			return err
		}
	}

	if len(c.URL) > 0 {

		urls := make([]*XMLTVChannelURL, len(c.URL))

		for index, url := range c.URL {
			urls[index] = &url
		}

		if err = g.appendChannelURL(cid, urls); err != nil {
			return err
		}
	}

	return nil
}

func (g *Guide) appendChannelDisplayNames(cid int64, d []*XMLTVChannelDisplayName) (err error) {

	for _, dn := range d {
		if _, err = g.stmt["cmdAppendChannelDisplayName"].Exec(&cid, &dn.Lang, &dn.Value); err != nil {
			return err
		}
	}

	return
}

func (g *Guide) appendChannelURL(cid int64, urls []*XMLTVChannelURL) (err error) {

	for _, url := range urls {
		if _, err = g.stmt["cmdAppendChannelURL"].Exec(&cid, &url.Value); err != nil {
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

	if err = g.checkAppendProgrammeTitle(pid, p); err != nil {
		return
	}

	if err = g.checkAppendProgrammeSubTitle(pid, p); err != nil {
		return
	}

	if err = g.checkAppendProgrammeDesc(pid, p); err != nil {
		return
	}

	if err = g.checkAppendProgrammeActors(pid, p); err != nil {
		return
	}

	if err = g.checkAppendProgrammeAdapters(pid, p); err != nil {
		return
	}

	if err = g.checkAppendProgrammeCommentators(pid, p); err != nil {
		return
	}

	if err = g.checkAppendProgrammeComposers(pid, p); err != nil {
		return
	}

	if err = g.checkAppendProgrammeDirectors(pid, p); err != nil {
		return
	}

	if err = g.checkAppendProgrammeEditors(pid, p); err != nil {
		return
	}

	if err = g.checkAppendProgrammeGuests(pid, p); err != nil {
		return
	}

	if err = g.checkAppendProgrammePresenters(pid, p); err != nil {
		return
	}

	if err = g.checkAppendProgrammeProducers(pid, p); err != nil {
		return
	}

	if err = g.checkAppendProgrammeWriters(pid, p); err != nil {
		return
	}

	if err = g.checkAppendProgrammeDates(pid, p); err != nil {
		return
	}

	if err = g.checkAppendProgrammeCategories(pid, p); err != nil {
		return
	}

	if err = g.checkAppendProgrammeKeywords(pid, p); err != nil {
		return
	}

	if err = g.checkAppendProgrammeLanguages(pid, p); err != nil {
		return
	}

	if err = g.checkAppendProgrammeOriginalLanguages(pid, p); err != nil {
		return
	}

	if err = g.checkAppendProgrammeLength(pid, p); err != nil {
		return
	}

	if err = g.checkAppendProgrammeIcon(pid, p); err != nil {
		return
	}

	if err = g.checkAppendProgrammeCountry(pid, p); err != nil {
		return
	}

	if err = g.checkAppendProgrammeEpisodeNum(pid, p); err != nil {
		return
	}

	if err = g.checkAppendProgrammeVideo(pid, p); err != nil {
		return
	}

	if err = g.checkAppendProgrammeAudio(pid, p); err != nil {
		return
	}

	if err = g.checkAppendProgrammePreviouslyShown(pid, p); err != nil {
		return
	}

	if err = g.checkAppendProgrammePremiere(pid, p); err != nil {
		return
	}

	if err = g.checkAppendProgrammeLastChance(pid, p); err != nil {
		return
	}

	if err = g.checkAppendProgrammeSubtitles(pid, p); err != nil {
		return
	}

	if err = g.checkAppendProgrammeRating(pid, p); err != nil {
		return
	}

	if err = g.checkAppendProgrammeStarRating(pid, p); err != nil {
		return
	}

	if err = g.checkAppendProgrammeReview(pid, p); err != nil {
		return
	}

	return
}

func (g *Guide) checkAppendProgrammeTitle(pid int64, p *XMLTVProgramme) (err error) {

	if len(p.Title) > 0 {

		titles := make([]*XMLTVProgrammeTitle, len(p.Title))

		for idx, t := range p.Title {
			titles[idx] = &t
		}

		if err = g.appendProgrammeTitle(pid, titles); err != nil {
			return
		}
	}

	return
}

func (g *Guide) checkAppendProgrammeSubTitle(pid int64, p *XMLTVProgramme) (err error) {

	if len(p.SubTitle) > 0 {

		subtitles := make([]*XMLTVProgrammeSubTitle, len(p.SubTitle))

		for idx, t := range p.SubTitle {
			subtitles[idx] = &t
		}

		if err = g.appendProgrammeSubTitle(pid, subtitles); err != nil {
			return
		}
	}

	return
}

func (g *Guide) checkAppendProgrammeDesc(pid int64, p *XMLTVProgramme) (err error) {

	if len(p.Desc) > 0 {

		desc := make([]*XMLTVProgrammeDesc, len(p.Desc))

		for idx, d := range p.Desc {
			desc[idx] = &d
		}

		if err = g.appendProgrammeDesc(pid, desc); err != nil {
			return
		}
	}

	return
}

func (g *Guide) checkAppendProgrammeActors(pid int64, p *XMLTVProgramme) (err error) {

	if len(p.Credits.Actors) > 0 {

		actors := make([]*XMLTVProgrammeActor, len(p.Credits.Actors))

		for idx, a := range p.Credits.Actors {
			actors[idx] = &a
		}

		if err = g.appendProgrammeActors(pid, actors); err != nil {
			return
		}
	}

	return
}

func (g *Guide) checkAppendProgrammeAdapters(pid int64, p *XMLTVProgramme) (err error) {

	if len(p.Credits.Adapters) > 0 {

		adapters := make([]*string, len(p.Credits.Adapters))

		for idx, a := range p.Credits.Adapters {
			adapters[idx] = &a
		}

		if err = g.appendProgrammeAdapters(pid, adapters); err != nil {
			return
		}
	}

	return
}

func (g *Guide) checkAppendProgrammeCommentators(pid int64, p *XMLTVProgramme) (err error) {

	if len(p.Credits.Commentators) > 0 {

		commentators := make([]*string, len(p.Credits.Commentators))

		for idx, c := range p.Credits.Commentators {
			commentators[idx] = &c
		}

		if err = g.appendProgrammeCommentators(pid, commentators); err != nil {
			return
		}
	}

	return
}

func (g *Guide) checkAppendProgrammeComposers(pid int64, p *XMLTVProgramme) (err error) {

	if len(p.Credits.Composers) > 0 {

		composers := make([]*string, len(p.Credits.Composers))

		for idx, c := range p.Credits.Composers {
			composers[idx] = &c
		}

		if err = g.appendProgrammeComposers(pid, composers); err != nil {
			return
		}
	}

	return
}

func (g *Guide) checkAppendProgrammeDirectors(pid int64, p *XMLTVProgramme) (err error) {

	if len(p.Credits.Directors) > 0 {

		directors := make([]*string, len(p.Credits.Directors))

		for idx, d := range p.Credits.Directors {
			directors[idx] = &d
		}

		if err = g.appendProgrammeDirectors(pid, directors); err != nil {
			return
		}
	}

	return
}

func (g *Guide) checkAppendProgrammeEditors(pid int64, p *XMLTVProgramme) (err error) {

	if len(p.Credits.Editors) > 0 {

		editors := make([]*string, len(p.Credits.Editors))

		for idx, e := range p.Credits.Editors {
			editors[idx] = &e
		}

		if err = g.appendProgrammeEditors(pid, editors); err != nil {
			return
		}
	}

	return
}

func (g *Guide) checkAppendProgrammeGuests(pid int64, p *XMLTVProgramme) (err error) {

	if len(p.Credits.Guests) > 0 {

		guests := make([]*string, len(p.Credits.Guests))

		for idx, gst := range p.Credits.Guests {
			guests[idx] = &gst
		}

		if err = g.appendProgrammeGuests(pid, guests); err != nil {
			return
		}
	}

	return
}

func (g *Guide) checkAppendProgrammePresenters(pid int64, p *XMLTVProgramme) (err error) {

	if len(p.Credits.Presenters) > 0 {

		presenters := make([]*string, len(p.Credits.Presenters))

		for idx, pr := range p.Credits.Presenters {
			presenters[idx] = &pr
		}

		if err = g.appendProgrammePresenters(pid, presenters); err != nil {
			return
		}
	}

	return
}

func (g *Guide) checkAppendProgrammeProducers(pid int64, p *XMLTVProgramme) (err error) {

	if len(p.Credits.Producers) > 0 {

		producers := make([]*string, len(p.Credits.Producers))

		for idx, pr := range p.Credits.Producers {
			producers[idx] = &pr
		}

		if err = g.appendProgrammeProducers(pid, producers); err != nil {
			return
		}
	}

	return
}

func (g *Guide) checkAppendProgrammeWriters(pid int64, p *XMLTVProgramme) (err error) {

	if len(p.Credits.Writers) > 0 {

		writers := make([]*string, len(p.Credits.Writers))

		for idx, w := range p.Credits.Writers {
			writers[idx] = &w
		}

		if err = g.appendProgrammeWriters(pid, writers); err != nil {
			return
		}
	}

	return
}

func (g *Guide) checkAppendProgrammeDates(pid int64, p *XMLTVProgramme) (err error) {

	if len(p.Dates) > 0 {

		dates := make([]*string, len(p.Dates))

		for idx, d := range p.Dates {
			dates[idx] = &d
		}

		if err = g.appendProgrammeDates(pid, dates); err != nil {
			return
		}
	}

	return
}

func (g *Guide) checkAppendProgrammeCategories(pid int64, p *XMLTVProgramme) (err error) {

	if len(p.Categories) > 0 {

		categories := make([]*XMLTVProgrammeCategory, len(p.Categories))

		for idx, c := range p.Categories {
			categories[idx] = &c
		}

		if err = g.appendProgrammeCategories(pid, categories); err != nil {
			return
		}
	}

	return
}

func (g *Guide) checkAppendProgrammeKeywords(pid int64, p *XMLTVProgramme) (err error) {

	if len(p.Keywords) > 0 {

		keywords := make([]*XMLTVProgrammeKeyword, len(p.Keywords))

		for idx, k := range p.Keywords {
			keywords[idx] = &k
		}

		if err = g.appendProgrammeKeywords(pid, keywords); err != nil {
			return
		}
	}

	return
}

func (g *Guide) checkAppendProgrammeLanguages(pid int64, p *XMLTVProgramme) (err error) {

	if len(p.Languages) > 0 {

		languages := make([]*XMLTVProgrammeLanguage, len(p.Languages))

		for idx, lang := range p.Languages {
			languages[idx] = &lang
		}

		if err = g.appendProgrammeLanguages(pid, languages); err != nil {
			return
		}
	}

	return
}

func (g *Guide) checkAppendProgrammeOriginalLanguages(pid int64, p *XMLTVProgramme) (err error) {

	if len(p.OriginalLanguages) > 0 {

		languages := make([]*XMLTVProgrammeOriginalLanguage, len(p.OriginalLanguages))

		for idx, lang := range p.OriginalLanguages {
			languages[idx] = &lang
		}

		if err = g.appendProgrammeOriginalLanguages(pid, languages); err != nil {
			return
		}
	}

	return
}

func (g *Guide) checkAppendProgrammeLength(pid int64, p *XMLTVProgramme) (err error) {

	if len(p.Length) > 0 {

		length := make([]*XMLTVProgrammeLength, len(p.Length))

		for idx, l := range p.Length {
			length[idx] = &l
		}

		if err = g.appendProgrammeLength(pid, length); err != nil {
			return
		}
	}

	return
}

func (g *Guide) checkAppendProgrammeIcon(pid int64, p *XMLTVProgramme) (err error) {

	if len(p.Icon) > 0 {

		icons := make([]*XMLTVProgrammeIcon, len(p.Icon))

		for idx, icon := range p.Icon {
			icons[idx] = &icon
		}

		if err = g.appendProgrammeIcon(pid, icons); err != nil {
			return
		}
	}

	return
}

func (g *Guide) checkAppendProgrammeCountry(pid int64, p *XMLTVProgramme) (err error) {

	if len(p.Country) > 0 {

		countries := make([]*XMLTVProgrammeCountry, len(p.Country))

		for idx, country := range p.Country {
			countries[idx] = &country
		}

		if err = g.appendProgrammeCountry(pid, countries); err != nil {
			return
		}
	}

	return
}

func (g *Guide) checkAppendProgrammeEpisodeNum(pid int64, p *XMLTVProgramme) (err error) {

	if len(p.EpisodeNum) > 0 {

		enums := make([]*XMLTVProgrammeEpisodeNum, len(p.EpisodeNum))

		for idx, enum := range p.EpisodeNum {
			enums[idx] = &enum
		}

		if err = g.appendProgrammeEpisodeNum(pid, enums); err != nil {
			return
		}
	}

	return
}

func (g *Guide) checkAppendProgrammeVideo(pid int64, p *XMLTVProgramme) (err error) {

	if len(p.Video) > 0 {

		video := make([]*XMLTVProgrammeVideo, len(p.Video))

		for idx, v := range p.Video {
			video[idx] = &v
		}

		if err = g.appendProgrammeVideo(pid, video); err != nil {
			return
		}
	}

	return
}

func (g *Guide) checkAppendProgrammeAudio(pid int64, p *XMLTVProgramme) (err error) {

	if len(p.Audio) > 0 {

		audio := make([]*XMLTVProgrammeAudio, len(p.Audio))

		for idx, a := range p.Audio {
			audio[idx] = &a
		}

		if err = g.appendProgrammeAudio(pid, audio); err != nil {
			return
		}
	}

	return
}

func (g *Guide) checkAppendProgrammePreviouslyShown(pid int64, p *XMLTVProgramme) (err error) {

	if len(p.PreviouslyShown) > 0 {

		shown := make([]*XMLTVProgrammePreviouslyShown, len(p.PreviouslyShown))

		for idx, s := range p.PreviouslyShown {
			shown[idx] = &s
		}

		if err = g.appendProgrammePreviouslyShown(pid, shown); err != nil {
			return
		}
	}

	return
}

func (g *Guide) checkAppendProgrammePremiere(pid int64, p *XMLTVProgramme) (err error) {

	if len(p.Premiere) > 0 {

		premiere := make([]*XMLTVProgrammePremiere, len(p.Premiere))

		for idx, prem := range p.Premiere {
			premiere[idx] = &prem
		}

		if err = g.appendProgrammePremiere(pid, premiere); err != nil {
			return
		}
	}

	return
}

func (g *Guide) checkAppendProgrammeLastChance(pid int64, p *XMLTVProgramme) (err error) {

	if len(p.LastChance) > 0 {

		lc := make([]*XMLTVProgrammmeLastChance, len(p.LastChance))

		for idx, l := range p.LastChance {
			lc[idx] = &l
		}

		if err = g.appendProgrammeLastChance(pid, lc); err != nil {
			return
		}
	}

	return
}

func (g *Guide) checkAppendProgrammeSubtitles(pid int64, p *XMLTVProgramme) (err error) {

	if len(p.Subtitles) > 0 {

		subtitles := make([]*XMLTVProgrammeSubtitles, len(p.Subtitles))

		for idx, s := range p.Subtitles {
			subtitles[idx] = &s
		}

		if err = g.appendProgrammeSubtitles(pid, subtitles); err != nil {
			return
		}
	}

	return
}

func (g *Guide) checkAppendProgrammeRating(pid int64, p *XMLTVProgramme) (err error) {

	if len(p.Rating) > 0 {

		rating := make([]*XMLTVProgrammeRating, len(p.Rating))

		for idx, r := range p.Rating {
			rating[idx] = &r
		}

		if err = g.appendProgrammeRating(pid, rating); err != nil {
			return
		}
	}

	return
}

func (g *Guide) checkAppendProgrammeStarRating(pid int64, p *XMLTVProgramme) (err error) {

	if len(p.StarRating) > 0 {

		rating := make([]*XMLTVProgrammeStarRating, len(p.StarRating))

		for idx, r := range p.StarRating {
			rating[idx] = &r
		}

		if err = g.appendProgrammeStarRating(pid, rating); err != nil {
			return
		}
	}

	return
}

func (g *Guide) checkAppendProgrammeReview(pid int64, p *XMLTVProgramme) (err error) {

	if len(p.Review) > 0 {

		review := make([]*XMLTVProgrammeReview, len(p.Review))

		for idx, r := range p.Review {
			review[idx] = &r
		}

		if err = g.appendProgrammeReview(pid, review); err != nil {
			return
		}
	}

	return
}

func (g *Guide) appendProgrammeRecord(p *XMLTVProgramme) (int64, error) {

	var pid int64 = -1

	cs := g.stmt["cmdAppendGuideProgramme"]
	us := g.stmt["cmdUpdateGuideProgrammePID"]

	start, err := timeOfProgramme(p.Start)

	if err != nil {
		return pid, err
	}

	stop, err := timeOfProgramme(p.Stop)

	if err != nil {
		return pid, err
	}

	res, err := cs.Exec(&p.Channel, &start, &stop, &p.PDCStart, &p.VPSStart,
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
		if _, err = g.stmt["cmdAppendProgrammeTitle"].Exec(&pid, &t.Lang, &t.Value); err != nil {
			return
		}
	}

	return
}

func (g *Guide) appendProgrammeSubTitle(pid int64, subtitles []*XMLTVProgrammeSubTitle) (err error) {

	for _, s := range subtitles {
		if _, err = g.stmt["cmdAppendProgrammeSubTitle"].Exec(&pid, &s.Lang, &s.Value); err != nil {
			return
		}
	}

	return
}

func (g *Guide) appendProgrammeDesc(pid int64, desc []*XMLTVProgrammeDesc) (err error) {

	for _, d := range desc {
		if _, err = g.stmt["cmdAppendProgrammeDesc"].Exec(&pid, &d.Lang, &d.Value); err != nil {
			return
		}
	}

	return
}

func (g *Guide) appendProgrammeActors(pid int64, actors []*XMLTVProgrammeActor) (err error) {

	for _, a := range actors {
		if _, err = g.stmt["cmdAppendProgrammeActors"].Exec(&pid, &a.Name, &a.Role); err != nil {
			return
		}
	}

	return
}

func (g *Guide) appendProgrammeAdapters(pid int64, adapters []*string) (err error) {

	for _, a := range adapters {
		if _, err = g.stmt["cmdAppendProgrammeAdapters"].Exec(&pid, &a); err != nil {
			return
		}
	}

	return
}

func (g *Guide) appendProgrammeCommentators(pid int64, commentators []*string) (err error) {

	for _, c := range commentators {
		if _, err = g.stmt["cmdAppendProgrammeCommentators"].Exec(&pid, &c); err != nil {
			return
		}
	}

	return
}

func (g *Guide) appendProgrammeComposers(pid int64, composers []*string) (err error) {

	for _, c := range composers {
		if _, err = g.stmt["cmdAppendProgrammeComposers"].Exec(&pid, &c); err != nil {
			return
		}
	}

	return
}

func (g *Guide) appendProgrammeDirectors(pid int64, directors []*string) (err error) {

	for _, d := range directors {
		if _, err = g.stmt["cmdAppendProgrammeDirectors"].Exec(&pid, &d); err != nil {
			return
		}
	}

	return
}

func (g *Guide) appendProgrammeEditors(pid int64, editors []*string) (err error) {

	for _, e := range editors {
		if _, err = g.stmt["cmdAppendProgrammeEditors"].Exec(&pid, &e); err != nil {
			return
		}
	}

	return
}

func (g *Guide) appendProgrammeGuests(pid int64, guests []*string) (err error) {

	for _, gst := range guests {
		if _, err = g.stmt["cmdAppendProgrammeGuests"].Exec(&pid, &gst); err != nil {
			return
		}
	}

	return
}

func (g *Guide) appendProgrammePresenters(pid int64, presenters []*string) (err error) {

	for _, p := range presenters {
		if _, err = g.stmt["cmdAppendProgrammePresenters"].Exec(&pid, &p); err != nil {
			return
		}
	}

	return
}

func (g *Guide) appendProgrammeProducers(pid int64, producers []*string) (err error) {

	for _, p := range producers {
		if _, err = g.stmt["cmdAppendProgrammeProducers"].Exec(&pid, &p); err != nil {
			return
		}
	}

	return
}

func (g *Guide) appendProgrammeWriters(pid int64, writers []*string) (err error) {

	for _, w := range writers {
		if _, err = g.stmt["cmdAppendProgrammeWriters"].Exec(&pid, &w); err != nil {
			return
		}
	}

	return
}

func (g *Guide) appendProgrammeDates(pid int64, dates []*string) (err error) {

	for _, d := range dates {

		if _, err = g.stmt["cmdAppendProgrammeDates"].Exec(&pid, &d); err != nil {
			return
		}
	}

	return
}

func (g *Guide) appendProgrammeCategories(pid int64, categories []*XMLTVProgrammeCategory) (err error) {

	for _, c := range categories {
		if _, err = g.stmt["cmdAppendProgrammeCategories"].Exec(&pid, &c.Lang, &c.Value); err != nil {
			return
		}
	}

	return
}

func (g *Guide) appendProgrammeKeywords(pid int64, keywords []*XMLTVProgrammeKeyword) (err error) {

	for _, k := range keywords {
		if _, err = g.stmt["cmdAppendProgrammeKeywords"].Exec(&pid, &k.Lang, &k.Value); err != nil {
			return
		}
	}

	return
}

func (g *Guide) appendProgrammeLanguages(pid int64, languages []*XMLTVProgrammeLanguage) (err error) {

	for _, lang := range languages {
		if _, err = g.stmt["cmdAppendProgrammeLanguage"].Exec(&pid, &lang.Lang, &lang.Value); err != nil {
			return
		}
	}

	return
}

func (g *Guide) appendProgrammeOriginalLanguages(pid int64, languages []*XMLTVProgrammeOriginalLanguage) (err error) {

	for _, lang := range languages {
		if _, err = g.stmt["cmdAppendProgrammeOriginalLanguage"].Exec(&pid, &lang.Lang, &lang.Value); err != nil {
			return
		}
	}

	return
}

func (g *Guide) appendProgrammeLength(pid int64, length []*XMLTVProgrammeLength) (err error) {

	for _, l := range length {
		if _, err = g.stmt["cmdAppendProgrammeLength"].Exec(&pid, &l.Value, &l.Units); err != nil {
			return
		}
	}

	return
}

func (g *Guide) appendProgrammeIcon(pid int64, icons []*XMLTVProgrammeIcon) (err error) {

	for _, icon := range icons {
		if _, err = g.stmt["cmdAppendProgrammeIcon"].Exec(&pid, &icon.Src, &icon.Width, &icon.Height); err != nil {
			return
		}
	}

	return
}

func (g *Guide) appendProgrammeCountry(pid int64, countries []*XMLTVProgrammeCountry) (err error) {

	for _, country := range countries {
		if _, err = g.stmt["cmdAppendProgrammeCountries"].Exec(&pid, &country.Lang, &country.Value); err != nil {
			return
		}
	}

	return
}

func (g *Guide) appendProgrammeEpisodeNum(pid int64, enums []*XMLTVProgrammeEpisodeNum) (err error) {

	for _, enum := range enums {
		if _, err = g.stmt["cmdAppendProgrammeEpisodeNum"].Exec(&pid, &enum.System, &enum.Value); err != nil {
			return
		}
	}

	return
}

func (g *Guide) appendProgrammeVideo(pid int64, video []*XMLTVProgrammeVideo) (err error) {

	for _, v := range video {
		if _, err = g.stmt["cmdAppendProgrammeVideo"].Exec(&pid, &v.Present, &v.Colour, &v.Aspect, &v.Quality); err != nil {
			return
		}
	}

	return
}

func (g *Guide) appendProgrammeAudio(pid int64, audio []*XMLTVProgrammeAudio) (err error) {

	for _, a := range audio {
		if _, err = g.stmt["cmdAppendProgrammeAudio"].Exec(&pid, &a.Present, &a.Stereo); err != nil {
			return
		}
	}

	return
}

func (g *Guide) appendProgrammePreviouslyShown(pid int64, shown []*XMLTVProgrammePreviouslyShown) (err error) {

	for _, s := range shown {
		if _, err = g.stmt["cmdAppendProgrammePreviouslyShown"].Exec(&pid, &s.Start, &s.Channel); err != nil {
			return
		}
	}

	return
}

func (g *Guide) appendProgrammePremiere(pid int64, premiere []*XMLTVProgrammePremiere) (err error) {

	for _, prem := range premiere {
		if _, err = g.stmt["cmdAppendProgrammePremiere"].Exec(&pid, &prem.Lang, &prem.Value); err != nil {
			return
		}
	}

	return
}

func (g *Guide) appendProgrammeLastChance(pid int64, ls []*XMLTVProgrammmeLastChance) (err error) {

	for _, l := range ls {
		if _, err = g.stmt["cmdAppendProgrammeLastChance"].Exec(&pid, &l.Lang, &l.Value); err != nil {
			return
		}
	}

	return
}

func (g *Guide) appendProgrammeSubtitles(pid int64, subtitles []*XMLTVProgrammeSubtitles) (err error) {

	var stype string

	for _, st := range subtitles {

		stype = st.Type

		if len(st.Language) > 0 {

			for _, lang := range st.Language {

				if _, err = g.stmt["cmdAppendProgrammeSubtitles"].Exec(&pid, &stype, &lang.Lang, &lang.Value); err != nil {
					return
				}
			}
		}
	}

	return
}

func (g *Guide) appendProgrammeRating(pid int64, rating []*XMLTVProgrammeRating) (err error) {

	var (
		s     string
		value string
		src   string
		w     string
		h     string
	)

	for _, r := range rating {

		s = r.System

		value = r.Value.Value
		src = r.Icon.Src
		w = r.Icon.Width
		h = r.Icon.Height

		if _, err = g.stmt["cmdAppendProgrammeRating"].Exec(&pid, &s, &value, &src, &w, &h); err != nil {
			return
		}
	}

	return
}

func (g *Guide) appendProgrammeStarRating(pid int64, rating []*XMLTVProgrammeStarRating) (err error) {

	var (
		s     string
		value string
		src   string
		w     string
		h     string
	)

	for _, r := range rating {

		s = r.System

		value = r.Value.Value
		src = r.Icon.Src
		w = r.Icon.Width
		h = r.Icon.Height

		if _, err = g.stmt["cmdAppendProgrammeStarRating"].Exec(&pid, &s, &value, &src, &w, &h); err != nil {
			return
		}
	}

	return
}

func (g *Guide) appendProgrammeReview(pid int64, review []*XMLTVProgrammeReview) (err error) {

	for _, r := range review {
		if _, err = g.stmt["cmdAppendProgrammeReview"].Exec(&pid, &r.Type, &r.Source, r.Reviewer, &r.Lang, &r.Value); err != nil {
			return
		}
	}

	return
}

func (g *Guide) appendProgrammeLangStat() (err error) {

	_, err = g.stmt["cmdAppendProgrammeLangStat"].Exec()
	return
}
