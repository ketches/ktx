/*
Copyright © 2025 Pone Ding <poneding@gmail.com>
*/
package cmd

import (
	"os"

	"github.com/fatih/color"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/poneding/ktx/internal/kubeconfig"
	"github.com/poneding/ktx/internal/output"
	"github.com/poneding/ktx/internal/prompt"
	"github.com/poneding/ktx/internal/types"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "list contexts in ~/.kube/config",
	Long:    `list contexts in ~/.kube/config`,
	Run: func(cmd *cobra.Command, args []string) {
		runList()
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}

func runList() {
	config := kubeconfig.Load()
	ctxs := kubeconfig.Contexts(config)

	if len(ctxs) == 0 {
		output.Note("No context found.")
		return
	}

	listContexts(ctxs)

	// 如果当前没有 context，那么提示用户选择一个 context
	if config.CurrentContext == "" {
		config.CurrentContext = prompt.ContextSelection("Select context to use", config)
		kubeconfig.Save(config)
	}
}

func listContexts(ctxs []*types.ContextProfile) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"", "name", "namespace", "cluster_name", "user_name", "server"})

	for _, ctx := range ctxs {
		appendRow(t, ctx)
	}
	t.SetStyle(table.StyleRounded)
	t.Render()
}

func appendRow(t table.Writer, ctx *types.ContextProfile) {
	var flag string
	if ctx.Current {
		ctx.Name = color.CyanString(ctx.Name)
		ctx.Namespace = color.CyanString(ctx.Namespace)
		ctx.Cluster = color.CyanString(ctx.Cluster)
		ctx.User = color.CyanString(ctx.User)
		ctx.Server = color.CyanString(ctx.Server)
		flag = color.CyanString("✔")
	}
	t.AppendRow(table.Row{flag, ctx.Name, ctx.Namespace, ctx.Cluster, ctx.User, ctx.Server}, table.RowConfig{
		AutoMerge: true,
	})
}
