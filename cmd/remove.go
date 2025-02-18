/*
Copyright © 2025 Pone Ding <poneding@gmail.com>
*/
package cmd

import (
	"fmt"

	"github.com/poneding/ktx/internal/completion"
	"github.com/poneding/ktx/internal/kube"
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
	Short:   "Remove context(s) from specified kubeconfig(~/.kube/config by default)",
	Long:    `Remove context(s) from specified kubeconfig(~/.kube/config by default)`,
	Run: func(cmd *cobra.Command, args []string) {
		runRemove(args)
	},
	ValidArgsFunction: completion.ContextArray,
}

func init() {
	rootCmd.AddCommand(removeCmd)
}

func runRemove(args []string) {
	config := kube.LoadConfigFromFile(rootFlag.kubeconfig)

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

	kube.SaveConfigToFile(config, rootFlag.kubeconfig)
	output.Done("Context <%s> removed.", dst)

	// 如果当前没有 context，那么提示用户选择一个 context
	if len(config.Contexts) > 0 && config.CurrentContext == "" {
		if len(config.Contexts) > 0 {
			new := prompt.ContextSelection("Select a context as current", config)
			config.CurrentContext = new
			kube.SaveConfigToFile(config, rootFlag.kubeconfig)
			output.Done("Switched to context <%s>.", new)
		}
	}
}
