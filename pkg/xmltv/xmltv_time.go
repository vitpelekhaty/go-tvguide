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

var layouts = [4]string{
	"20060102150405 -0700",
	"2006-01-02 15:04:05",
	"20060102150405",
	"2006"}

// TimeOfProgramme return result of time parsing
func TimeOfProgramme(st string) (time.Time, error) {

	var et time.Time

	if strings.Trim(st, " ") == "" {
		return et, nil
	}

	for _, layout := range layouts {

		if t, err := time.Parse(layout, st); err == nil {
			return t, err
		}
	}

	return et, fmt.Errorf(`string "%q" parsing error`, st)
}
