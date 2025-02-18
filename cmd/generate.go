/*
Copyright © 2025 Pone Ding <poneding@gmail.com>
*/
package cmd

import (
	"github.com/poneding/ktx/internal/completion"
	kubeconfig "github.com/poneding/ktx/internal/kube"
	"github.com/spf13/cobra"
)

type generateFlags struct {
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

	generateCmd.Flags().StringVarP(&generateFlag.namespace, "namespace", "n", "default", "Namespace")
	generateCmd.Flags().StringVar(&generateFlag.serviceAccount, "service-account", "", "ServiceAccount")
	generateCmd.Flags().StringVarP(&generateFlag.output, "output", "o", "", "Output kube config file")

	generateCmd.RegisterFlagCompletionFunc("namespace", completion.Namespace)
	generateCmd.RegisterFlagCompletionFunc("service-account", completion.ServiceAccount)

	generateCmd.MarkFlagRequired("service-account")
}

func runGenerate() {
	config := kubeconfig.GenerateConfigForServiceAccount(generateFlag.serviceAccount, generateFlag.namespace)
	if generateFlag.output == "" {
		kubeconfig.PrintConfig(config)
	} else {
		kubeconfig.SaveConfigToFile(config, generateFlag.output)
	}
}
