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
	"testing"
	"time"
)

type wanttime struct {
	year     int
	month    time.Month
	day      int
	hour     int
	minute   int
	second   int
	tzname   string
	tzoffset int
}

func TestTimeOfProgramme(t *testing.T) {

	var tests = []struct {
		input string
		want  wanttime
	}{
		{"20181027030000", wanttime{year: 2018, month: time.October, day: 27, hour: 3, minute: 0, second: 0, tzname: "UTC"}},
		{"20181027030000 +0300", wanttime{year: 2018, month: time.October, day: 27, hour: 3, minute: 0, second: 0, tzoffset: 10800}},
		{"2018", wanttime{year: 2018, month: time.January, day: 1, hour: 0, minute: 0, second: 0, tzname: "UTC"}},
		{"", wanttime{year: 1, month: time.January, day: 1, hour: 0, minute: 0, second: 0, tzname: "UTC"}},
		{"2018-10-27 02:05:00", wanttime{year: 2018, month: time.October, day: 27, hour: 2, minute: 5, second: 0, tzname: "UTC"}},
	}

	for _, test := range tests {

		got, err := TimeOfProgramme(test.input)

		if err != nil {
			t.Errorf("TimeOfProgramme(%q) = %v (%q)", test.input, false, err)
		}

		tzname, tzoffset := got.Zone()

		done := got.Year() == test.want.year && got.Month() == test.want.month && got.Day() == test.want.day &&
			got.Hour() == test.want.hour && got.Minute() == test.want.minute && got.Second() == test.want.second &&
			tzname == test.want.tzname && tzoffset == test.want.tzoffset

		if !done {
			t.Errorf("timeOfProgramme(%q) = %v", test.input, done)
		}
	}
}
