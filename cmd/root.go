/*
Copyright © 2025 Pone Ding <poneding@gmail.com>
*/
package cmd

import (
	"os"

	"github.com/poneding/ktx/internal/kube"
	"github.com/spf13/cobra"
)

type rootFlags struct {
	kubeconfig string
}

var rootFlag rootFlags

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "ktx",
	Short: "ktx is a tool to manage kubernetes contexts.",
	Long:  `ktx is a tool to manage kubernetes contexts.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		// 默认运行 switch 子命令
		runSwitch(args)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	// 如果 kubeconfig 为默认值，则检查或初始化 kubeconfig
	if rootFlag.kubeconfig == kube.DefaultConfigFile {
		kube.CheckOrInitConfig()
	}

	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&rootFlag.kubeconfig, "kubeconfig", kube.DefaultConfigFile, "kubeconfig file")
}
