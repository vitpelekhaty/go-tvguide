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

import "time"

// Programme contains info about tv programme
type Programme struct {
	PID   int
	Start time.Time
	Stop  time.Time
	Title string
}

// StartHour returns the hour of the TV program start
func (gi *Programme) StartHour() int {
	return gi.Start.Hour()
}

// StartMinute returns the minute of the TV program start
func (gi *Programme) StartMinute() int {
	return gi.Start.Minute()
}

// StopHour returns the end hour of the TV program
func (gi *Programme) StopHour() int {
	return gi.Stop.Hour()
}

// StopMinute returns the end minute of the TV program
func (gi *Programme) StopMinute() int {
	return gi.Stop.Minute()
}
