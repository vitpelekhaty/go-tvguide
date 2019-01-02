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
	"fmt"
	"strings"
	"time"
)

// ProgrammeRating - programme rating
type ProgrammeRating struct {
	System string
	Rating string
}

// ProgrammeActor - actor
type ProgrammeActor struct {
	Actor string
	Role  string
}

// ProgrammeDescription - description of the programme
type ProgrammeDescription struct {
	Programme
	SubTitle    string
	Description string
	Category    []*string
	Country     []*string
	Directors   []*string
	Actors      []*ProgrammeActor
	Rating      []*ProgrammeRating
}

// ProgrammeTimeDescription return text description of the programme start and stop times and duration
func (pd *ProgrammeDescription) ProgrammeTimeDescription(timeCurrent time.Time) string {

	var (
		inTime bool
		d, cd  time.Duration
	)

	inTime = (timeCurrent.After(pd.Start) || timeCurrent.Equal(pd.Start)) && timeCurrent.Before(pd.Stop)

	d = pd.Stop.Sub(pd.Start)
	cd = pd.Stop.Sub(timeCurrent)

	if inTime {
		return fmt.Sprintf("%02d:%02d - %02d:%02d / %vmin / +%vmin", pd.StartHour(), pd.StartMinute(), pd.StopHour(), pd.StopMinute(),
			d.Minutes(), cd.Minutes())
	}

	return fmt.Sprintf("%02d:%02d - %02d:%02d / %vmin", pd.StartHour(), pd.StartMinute(), pd.StopHour(), pd.StopMinute(),
		d.Minutes())
}

// ProgrammeDirectors returns list of directors of the programme represented as string
func (pd *ProgrammeDescription) ProgrammeDirectors() string {

	directors := make([]string, len(pd.Directors))

	for index, director := range pd.Directors {
		directors[index] = *director
	}

	if len(directors) > 0 {
		return strings.Join(directors, ", ")
	}

	return ""
}

// ProgrammeActors returns list of actors represented as string
func (pd *ProgrammeDescription) ProgrammeActors() string {

	actors := make([]string, len(pd.Actors))

	for index, actor := range pd.Actors {
		actors[index] = actor.ToString()
	}

	if len(actors) > 0 {
		return strings.Join(actors, ", ")
	}

	return ""
}

// ProgrammeRatings returns rating of the programme
func (pd *ProgrammeDescription) ProgrammeRatings() string {

	ratings := make([]string, len(pd.Rating))

	for index, rating := range pd.Rating {
		ratings[index] = rating.ToString()
	}

	if len(ratings) > 0 {
		return strings.Join(ratings, ", ")
	}

	return ""

}

// ProgrammeCountries returns list of countries represented as string
func (pd *ProgrammeDescription) ProgrammeCountries() string {

	countries := make([]string, len(pd.Country))

	for index, country := range pd.Country {
		countries[index] = *country
	}

	if len(countries) > 0 {
		return strings.Join(countries, ", ")
	}

	return ""
}

// ProgrammeCategories returns list of categories represented as string
func (pd *ProgrammeDescription) ProgrammeCategories() string {

	categories := make([]string, len(pd.Category))

	for index, category := range pd.Category {
		categories[index] = *category
	}

	if len(categories) > 0 {
		return strings.Join(categories, ", ")
	}

	return ""
}

// ToString returns actor and role represented as string
func (pa *ProgrammeActor) ToString() string {

	if strings.Trim(pa.Actor, " ") == "" {
		return ""
	}

	if strings.Trim(pa.Role, " ") == "" {
		return pa.Actor
	}

	return fmt.Sprintf("%s (%s)", pa.Actor, pa.Role)
}

// ToString returns rating of the programme represented as string
func (pr *ProgrammeRating) ToString() string {

	if strings.Trim(pr.Rating, " ") == "" {
		return ""
	}

	if strings.Trim(pr.System, " ") == "" {
		return pr.Rating
	}

	return fmt.Sprintf("%s (%s)", pr.Rating, pr.System)

}
