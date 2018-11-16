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
