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

package strutils

import "testing"

type StringListItemValue struct {
	Index int
	Value string
}

type WantStringList struct {
	Count  int
	Values []StringListItemValue
}

type StringListTest struct {
	TestValue   string
	RightMargin int
	Want        WantStringList
}

var stringlisttests = []StringListTest{
	{
		TestValue: "string 1\nstring 2\nstring 3", RightMargin: 0,
		Want: WantStringList{Count: 3,
			Values: []StringListItemValue{
				StringListItemValue{Index: 0, Value: "string 1"},
				StringListItemValue{Index: 1, Value: "string 2"},
				StringListItemValue{Index: 2, Value: "string 3"},
			},
		},
	},
	{
		TestValue: "lorem ipsum lorem ipsum lorem ipsum lorem ipsum", RightMargin: 11,
		Want: WantStringList{Count: 4,
			Values: []StringListItemValue{
				StringListItemValue{Index: 0, Value: "lorem ipsum"},
				StringListItemValue{Index: 1, Value: "lorem ipsum"},
				StringListItemValue{Index: 2, Value: "lorem ipsum"},
				StringListItemValue{Index: 3, Value: "lorem ipsum"},
			},
		},
	},
	{
		TestValue: "lorem ipsum lorem ipsum lorem ipsum lorem ipsum", RightMargin: 15,
		Want: WantStringList{Count: 4,
			Values: []StringListItemValue{
				StringListItemValue{Index: 0, Value: "lorem ipsum"},
				StringListItemValue{Index: 1, Value: "lorem ipsum"},
				StringListItemValue{Index: 2, Value: "lorem ipsum"},
				StringListItemValue{Index: 3, Value: "lorem ipsum"},
			},
		},
	},
	{
		TestValue: "lorem ipsum lorem ipsum\nlorem ipsum lorem ipsum", RightMargin: 15,
		Want: WantStringList{Count: 4,
			Values: []StringListItemValue{
				StringListItemValue{Index: 0, Value: "lorem ipsum"},
				StringListItemValue{Index: 1, Value: "lorem ipsum"},
				StringListItemValue{Index: 2, Value: "lorem ipsum"},
				StringListItemValue{Index: 3, Value: "lorem ipsum"},
			},
		},
	},
	{
		TestValue: "lorem ipsum lorem\nipsum lorem ipsum lorem ipsum", RightMargin: 20,
		Want: WantStringList{Count: 3,
			Values: []StringListItemValue{
				StringListItemValue{Index: 0, Value: "lorem ipsum lorem"},
				StringListItemValue{Index: 1, Value: "ipsum lorem ipsum"},
				StringListItemValue{Index: 2, Value: "lorem ipsum"},
			},
		},
	},
	{
		TestValue: "lorem ipsum lorem\nipsum lorem ipsum lorem ipsum", RightMargin: 24,
		Want: WantStringList{Count: 3,
			Values: []StringListItemValue{
				StringListItemValue{Index: 0, Value: "lorem ipsum lorem"},
				StringListItemValue{Index: 1, Value: "ipsum lorem ipsum lorem"},
				StringListItemValue{Index: 2, Value: "ipsum"},
			},
		},
	},
}

func TestStringListSetText(t *testing.T) {

	for testIndex, test := range stringlisttests {

		var (
			input     string
			margin    int
			wantCount int
			count     int
		)

		input = test.TestValue
		margin = test.RightMargin
		wantCount = test.Want.Count

		sl := StringList{RightMargin: margin}

		err := sl.SetText(input)

		if err != nil {
			t.Errorf("test %v:\tStringList.SetText() - %v (%q)\n", testIndex+1, false, err)
		}

		count = sl.Count()

		if count != wantCount {
			t.Errorf("test %v:\tStringList.Count() - %v (current %v, must %v)\n", testIndex+1, false, count, wantCount)
		}

		for _, val := range test.Want.Values {

			var (
				wantIndex int
				wantValue string
				item      string
			)

			wantIndex = val.Index
			wantValue = val.Value

			item, err = sl.Item(wantIndex)

			if err != nil {
				t.Errorf("test %v:\tStringList.Item(%v) - %v (%q)\n", testIndex+1, wantIndex, false, err)
			}

			if item != wantValue {
				t.Errorf("test %v:\tStringList.Item(%v) - %v (current %s, must %s)\n", testIndex+1, wantIndex, false, item, wantValue)
			}
		}
	}
}
