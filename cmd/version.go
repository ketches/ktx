/*
Copyright © 2025 The Ketches Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cmd

import (
	"github.com/ketches/ktx/internal/completion"
	"github.com/ketches/ktx/internal/output"
	"github.com/spf13/cobra"
)

const VERSION = "v0.3.1"

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
