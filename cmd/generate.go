/*
Copyright © 2025 Pone Ding <poneding@gmail.com>
*/
package cmd

import (
	"github.com/poneding/ktx/internal/completion"
	kubeconfig "github.com/poneding/ktx/internal/kube"
	"github.com/spf13/cobra"
)

type genFlags struct {
	namespace      string
	serviceAccount string
	output         string
}

var genFlag genFlags

// genCmd represents the gen command
var genCmd = &cobra.Command{
	Use:     "generate",
	Aliases: []string{"gen"},
	Short:   "Generate a new context from ServiceAccount",
	Long:    `Generate a new context from ServiceAccount.`,
	Run: func(cmd *cobra.Command, args []string) {
		runGen()
	},
	ValidArgsFunction: completion.None,
}

func init() {
	rootCmd.AddCommand(genCmd)

	genCmd.Flags().StringVarP(&genFlag.namespace, "namespace", "n", "default", "Namespace")
	genCmd.Flags().StringVar(&genFlag.serviceAccount, "service-account", "", "ServiceAccount")
	genCmd.Flags().StringVarP(&genFlag.output, "output", "o", "", "Output kube config file")

	genCmd.RegisterFlagCompletionFunc("namespace", completion.Namespace)
	genCmd.RegisterFlagCompletionFunc("service-account", completion.ServiceAccount)

	genCmd.MarkFlagRequired("service-account")
}

func runGen() {
	config := kubeconfig.GenerateConfigForServiceAccount(genFlag.serviceAccount, genFlag.namespace)
	if genFlag.output == "" {
		kubeconfig.PrintConfig(config)
	} else {
		kubeconfig.SaveConfigToFile(config, genFlag.output)
	}
}
