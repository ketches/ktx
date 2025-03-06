/*
Copyright © 2025 Pone Ding <poneding@gmail.com>

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

	"github.com/poneding/ktx/internal/completion"
	"github.com/poneding/ktx/internal/kube"
	"github.com/poneding/ktx/internal/output"
	"github.com/poneding/ktx/internal/prompt"
	"github.com/spf13/cobra"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

// renameCmd represents the rename command
var renameCmd = &cobra.Command{
	Use:   "rename",
	Short: "Rename context in specified kubeconfig(~/.kube/config by default)",
	Long:  `Rename context in specified kubeconfig(~/.kube/config by default)`,
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		runRename(args)
	},
	ValidArgsFunction: completion.Context,
}

func init() {
	rootCmd.AddCommand(renameCmd)
}

func runRename(args []string) {
	config := kube.LoadConfigFromFile(rootFlag.kubeconfig)

	var dst string
	if len(args) == 0 {
		dst = prompt.ContextSelection("Select context to rename", config)
	} else {
		dst = args[0]
	}

	renameContext(dst, config)
}

func renameContext(oldCtxName string, config *clientcmdapi.Config) {
	dstCtx, ok := config.Contexts[oldCtxName]
	if !ok {
		output.Fatal("Context <%s> not found.", oldCtxName)
	}

	var (
		oldCluster = dstCtx.Cluster
		oldUser    = dstCtx.AuthInfo
		cluster    = config.Clusters[dstCtx.Cluster]
		user       = config.AuthInfos[dstCtx.AuthInfo]
	)

	newCtxName := prompt.TextInput("Enter a new name", oldCtxName)
	for contextNameConflict(newCtxName, config) {
		if newCtxName == oldCtxName {
			output.Done("Context <%s> not changed.", oldCtxName)
			return
		}

		newCtxName = prompt.TextInput(fmt.Sprintf("Context name <%s> already exists, enter a new name", newCtxName), newCtxName)
		dstCtx.Cluster = "cluster-" + newCtxName
		dstCtx.AuthInfo = "user-" + newCtxName
	}

	if dstCtx.Cluster != oldCluster {
		dstCtx.Cluster = "cluster-" + newCtxName
		config.Clusters[dstCtx.Cluster] = cluster
		delete(config.Clusters, oldCluster)
	}

	if dstCtx.AuthInfo != oldUser {
		dstCtx.AuthInfo = "user-" + newCtxName
		config.AuthInfos[dstCtx.AuthInfo] = user
		delete(config.AuthInfos, oldUser)
	}
	config.Contexts[newCtxName] = dstCtx
	if config.CurrentContext == oldCtxName {
		config.CurrentContext = newCtxName
	}
	delete(config.Contexts, oldCtxName)

	kube.SaveConfigToFile(config, rootFlag.kubeconfig)
	output.Done("Context <%s> renamed to <%s>.", oldCtxName, newCtxName)
}
