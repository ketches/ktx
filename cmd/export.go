/*
Copyright © 2025 Pone Ding <poneding@gmail.com>
*/
package cmd

import (
	completion "github.com/poneding/ktx/internal/completion"
	"github.com/poneding/ktx/internal/kube"
	kubeconfig "github.com/poneding/ktx/internal/kube"
	"github.com/poneding/ktx/internal/output"
	"github.com/poneding/ktx/internal/prompt"
	"github.com/spf13/cobra"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

type exportFlags struct {
	output string
}

var exportFlag exportFlags

// exportCmd represents the export command
var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export context(s) from specified kubeconfig(~/.kube/config by default)",
	Long:  `Export context(s) from specified kubeconfig(~/.kube/config by default)`,
	Run: func(cmd *cobra.Command, args []string) {
		runExport(args)
	},
	ValidArgsFunction: completion.ContextArray,
}

func init() {
	rootCmd.AddCommand(exportCmd)

	exportCmd.Flags().StringVarP(&exportFlag.output, "output", "o", "", "Output kube config file")
}

func runExport(args []string) {
	config := kube.LoadConfigFromFile(rootFlag.kubeconfig)

	dsts := args
	if len(dsts) == 0 {
		dsts = []string{prompt.ContextSelection("Select context to export", config)}
	}

	exportContext(config, dsts)
}

func exportContext(config *clientcmdapi.Config, dsts []string) {
	dstConfig := clientcmdapi.NewConfig()
	for _, dst := range dsts {
		dstCtx, ok := config.Contexts[dst]
		if !ok {
			output.Fatal("Context <%s> not found.", dst)
		}

		dstCluster, ok := config.Clusters[dstCtx.Cluster]
		if !ok {
			output.Fatal("Cluster not found for context <%s>.", dstCtx.Cluster, dst)
		}

		dstUser, ok := config.AuthInfos[dstCtx.AuthInfo]
		if !ok {
			output.Fatal("User not found for context <%s>.", dstCtx.AuthInfo, dst)
		}
		dstConfig.Contexts[dst] = dstCtx
		dstConfig.Clusters[dstCtx.Cluster] = dstCluster
		dstConfig.AuthInfos[dstCtx.AuthInfo] = dstUser
	}

	// 设置当前上下文，默认第一个
	if len(dstConfig.Contexts) > 0 {
		dstConfig.CurrentContext = dsts[0]
	}

	if genFlag.output == "" {
		kubeconfig.PrintConfig(dstConfig)
	} else {
		kubeconfig.SaveConfigToFile(config, genFlag.output)
		output.Done("Context exported to %s.", genFlag.output)
	}
}
