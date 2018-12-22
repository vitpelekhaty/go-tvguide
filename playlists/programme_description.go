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

	"github.com/logrusorgru/aurora"
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

const fProgrammeDescription = "\n\t%s\n\t%s\n\n\t%s\n\t%s\n"

// ToString return the textual description of the programme
func (pd *ProgrammeDescription) ToString() string {

	var (
		text     string
		title    string
		subtitle string
		duration string
		desc     string
	)

	title = pd.Title
	subtitle = pd.SubTitle
	desc = strings.Repeat("lorem ipsum ", 1000) //pd.Description

	text = fmt.Sprintf(fProgrammeDescription, aurora.Bold(title), duration, subtitle, desc)

	return text
}
