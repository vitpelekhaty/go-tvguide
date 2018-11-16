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
	"bytes"
	"encoding/xml"
	"io"
)

// XMLTVHead - the root element
type XMLTVHead struct {
	XMLName           xml.Name `xml:"tv"`
	GeneratorInfoName string   `xml:"generator-info-name,attr"`
	GeneratorInfoURL  string   `xml:"generator-info-url,attr"`
	SourceInfoURL     string   `xml:"source-info-url,attr"`
	SourceInfoName    string   `xml:"source-info-name,attr"`
	SourceDataURL     string   `xml:"source-data-url,attr"`
}

// XMLTVChannelDisplayName - a user-friendly name for the channel
type XMLTVChannelDisplayName struct {
	XMLName xml.Name `xml:"display-name"`
	Lang    string   `xml:"lang,attr"`
	Value   string   `xml:",chardata"`
}

// XMLTVChannelURL contains URLs of the channel
type XMLTVChannelURL struct {
	XMLName xml.Name `xml:"url"`
	Value   string   `xml:",chardata"`
}

// XMLTVChannel contains info about channel
type XMLTVChannel struct {
	XMLName     xml.Name                  `xml:"channel"`
	ID          string                    `xml:"id,attr"`
	DisplayName []XMLTVChannelDisplayName `xml:"display-name"`
	URL         []XMLTVChannelURL         `xml:"url"`
}

// XMLTVProgramme contains info about tv programme
type XMLTVProgramme struct {
	XMLName           xml.Name                         `xml:"programme"`
	Channel           string                           `xml:"channel"`
	Start             string                           `xml:"start"`
	Stop              string                           `xml:"stop"`
	PDCStart          string                           `xml:"pdc-start"`
	VPSStart          string                           `xml:"vps-start"`
	ShowView          string                           `xml:"showview"`
	VideoPlus         string                           `xml:"videoplus"`
	ClumpIdx          string                           `xml:"clumpidx"`
	Title             []XMLTVProgrammeTitle            `xml:"title"`
	SubTitle          []XMLTVProgrammeSubTitle         `xml:"sub-title"`
	Desc              []XMLTVProgrammeDesc             `xml:"desc"`
	Credits           XMLTVProgrammeCredits            `xml:"credits"`
	Dates             []string                         `xml:"date"`
	Categories        []XMLTVProgrammeCategory         `xml:"category"`
	Keywords          []XMLTVProgrammeKeyword          `xml:"keyword"`
	Languages         []XMLTVProgrammeLanguage         `xml:"language"`
	OriginalLanguages []XMLTVProgrammeOriginalLanguage `xml:"orig-language"`
	Length            XMLTVProgrammeLength             `xml:"length"`
	Icon              XMLTVProgrammeIcon               `xml:"icon"`
	Country           []XMLTVProgrammeCountry          `xml:"country"`
	EpisodeNum        XMLTVProgrammeEpisodeNum         `xml:"episode-num"`
	Video             XMLTVProgrammeVideo              `xml:"video"`
	Audio             XMLTVProgrammeAudio              `xml:"audio"`
	PreviouslyShown   XMLTVProgrammePreviouslyShown    `xml:"previously-shown"`
	Premiere          XMLTVProgrammePremiere           `xml:"premiere"`
	LastChance        XMLTVProgrammmeLastChance        `xml:"last-chance"`
	Subtitles         []XMLTVProgrammeSubtitles        `xml:"subtitles"`
	Rating            []XMLTVProgrammeRating           `xml:"rating"`
	StarRating        []XMLTVProgrammeStarRating       `xml:"star-rating"`
	Review            []XMLTVProgrammeReview           `xml:"review"`
}

type XMLTVProgrammeValue struct {
	XMLName xml.Name `xml:"value"`
	Value   string   `xml:",chardata"`
}

