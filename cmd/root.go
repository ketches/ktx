/*
Copyright © 2025 Pone Ding <poneding@gmail.com>
*/
package cmd

import (
	"fmt"
	"os"
	"slices"

	"github.com/poneding/ktx/internal/kubeconfig"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "ktx",
	Short: "ktx is a tool to manage kubernetes contexts.",
	Long:  `ktx is a tool to manage kubernetes contexts.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func completeWithContextProfile(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) > 0 {
		return nil, cobra.ShellCompDirectiveNoSpace | cobra.ShellCompDirectiveNoFileComp
	}

	current := kubeconfig.Contexts(kubeconfig.Load())
	var completions []string
	for _, context := range current {
		completions = append(completions, fmt.Sprintf("%s\t[%s] %s - %s", context.Name, context.Emoji, context.Namespace, context.Server))
	}

	return completions, cobra.ShellCompDirectiveNoSpace | cobra.ShellCompDirectiveNoFileComp
}

func completeWithContextProfiles(_cmd *cobra.Command, args []string, _toComplete string) ([]string, cobra.ShellCompDirective) {
	current := kubeconfig.Contexts(kubeconfig.Load())

	var completions []string
	for _, context := range current {
		if slices.Contains(args, context.Name) {
			continue
		}
		completions = append(completions, fmt.Sprintf("%s\t[%s] %s - %s", context.Name, context.Emoji, context.Namespace, context.Server))
	}

	return completions, cobra.ShellCompDirectiveNoSpace | cobra.ShellCompDirectiveNoFileComp
}
