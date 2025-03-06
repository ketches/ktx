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
	"github.com/poneding/ktx/internal/completion"
	"github.com/poneding/ktx/internal/kube"
	"github.com/spf13/cobra"
)

type generateFlags struct {
	context        string
	namespace      string
	serviceAccount string
	output         string
}

var generateFlag generateFlags

// generateCmd represents the gen command
var generateCmd = &cobra.Command{
	Use:     "generate",
	Aliases: []string{"gen"},
	Short:   "Generate a new context from ServiceAccount",
	Long:    `Generate a new context from ServiceAccount.`,
	Run: func(cmd *cobra.Command, args []string) {
		runGenerate()
	},
	ValidArgsFunction: completion.None,
}

func init() {
	rootCmd.AddCommand(generateCmd)

	generateCmd.Flags().StringVarP(&generateFlag.context, "context", "c", "", "Context")
	generateCmd.Flags().StringVarP(&generateFlag.namespace, "namespace", "n", kube.DefaultNamespace, "Namespace")
	generateCmd.Flags().StringVar(&generateFlag.serviceAccount, "service-account", "", "ServiceAccount")
	generateCmd.Flags().StringVarP(&generateFlag.output, "output", "o", "", "Output kube config file")

	generateCmd.RegisterFlagCompletionFunc("context", completion.Context)
	generateCmd.RegisterFlagCompletionFunc("namespace", completion.Namespace)
	generateCmd.RegisterFlagCompletionFunc("service-account", completion.ServiceAccount)

	generateCmd.MarkFlagRequired("service-account")
}

func runGenerate() {
	generateContext(rootFlag.kubeconfig, generateFlag.context, generateFlag.namespace, generateFlag.serviceAccount)
}

func generateContext(kubeconfig, context, namespace, serviceAccount string) {
	config := kube.GenerateConfigForServiceAccount(kubeconfig, context, namespace, serviceAccount)
	if len(generateFlag.output) == 0 {
		kube.PrintConfig(config)
	} else {
		kube.SaveConfigToFile(config, generateFlag.output)
	}
}
