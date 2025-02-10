/*
Copyright © 2025 Pone Ding <poneding@gmail.com>
*/
package cmd

import (
	"github.com/poneding/ktx/internal/kubeconfig"
	"github.com/poneding/ktx/internal/output"
	"github.com/poneding/ktx/internal/prompt"
	"github.com/spf13/cobra"
)

// useCmd represents the use command
var useCmd = &cobra.Command{
	Use:   "use",
	Short: "use context in ~/.kube/config",
	Long:  `use context in ~/.kube/config`,
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		runUse(args)
	},
	ValidArgsFunction: completeWithContextProfile,
}

func init() {
	rootCmd.AddCommand(useCmd)
}

func runUse(args []string) {
	config := kubeconfig.Load()

	var dst string
	if len(args) == 0 {
		dst = prompt.ContextSelection("Select context to use", config)
	} else {
		dst = args[0]
	}

	useContext(dst)
}

func useContext(dst string) {
	current := kubeconfig.Load()

	_, ok := current.Contexts[dst]
	if !ok {
		output.Fatal("Context <%s> not found.", dst)
	}

	current.CurrentContext = dst
	kubeconfig.Save(current)
	output.Done("Context <%s> is now in use.", dst)
}
