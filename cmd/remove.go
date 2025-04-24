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
	"fmt"

	"github.com/ketches/ktx/internal/completion"
	"github.com/ketches/ktx/internal/kube"
	"github.com/ketches/ktx/internal/output"
	"github.com/ketches/ktx/internal/prompt"
	"github.com/spf13/cobra"
	"k8s.io/client-go/tools/clientcmd/api"
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

	if !prompt.YesNo(fmt.Sprintf("Are you sure you want to remove context %s", dst)) {
		return
	}

	delete(config.Clusters, dstCtx.Cluster)
	delete(config.AuthInfos, dstCtx.AuthInfo)
	delete(config.Contexts, dst)

	// 如果删除的是 current context，那么清空 current context
	if config.CurrentContext == dst {
		config.CurrentContext = ""
	}

	kube.SaveConfigToFile(config, rootFlag.kubeconfig)
	output.Done("Context <%s> removed.", dst)

	// 如果当前没有 context，那么提示用户选择一个 context
	if len(config.CurrentContext) == 0 && len(config.Contexts) > 0 {
		if len(config.Contexts) > 0 {
			new := prompt.ContextSelection("Select a context as current", config)
			config.CurrentContext = new
			kube.SaveConfigToFile(config, rootFlag.kubeconfig)
			output.Done("Switched to context <%s>.", new)
		}
	}
}
