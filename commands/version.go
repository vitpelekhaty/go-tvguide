package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

var cmdVersion = &cobra.Command{
	Use:   "version",
	Short: "Print version number of tvguide",
	Long:  "Most of applications have version number. TVGuide is one of them",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("TVGuide v2018.10.0.1a")
	},
}
