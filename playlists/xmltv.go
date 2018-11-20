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
)

// Description of XMLTV guide format
// For more details see https://github.com/XMLTV/xmltv/blob/master/xmltv.dtd

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

// XMLTVChannelURL - an URL where you can find out more about the channel
type XMLTVChannelURL struct {
	XMLName xml.Name `xml:"url"`
	Value   string   `xml:",chardata"`
}

// XMLTVChannel - details of the channel
type XMLTVChannel struct {
	XMLName     xml.Name                  `xml:"channel"`
	ID          string                    `xml:"id,attr"`
	DisplayName []XMLTVChannelDisplayName `xml:"display-name"`
	URL         []XMLTVChannelURL         `xml:"url"`
}

// XMLTVProgramme - details of the single programme transmission
type XMLTVProgramme struct {
	XMLName           xml.Name                         `xml:"programme"`
	Channel           string                           `xml:"channel,attr"`
	Start             string                           `xml:"start,attr"`
	Stop              string                           `xml:"stop,attr"`
	PDCStart          string                           `xml:"pdc-start,attr"`
	VPSStart          string                           `xml:"vps-start,attr"`
	ShowView          string                           `xml:"showview,attr"`
	VideoPlus         string                           `xml:"videoplus,attr"`
	ClumpIdx          string                           `xml:"clumpidx,attr"`
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

// XMLTVProgrammeValue - the value of the element that contains it
type XMLTVProgrammeValue struct {
	XMLName xml.Name `xml:"value"`
	Value   string   `xml:",chardata"`
}

// XMLTVProgrammeTitle - programme title
type XMLTVProgrammeTitle struct {
	XMLName xml.Name `xml:"title"`
	Lang    string   `xml:"lang,attr"`
	Value   string   `xml:",chardata"`
}

// XMLTVProgrammeSubTitle - sub-title or episode title
type XMLTVProgrammeSubTitle struct {
	XMLName xml.Name `xml:"sub-title"`
	Lang    string   `xml:"lang,attr"`
	Value   string   `xml:",chardata"`
}

// XMLTVProgrammeDesc - description of the programme or episode
type XMLTVProgrammeDesc struct {
	XMLName xml.Name `xml:"desc"`
	Lang    string   `xml:"lang,attr"`
	Value   string   `xml:",chardata"`
}

// XMLTVProgrammeCredits - credits for the programme
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

// XMLTVProgrammeActor - item of the list of actors
type XMLTVProgrammeActor struct {
	XMLName xml.Name `xml:"actor"`
	Role    string   `xml:"role,attr"`
	Name    string   `xml:",chardata"`
}

// XMLTVProgrammeCategory - type of programme
type XMLTVProgrammeCategory struct {
	XMLName xml.Name `xml:"category"`
	Lang    string   `xml:"lang,attr"`
	Value   string   `xml:",chardata"`
}

// XMLTVProgrammeKeyword - keywords for the programme
type XMLTVProgrammeKeyword struct {
	XMLName xml.Name `xml:"keyword"`
	Lang    string   `xml:"lang,attr"`
	Value   string   `xml:",chardata"`
}

// XMLTVProgrammeLanguage - the language the programme will be broadcast in
type XMLTVProgrammeLanguage struct {
	XMLName xml.Name `xml:"language"`
	Lang    string   `xml:"lang,attr"`
	Value   string   `xml:",chardata"`
}

// XMLTVProgrammeOriginalLanguage - the original language, before dubbing
type XMLTVProgrammeOriginalLanguage struct {
	XMLName xml.Name `xml:"orig-language"`
	Lang    string   `xml:"lang,attr"`
	Value   string   `xnl:",chardata"`
}

// XMLTVProgrammeLength - the true length of the programme, not counting advertisements
// or trailers
type XMLTVProgrammeLength struct {
	XMLName xml.Name `xml:"length"`
	Units   string   `xml:"units,attr"`
	Value   string   `xml:",chardata"`
}

// XMLTVProgrammeIcon - an icon associated with the element that contains it
type XMLTVProgrammeIcon struct {
	XMLName xml.Name `xml:"icon"`
	Src     string   `xml:"src,attr"`
	Width   string   `xml:"width,attr"`
	Height  string   `xml:"height,attr"`
}

// XMLTVProgrammeCountry - the country where the programme was made or one of the countries in
// a joint production
type XMLTVProgrammeCountry struct {
	XMLName xml.Name `xml:"country"`
	Lang    string   `xml:"lang,attr"`
	Value   string   `xml:",chardata"`
}

// XMLTVProgrammeEpisodeNum - episode number
type XMLTVProgrammeEpisodeNum struct {
	XMLName xml.Name `xml:"episode-num"`
	System  string   `xml:"system,attr"`
	Value   string   `xml:",chardata"`
}

// XMLTVProgrammeVideo - video details
type XMLTVProgrammeVideo struct {
	XMLName xml.Name `xml:"video"`
	Present string   `xml:"present"`
	Colour  string   `xml:"colour"`
	Aspect  string   `xml:"aspect"`
	Quality string   `xml:"quality"`
}

// XMLTVProgrammeAudio - audio details
type XMLTVProgrammeAudio struct {
	XMLName xml.Name `xml:"audio"`
	Present string   `xml:"present"`
	Stereo  string   `xml:"stereo"`
}

// XMLTVProgrammePreviouslyShown - when and where the programme was
// last shown, if known
type XMLTVProgrammePreviouslyShown struct {
	XMLName xml.Name `xml:"previously-shown"`
	Start   string   `xml:"start,attr"`
	Channel string   `xml:"channel,attr"`
}

// XMLTVProgrammePremiere - premiere, if known
type XMLTVProgrammePremiere struct {
	XMLName xml.Name `xml:"premiere"`
	Lang    string   `xml:"lang,attr"`
	Value   string   `xml:",chardata"`
}

// XMLTVProgrammmeLastChance - in a way this is the opposite of premiere
type XMLTVProgrammmeLastChance struct {
	XMLName xml.Name `xml:"last-chance"`
	Lang    string   `xml:"lang,attr"`
	Value   string   `xml:",chardata"`
}

// XMLTVProgrammeSubtitles - subtitles
type XMLTVProgrammeSubtitles struct {
	XMLName  xml.Name                 `xml:"subtitles"`
	Type     string                   `xml:"type,attr"`
	Language []XMLTVProgrammeLanguage `xml:"language"`
}

// XMLTVProgrammeRating - rating
type XMLTVProgrammeRating struct {
	XMLName xml.Name            `xml:"rating"`
	System  string              `xml:"system"`
	Value   XMLTVProgrammeValue `xml:"value"`
	Icon    XMLTVProgrammeIcon  `xml:"icon"`
}

// XMLTVProgrammeStarRating - star rating
type XMLTVProgrammeStarRating struct {
	XMLName xml.Name            `xml:"star-rating"`
	System  string              `xml:"system"`
	Value   XMLTVProgrammeValue `xml:"value"`
	Icon    XMLTVProgrammeIcon  `xml:"icon"`
}

// XMLTVProgrammeReview - review
type XMLTVProgrammeReview struct {
	XMLName  xml.Name `xml:"review"`
	Type     string   `xml:"type,attr"`
	Source   string   `xml:"source,attr"`
	Reviewer string   `xml:"reviewer,attr"`
	Lang     string   `xml:"lang,attr"`
	Value    string   `xml:",chardata"`
}

// OnHeadEvent an event that fires when guide header is read
type OnHeadEvent func(h *XMLTVHead) error

// OnChannelEvent an event that fires when channel info is read
type OnChannelEvent func(ch *XMLTVChannel) error

// OnProgrammeEvent an event that fires when programme info is read
type OnProgrammeEvent func(p *XMLTVProgramme) error

// XMLTVParser is parser of tv guide with xmltv format (github.com/xmltv)
type XMLTVParser struct {
	OnHead      OnHeadEvent
	OnChannel   OnChannelEvent
	OnProgramme OnProgrammeEvent
}

func (parser *XMLTVParser) doHead(h *XMLTVHead) error {

	if parser.OnHead != nil {
		return parser.OnHead(h)
	}

	return nil
}

func (parser *XMLTVParser) doChannel(ch *XMLTVChannel) error {

	if parser.OnChannel != nil {
		return parser.OnChannel(ch)
	}

	return nil
}

func (parser *XMLTVParser) doProgramme(p *XMLTVProgramme) error {

	if parser.OnProgramme != nil {
		return parser.OnProgramme(p)
	}

	return nil
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
		token, _ = decoder.Token()

		if token == nil {
			return
		}

		switch elem := token.(type) {
		case xml.StartElement:

			ename = elem.Name.Local

			if ename == "tv" {

				var h XMLTVHead

				err = decoder.DecodeElement(&h, &elem)

				if err != nil {
					return
				}

				err = parser.doHead(&h)

				if err != nil {
					return
				}
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
		token, _ = decoder.Token()

		if token == nil {
			return
		}

		switch elem := token.(type) {
		case xml.StartElement:

			ename = elem.Name.Local

			switch ename {

			case "channel":

				var c XMLTVChannel

				err = decoder.DecodeElement(&c, &elem)

				if err != nil {
					return
				}

				err = parser.doChannel(&c)

				if err != nil {
					return
				}

			case "programme":

				var p XMLTVProgramme

				err = decoder.DecodeElement(&p, &elem)

				if err != nil {
					return
				}

				err = parser.doProgramme(&p)

				if err != nil {
					return
				}

			}

		default:
		}
	}
}
