/*
Copyright © 2025 The Ketches Authors.

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
	"os"

	"github.com/fatih/color"
	"github.com/jedib0t/go-pretty/v6/table"
	completion "github.com/ketches/ktx/internal/completion"
	"github.com/ketches/ktx/internal/kube"
	"github.com/ketches/ktx/internal/output"
	"github.com/ketches/ktx/internal/prompt"
	"github.com/ketches/ktx/internal/types"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List contexts in ~/.kube/config",
	Long:    `List contexts in ~/.kube/config`,
	Run: func(cmd *cobra.Command, args []string) {
		runList()
	},
	ValidArgsFunction: completion.None,
}

func init() {
	rootCmd.AddCommand(listCmd)
}

func runList() {
	config := kube.LoadConfigFromFile(rootFlag.kubeconfig)
	ctxs := kube.ListContexts(config)

	if len(ctxs) == 0 {
		output.Note("No context found.")
		return
	}

	listContexts(ctxs)

	// 如果当前没有 context，那么提示用户选择一个 context
	if len(config.CurrentContext) == 0 {
		config.CurrentContext = prompt.ContextSelection("Select a context as current", config)
		kube.SaveConfigToFile(config, rootFlag.kubeconfig)
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
	if ctx.Current {
		ctx.Name = color.CyanString(ctx.Name)
		ctx.Namespace = color.CyanString(ctx.Namespace)
		ctx.Cluster = color.CyanString(ctx.Cluster)
		ctx.User = color.CyanString(ctx.User)
		ctx.Server = color.CyanString(ctx.Server)
	}
	t.AppendRow(table.Row{ctx.Emoji, ctx.Name, ctx.Namespace, ctx.Cluster, ctx.User, ctx.Server}, table.RowConfig{
		AutoMerge: true,
	})
}
