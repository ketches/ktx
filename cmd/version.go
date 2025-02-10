/*
Copyright © 2025 Pone Ding <poneding@gmail.com>
*/
package cmd

import (
	"github.com/poneding/ktx/internal/output"
	"github.com/spf13/cobra"
)

const VERSION = "0.1.1"

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of ktx",
	Long:  `Print the version number of ktx`,
	Run: func(cmd *cobra.Command, args []string) {
		output.Done("Version: v%s", VERSION)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