// XMLTVProgrammeTitle contans title of the programme
type XMLTVProgrammeTitle struct {
	XMLName xml.Name `xml:"title"`
	Lang    string   `xml:"lang,attr"`
	Value   string   `xml:",chardata"`
}

type XMLTVProgrammeSubTitle struct {
	XMLName xml.Name `xml:"sub-title"`
	Lang    string   `xml:"lang,attr"`
	Value   string   `xml:",chardata"`
}

type XMLTVProgrammeDesc struct {
	XMLName xml.Name `xml:"desc"`
	Lang    string   `xml:"lang,attr"`
	Value   string   `xml:",chardata"`
}

type XMLTVProgrammeCredits struct {
	XMLName      xml.Name              `xml:"credits"`
	Directors    []string              `xml:"director"`
	Actors       []XMLTVProgrammeActor `xml:"actor"`
	Writers      []string              `xml:"writer"`
	Adapters     []string              `xml:"adapter"`
	Producers    []string              `xml:"producer"`
	Composers    []string              `xml:"composer"`
	Editors      []string              `xml:"editor"`
	Presenters   []string              `xml:"presenter"`
	Commentators []string              `xml:"commentator"`
	Guests       []string              `xml:"guest"`
}

type XMLTVProgrammeActor struct {
	XMLName xml.Name `xml:"actor"`
	Role    string   `xml:"role,attr"`
	Name    string   `xml:",chardata"`
}

type XMLTVProgrammeCategory struct {
	XMLName xml.Name `xml:"category"`
	Lang    string   `xml:"lang,attr"`
	Value   string   `xml:",chardata"`
}

type XMLTVProgrammeKeyword struct {
	XMLName xml.Name `xml:"keyword"`
	Lang    string   `xml:"lang,attr"`
	Value   string   `xml:",chardata"`
}

type XMLTVProgrammeLanguage struct {
	XMLName xml.Name `xml:"language"`
	Lang    string   `xml:"lang,attr"`
	Value   string   `xml:",chardata"`
}

type XMLTVProgrammeOriginalLanguage struct {
	XMLName xml.Name `xml:"orig-language"`
	Lang    string   `xml:"lang,attr"`
	Value   string   `xnl:",chardata"`
}

type XMLTVProgrammeLength struct {
	XMLName xml.Name `xml:"length"`
	Units   string   `xml:"units,attr"`
	Value   string   `xml:",chardata"`
}

type XMLTVProgrammeIcon struct {
	XMLName xml.Name `xml:"icon"`
	Src     string   `xml:"src,attr"`
	Width   string   `xml:"width,attr"`
	Height  string   `xml:"height,attr"`
}

type XMLTVProgrammeCountry struct {
	XMLName xml.Name `xml:"country"`
	Lang    string   `xml:"lang,attr"`
	Value   string   `xml:",chardata"`
}

type XMLTVProgrammeEpisodeNum struct {
	XMLName xml.Name `xml:"episode-num"`
	System  string   `xml:"system,attr"`
	Value   string   `xml:",chardata"`
}

type XMLTVProgrammeVideo struct {
	XMLName xml.Name `xml:"video"`
	Present string   `xml:"present"`
	Colour  string   `xml:"colour"`
	Aspect  string   `xml:"aspect"`
	Quality string   `xml:"quality"`
}

type XMLTVProgrammeAudio struct {
	XMLName xml.Name `xml:"audio"`
	Present string   `xml:"present"`
	Stereo  string   `xml:"stereo"`
}

type XMLTVProgrammePreviouslyShown struct {
	XMLName xml.Name `xml:"previously-shown"`
	Start   string   `xml:"start,attr"`
	Channel string   `xml:"channel,attr"`
}

type XMLTVProgrammePremiere struct {
	XMLName xml.Name `xml:"premiere"`
	Lang    string   `xml:"lang,attr"`
	Value   string   `xml:",chardata"`
}

