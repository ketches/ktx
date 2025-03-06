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

type setServerFlags struct {
	context string
	server  string
}

var setServerFlag setServerFlags

// setServerCmd represents the set-server command
var setServerCmd = &cobra.Command{
	Use:   "set-server",
	Short: "Set context server host in specified kubeconfig(~/.kube/config by default)",
	Long:  `Set context server host in specified kubeconfig(~/.kube/config by default)`,
	Args:  cobra.MaximumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		runSetServer()
	},
	ValidArgsFunction: completion.None,
}

func init() {
	rootCmd.AddCommand(setServerCmd)

	setServerCmd.Flags().StringVarP(&setServerFlag.context, "context", "c", "", "Context")
	setServerCmd.Flags().StringVarP(&setServerFlag.server, "server", "s", "", "Server host")

	setServerCmd.RegisterFlagCompletionFunc("context", completion.Context)
	setServerCmd.RegisterFlagCompletionFunc("server", completion.Server)
	setServerCmd.MarkFlagRequired("server")
}

func runSetServer() {
	config := kube.LoadConfigFromFile(rootFlag.kubeconfig)

	ctxName := setServerFlag.context
	if len(ctxName) == 0 {
		ctxName = config.CurrentContext
	}

	setContextServer(config, ctxName, setServerFlag.server)
}

func setContextServer(config *clientcmdapi.Config, ctx, server string) {
	dstCtx, ok := config.Contexts[ctx]
	if !ok {
		output.Fatal("Context <%s> not found.", ctx)
	}

	oldServer := config.Clusters[dstCtx.Cluster].Server

	if oldServer == server {
		output.Done("Context <%s> server not changed.", ctx)
		return
	}

	config.Clusters[dstCtx.Cluster].Server = server

	kube.SaveConfigToFile(config, rootFlag.kubeconfig)
	output.Done("Context <%s> set server from <%s> to <%s>.", ctx, oldServer, server)
}
