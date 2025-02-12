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
	"k8s.io/client-go/tools/clientcmd/api"
)

var (
	removeName string
)

// removeCmd represents the remove command
var removeCmd = &cobra.Command{
	Use:     "remove",
	Aliases: []string{"rm"},
	Short:   "remove context from ~/.kube/config",
	Long:    `remove context from ~/.kube/config`,
	Run: func(cmd *cobra.Command, args []string) {
		runRemove(args)
	},
	ValidArgsFunction: completeWithContextProfiles,
}

func init() {
	rootCmd.AddCommand(removeCmd)
}

func runRemove(args []string) {
	config := kubeconfig.Load()

	dsts := args
	if len(dsts) == 0 {
		dsts = []string{prompt.ContextSelection("Select context to remove", config)}
	}

	for _, dst := range dsts {
		removeContext(config, dst)
	}
}

func removeContext(config *api.Config, dst string) {
	dstCtx, ok := config.Contexts[dst]
	if !ok {
		output.Fatal("Context <%s> not found.", dst)
	}

	if prompt.YesNo(fmt.Sprintf("Are you sure you want to remove context %s", dst)) != "Yes" {
		return
	}

	// 如果没有其他 context 引用这个 cluster 或者 user，那么删除这个 cluster 和 user
	var keepCluster, keepUser bool
	for ctxName, ctx := range config.Contexts {
		if ctxName != dst {
			if ctx.Cluster == dstCtx.Cluster {
				keepCluster = true
			}
			if ctx.AuthInfo == dstCtx.AuthInfo {
				keepUser = true
			}
		}
	}
	if !keepCluster {
		delete(config.Clusters, dstCtx.Cluster)
	}
	if !keepUser {
		delete(config.AuthInfos, dstCtx.AuthInfo)
	}

	delete(config.Contexts, dst)

	// 如果删除的是 current context，那么清空 current context
	if config.CurrentContext == dst {
		config.CurrentContext = ""
	}

	kubeconfig.Save(config)
	output.Done("Context <%s> removed.", dst)

	// 如果当前没有 context，那么提示用户选择一个 context
	if len(config.Contexts) > 0 && config.CurrentContext == "" {
		if len(config.Contexts) > 0 {
			new := prompt.ContextSelection("Remove current context, select another one", config)
			config.CurrentContext = new
			kubeconfig.Save(config)
			output.Done("Context <%s> is now in use.", new)
		}
	}
}
