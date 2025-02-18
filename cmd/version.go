/*
Copyright © 2025 Pone Ding <poneding@gmail.com>
*/
package cmd

import (
	"github.com/poneding/ktx/internal/completion"
	"github.com/poneding/ktx/internal/output"
	"github.com/spf13/cobra"
)

const VERSION = "v0.2.1"

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:     "version",
	Aliases: []string{"v"},
	Short:   "Print the version number of ktx",
	Long:    `Print the version number of ktx`,
	Run: func(cmd *cobra.Command, args []string) {
		output.Done("Version: %s", VERSION)
	},
	ValidArgsFunction: completion.None,
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
