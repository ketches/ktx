/*
Copyright © 2025 Pone Ding <poneding@gmail.com>
*/
package cmd

import (
	"fmt"

	"github.com/poneding/ktx/internal/kubeconfig"
	"github.com/poneding/ktx/internal/output"
	"github.com/poneding/ktx/internal/prompt"
	"github.com/spf13/cobra"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

// renameCmd represents the rename command
var renameCmd = &cobra.Command{
	Use:   "rename",
	Short: "rename context in ~/.kube/config",
	Long:  `rename context in ~/.kube/config`,
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		runRename(args)
	},
	ValidArgsFunction: completeWithContextProfile,
}

func init() {
	rootCmd.AddCommand(renameCmd)
}

func runRename(args []string) {
	config := kubeconfig.Load()

	var dst string
	if len(args) == 0 {
		dst = prompt.ContextSelection("Select context to rename", config)
	} else {
		dst = args[0]
	}

	renameContext(dst, config)
}

func renameContext(dst string, config *clientcmdapi.Config) {
	dstCtx, ok := config.Contexts[dst]
	if !ok {
		output.Fatal("Context <%s> not found.", dst)
	}

	new := prompt.TextInput("Enter a new name")
	if contextNameConflict(new, config) {
		if prompt.YesNo(fmt.Sprintf("Context name <%s> already exists, rename it", new)) != "Yes" {
			return
		}
		new = prompt.TextInput("Enter a new name")
	}

	delete(config.Contexts, dst)
	if config.CurrentContext == dst {
		config.CurrentContext = new
	}
	config.Contexts[new] = dstCtx

	kubeconfig.Save(config)
	output.Done("Context <%s> renamed to <%s>.", dst, new)
}