type XMLTVProgrammmeLastChance struct {
	XMLName xml.Name `xml:"last-chance"`
	Lang    string   `xml:"lang,attr"`
	Value   string   `xml:",chardata"`
}

type XMLTVProgrammeSubtitles struct {
	XMLName  xml.Name                 `xml:"subtitles"`
	Type     string                   `xml:"type,attr"`
	Language []XMLTVProgrammeLanguage `xml:"language"`
}

type XMLTVProgrammeRating struct {
	XMLName xml.Name            `xml:"rating"`
	System  string              `xml:"system"`
	Value   XMLTVProgrammeValue `xml:"value"`
	Icon    XMLTVProgrammeIcon  `xml:"icon"`
}

type XMLTVProgrammeStarRating struct {
	XMLName xml.Name            `xml:"star-rating"`
	System  string              `xml:"system"`
	Value   XMLTVProgrammeValue `xml:"value"`
	Icon    XMLTVProgrammeIcon  `xml:"icon"`
}

type XMLTVProgrammeReview struct {
	XMLName  xml.Name `xml:"review"`
	Type     string   `xml:"type,attr"`
	Source   string   `xml:"source,attr"`
	Reviewer string   `xml:"reviewer,attr"`
	Lang     string   `xml:"lang,attr"`
	Value    string   `xml:",chardata"`
}

// OnHeadEvent an event that fires when guide header is read
type OnHeadEvent func(h *XMLTVHead)

// OnChannelEvent an event that fires when channel info is read
type OnChannelEvent func(ch *XMLTVChannel)

// OnProgrammeEvent an event that fires when programme info is read
type OnProgrammeEvent func(p *XMLTVProgramme)

// XMLTVParser is parser of tv guide with xmltv format (github.com/xmltv)
type XMLTVParser struct {
	OnHead      OnHeadEvent
	OnChannel   OnChannelEvent
	OnProgramme OnProgrammeEvent
}

func (parser *XMLTVParser) doHead(h *XMLTVHead) {

	if parser.OnHead != nil {
		parser.OnHead(h)
	}
}

func (parser *XMLTVParser) doChannel(ch *XMLTVChannel) {

	if parser.OnChannel != nil {
		parser.OnChannel(ch)
	}
}

func (parser *XMLTVParser) doProgramme(p *XMLTVProgramme) {

	if parser.OnProgramme != nil {
		parser.OnProgramme(p)
	}
}

// Parse parses XMLTV guide data
func (parser *XMLTVParser) Parse(data []byte) error {

	err := parser.parseHeader(data)

	if err != nil {
		return err
	}

	err = parser.parseTVListing(data)

	return err
}

func (parser *XMLTVParser) parseHeader(data []byte) (err error) {

	r := bytes.NewReader(data)
	decoder := xml.NewDecoder(r)

	var ename string
	var token xml.Token

	for {
		token, err = decoder.Token()

		if token == nil || (err != nil && err != io.EOF) {
			return err
		}

		switch elem := token.(type) {
		case xml.StartElement:

			ename = elem.Name.Local

			if ename == "tv" {

				var h XMLTVHead

				_ = decoder.DecodeElement(&h, &elem)
				parser.doHead(&h)
			}

		default:
		}
	}
}

func (parser *XMLTVParser) parseTVListing(data []byte) (err error) {

	r := bytes.NewReader(data)
	decoder := xml.NewDecoder(r)

	var ename string
	var token xml.Token

	for {
		token, err = decoder.Token()

		if token == nil || (err != nil && err != io.EOF) {
			return err
		}

		switch elem := token.(type) {
		case xml.StartElement:

			ename = elem.Name.Local

			switch ename {

			case "channel":

				var c XMLTVChannel

				_ = decoder.DecodeElement(&c, &elem)
				parser.doChannel(&c)

			case "programme":

				var p XMLTVProgramme

				_ = decoder.DecodeElement(&p, &elem)
				parser.doProgramme(&p)

			}

		default:
		}
	}
}
