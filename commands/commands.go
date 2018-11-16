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

package commands

import (
	"github.com/spf13/cobra"
)

var rootCommand = &cobra.Command{
	Use:   "tvguide",
	Short: "IPTV guide viewer",
	Long:  "Another yet (maybe) IPTV guide viewer",
	Run: func(cmd *cobra.Command, args []string) {

	},
}

// PlaylistPath - path or URL of the playlist
var PlaylistPath string

func init() {

	cmdView.Flags().StringVarP(&PlaylistPath, "playlist", "p", "", "path or URL of the playlist (required)")
	cmdView.MarkFlagRequired("playlist")

	rootCommand.AddCommand(cmdView, cmdVersion)
}

// Execute is a enter point into application commands
func Execute() error {
	err := rootCommand.Execute()
	return err
}
