// IPTV guide viewer
//
// Copyright 2018 Vitaly Pelekhaty
// Based on https://github.com/antavelos/terminews/blob/master/ui.go
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

package strutils

import (
	"bufio"
	"errors"
	"strings"
	"unicode/utf8"
)

// StringList maintains a list of strings
type StringList struct {
	RightMargin int
	items       []*string
}

// Clear deletes all the strings from the list
func (sl *StringList) Clear() {
	sl.items = make([]*string, 0)
}

// Count returns the number of strings in the list
func (sl *StringList) Count() int {

	return len(sl.items)
}

// Add adds a new string to the list
func (sl *StringList) Add(s string) int {

	sl.items = append(sl.items, &s)
	return len(sl.items)
}

// SetText ...
func (sl *StringList) SetText(text string) (err error) {

	var lines []*string

	lines = make([]*string, 0)
	lines, err = sl.splitText(text, lines)

	if err != nil {
		return
	}

	sl.items = make([]*string, len(lines))

	for index, line := range lines {
		sl.items[index] = line
	}

	return
}

func (sl *StringList) splitText(text string, lines []*string) (output []*string, err error) {

	var (
		reader  *strings.Reader
		scanner *bufio.Scanner
	)

	output = make([]*string, len(lines))

	for idx, line := range lines {
		output[idx] = line
	}

	reader = strings.NewReader(text)
	scanner = bufio.NewScanner(reader)

	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {

		var line string
		line = scanner.Text()

		if sl.RightMargin > 0 {

			output, err = sl.splitLine(line, output)

			if err != nil {
				return
			}
		} else {

			var item string
			item = line

			output = append(output, &item)
		}
	}

	return
}

func (sl *StringList) splitLine(line string, lines []*string) (output []*string, err error) {

	var (
		reader  *strings.Reader
		scanner *bufio.Scanner
		word    string
		l       string
		margin  int
		rc      int
	)

	margin = sl.RightMargin

	output = make([]*string, len(lines))

	for idx, line := range lines {
		output[idx] = line
	}

	reader = strings.NewReader(line)
	scanner = bufio.NewScanner(reader)

	scanner.Split(bufio.ScanWords)

	for scanner.Scan() {

		word = scanner.Text()
		rc = utf8.RuneCountInString(word) + utf8.RuneCountInString(l) + 1

		if rc > margin {

			var s string

			s = l
			output = append(output, &s)

			l = ""
		}

		if len(l) == 0 {
			l = word
		} else {
			l = strings.Join([]string{l, word}, " ")
		}
	}

	if len(l) > 0 {

		var s string
		s = l

		output = append(output, &s)
	}

	return
}

// Text lists the strings in the StringList as a single string
func (sl *StringList) Text() string {

	buffer := make([]string, len(sl.items))

	for index, item := range sl.items {
		buffer[index] = *item
	}

	return strings.Join(buffer, "\n")
}

// IndexOf returns the position of a string in the list
func (sl *StringList) IndexOf(s string) int {

	for index, item := range sl.items {

		if *item == s {
			return index
		}
	}

	return -1
}

// Delete removes the string specified by the index parameter
func (sl *StringList) Delete(index int) error {

	if index < 0 || index > (len(sl.items)-1) {
		return errors.New("Index out of bounds")
	}

	before := sl.items[:index-1]
	after := sl.items[index+1:]

	var lb int
	lb = len(before)

	sl.items = make([]*string, lb+len(after))

	for idx, item := range before {
		sl.items[idx] = item
	}

	for idx, item := range after {
		sl.items[lb+idx] = item
	}

	return nil
}

// Item returns the item of the string list
func (sl *StringList) Item(index int) (string, error) {

	if index < 0 || index > (len(sl.items)-1) {
		return "", errors.New("Index out of bounds")
	}

	return *sl.items[index], nil
}
