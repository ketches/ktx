/*
Copyright © 2025 Pone Ding <poneding@gmail.com>
*/
package cmd

import (
	"github.com/poneding/ktx/internal/kubeconfig"
	"github.com/poneding/ktx/internal/output"
	"github.com/poneding/ktx/internal/prompt"
	"github.com/spf13/cobra"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

var (
	exportFile string
)

// exportCmd represents the export command
var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "export context to kubeconfig file from ~/.kube/config",
	Long:  `export context to kubeconfig file from ~/.kube/config`,
	Run: func(cmd *cobra.Command, args []string) {
		runExport(args)
	},
	ValidArgsFunction: completeWithContextProfile,
}

func init() {
	rootCmd.AddCommand(exportCmd)

	exportCmd.Flags().StringVarP(&exportFile, "file", "f", "", "kubeconfig file")

	exportCmd.MarkFlagRequired("file")
}

func runExport(args []string) {
	config := kubeconfig.Load()

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

	// 如果只有一个 context，那么设置为当前 context
	if len(dstConfig.Contexts) == 1 {
		dstConfig.CurrentContext = dsts[0]
	}

	kubeconfig.SaveToFile(dstConfig, exportFile)
	output.Done("Context exported to %s.", exportFile)

	// 如果导出的 kubeconfig 文件没有 current context，那么提示用户选择一个 current context
	if dstConfig.CurrentContext == "" {
		currentCtx := prompt.ContextSelection("Select current context fro exported file", dstConfig)
		dstConfig.CurrentContext = currentCtx
		kubeconfig.SaveToFile(dstConfig, exportFile)
		output.Done("Current context set to %s for exported file %s.", currentCtx, exportFile)
	}
}
