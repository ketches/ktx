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
	"github.com/poneding/ktx/internal/output"
	"github.com/poneding/ktx/internal/prompt"
	"github.com/spf13/cobra"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

// switchCmd represents the switch command
var switchCmd = &cobra.Command{
	Use:     "switch",
	Aliases: []string{"s"},
	Short:   "Switch context in specified kubeconfig(~/.kube/config by default)",
	Long:    `Switch context in specified kubeconfig(~/.kube/config by default)`,
	Args:    cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		runSwitch(args)
	},
	ValidArgsFunction: completion.Context,
}

func init() {
	rootCmd.AddCommand(switchCmd)
}

func runSwitch(args []string) {
	config := kube.LoadConfigFromFile(rootFlag.kubeconfig)

	var dst string
	if len(args) == 0 {
		dst = prompt.ContextSelection("Switch to context", config)
	} else {
		dst = args[0]
	}

	switchContext(config, dst)
}

func switchContext(config *clientcmdapi.Config, dst string) {
	_, ok := config.Contexts[dst]
	if !ok {
		output.Fatal("Context <%s> not found.", dst)
	}

	config.CurrentContext = dst
	kube.SaveConfigToFile(config, rootFlag.kubeconfig)
	output.Done("Switched to context <%s>.", dst)
}
