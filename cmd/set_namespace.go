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
	completion "github.com/poneding/ktx/internal/completion"
	"github.com/poneding/ktx/internal/kube"
	"github.com/poneding/ktx/internal/output"
	"github.com/spf13/cobra"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

type setNamespaceFlags struct {
	context   string
	namespace string
}

var setNamespaceFlag setNamespaceFlags

// setNamespaceCmd represents the set-namespace command
var setNamespaceCmd = &cobra.Command{
	Use:   "set-namespace",
	Short: "Set context namespace in specified kubeconfig(~/.kube/config by default)",
	Long:  `Set context namespace in specified kubeconfig(~/.kube/config by default)`,
	Args:  cobra.MaximumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		runSetNamespace()
	},
	ValidArgsFunction: completion.None,
}

func init() {
	rootCmd.AddCommand(setNamespaceCmd)

	setNamespaceCmd.Flags().StringVarP(&setNamespaceFlag.context, "context", "c", "", "Context")
	setNamespaceCmd.Flags().StringVarP(&setNamespaceFlag.namespace, "namespace", "n", kube.DefaultNamespace, "Namespace")

	setNamespaceCmd.RegisterFlagCompletionFunc("context", completion.Context)
	setNamespaceCmd.RegisterFlagCompletionFunc("namespace", completion.Namespace)

	setNamespaceCmd.MarkFlagRequired("namespace")
}

func runSetNamespace() {
	config := kube.LoadConfigFromFile(rootFlag.kubeconfig)

	ctxName := setServerFlag.context
	if len(ctxName) == 0 {
		ctxName = config.CurrentContext
	}

	setContextNamespace(config, ctxName, setNamespaceFlag.namespace)
}

func setContextNamespace(config *clientcmdapi.Config, ctx, namespace string) {
	dstCtx, ok := config.Contexts[ctx]
	if !ok {
		output.Fatal("Context <%s> not found.", ctx)
	}

	oldNamespace := dstCtx.Namespace

	if oldNamespace == namespace {
		output.Done("Context <%s> namespace not changed.", ctx)
		return
	}

	config.Contexts[ctx].Namespace = namespace

	kube.SaveConfigToFile(config, rootFlag.kubeconfig)
	output.Done("Context <%s> set namespace from <%s> to <%s>.", ctx, oldNamespace, namespace)
}
