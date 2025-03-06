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

	completion "github.com/poneding/ktx/internal/completion"
	"github.com/poneding/ktx/internal/kube"
	"github.com/poneding/ktx/internal/output"
	"github.com/poneding/ktx/internal/prompt"
	"github.com/poneding/ktx/internal/util"
	"github.com/spf13/cobra"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

var (
	addFile string
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add context from kubeconfig file to ~/.kube/config",
	Long:  `Add context from kubeconfig file to ~/.kube/config`,
	Run: func(cmd *cobra.Command, args []string) {
		runAdd()
	},
	ValidArgsFunction: completion.None,
}

func init() {
	rootCmd.AddCommand(addCmd)

	addCmd.Flags().StringVarP(&addFile, "file", "f", "", "kubeconfig file")

	addCmd.MarkFlagRequired("file")
}

func runAdd() {
	if !util.IsFileExist(addFile) {
		output.Fatal("File %s not found.", addFile)
	}

	config := kube.LoadConfigFromFile(rootFlag.kubeconfig)
	kube.StandardizeConfig(config)

	new := kube.LoadConfigFromFile(addFile)
	kube.StandardizeConfig(new)

	merge(config, new)
}

func merge(config, new *clientcmdapi.Config) {
	for newCtxName, newCtx := range new.Contexts {
		newCluster, ok := new.Clusters[newCtx.Cluster]
		if !ok {
			output.Note("Cluster not found for context <%s> in file %s, skipped.", newCtxName, addFile)
			continue
		}

		newUser, ok := new.AuthInfos[newCtx.AuthInfo]
		if !ok {
			output.Note("User not found for context <%s> in file %s, skipped.", newCtxName, addFile)
			continue
		}

		// 如果 context 名称已经存在，要求用户输入新的 context 名称
		var quitWithConflict bool
		for contextNameConflict(newCtxName, config) {
			if prompt.YesNo(fmt.Sprintf("Context name <%s> already exists, rename it", newCtxName)) {
				newCtxName = prompt.TextInput("Enter a new context name", newCtxName)
				newCtx.Cluster = "cluster-" + newCtxName
				newCtx.AuthInfo = "user-" + newCtxName
			} else {
				quitWithConflict = true
				break
			}
		}

		if quitWithConflict {
			continue
		}

		mf := &mergeFrom{
			contextName: newCtxName,
			clusterName: newCtx.Cluster,
			userName:    newCtx.AuthInfo,
			context:     newCtx,
			cluster:     newCluster,
			user:        newUser,
		}

		handleMerge(config, mf)
	}

	kube.SaveConfigToFile(config, rootFlag.kubeconfig)

	// 如果当前没有 context，那么提示用户选择一个 context
	if config.CurrentContext == "" && len(config.Contexts) > 0 {
		if len(config.Contexts) == 1 {
			for ctxName := range config.Contexts {
				config.CurrentContext = ctxName
				break
			}
		} else {
			config.CurrentContext = prompt.ContextSelection("Select a context to as current", config)
		}
		kube.SaveConfigToFile(config, rootFlag.kubeconfig)
	}
}

type mergeFrom struct {
	contextName, clusterName, userName string
	context                            *clientcmdapi.Context
	cluster                            *clientcmdapi.Cluster
	user                               *clientcmdapi.AuthInfo
}

func handleMerge(config *clientcmdapi.Config, mf *mergeFrom) {
	for contextNameConflict(mf.contextName, config) {
		mf.contextName = prompt.TextInput(fmt.Sprintf("Context name <%s> already exists, enter a new name", mf.contextName), mf.contextName)
		mf.clusterName = "cluster-" + mf.contextName
		mf.userName = "user-" + mf.contextName
	}

	config.Clusters[mf.clusterName] = mf.cluster
	config.AuthInfos[mf.userName] = mf.user
	config.Contexts[mf.contextName] = mf.context

	output.Done("Context <%s> added.", mf.contextName)
}

func contextNameConflict(name string, config *clientcmdapi.Config) bool {
	_, ok := config.Contexts[name]
	return ok
}
