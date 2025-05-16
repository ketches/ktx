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
	"sync"
	"time"

	"github.com/fatih/color"
	"github.com/jedib0t/go-pretty/v6/table"
	completion "github.com/ketches/ktx/internal/completion"
	"github.com/ketches/ktx/internal/kube"
	"github.com/ketches/ktx/internal/output"
	"github.com/ketches/ktx/internal/prompt"
	"github.com/ketches/ktx/internal/types"
	"github.com/spf13/cobra"
)

type listFlags struct {
	clusterInfo bool
}

var listFlag listFlags

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

	listCmd.Flags().BoolVar(&listFlag.clusterInfo, "cluster-info", false, "Show cluster info eg. status, version, and more.")
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

var tableStyle = table.Style{
	Name:    "KrbTableStyle",
	Box:     table.StyleBoxDefault,
	Color:   table.ColorOptionsDefault,
	Format:  table.FormatOptionsDefault,
	HTML:    table.DefaultHTMLOptions,
	Options: table.OptionsNoBordersAndSeparators,
	Size:    table.SizeOptionsDefault,
	Title:   table.TitleOptionsDefault,
}

func listContexts(ctxs []*types.ContextProfile) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	row := table.Row{"", "name", "namespace", "server"}
	if listFlag.clusterInfo {
		row = append(row, "status", "version")

		var wg sync.WaitGroup
		for i, ctx := range ctxs {
			wg.Add(1)
			go func(i int, ctx *types.ContextProfile) {
				defer wg.Done()
				var (
					clusterStatus  = types.ClusterStatusUnavailable
					clusterVersion = "-"
				)
				timer := time.NewTimer(time.Second * 5)
				defer timer.Stop()
				select {
				case <-timer.C:
					clusterStatus = types.ClusterStatusTimeout
				default:
					dc, _ := kube.DiscoveryClient(rootFlag.kubeconfig, ctx.Name)
					if dc != nil {
						cv, _ := kube.Version(dc)
						if cv != "" {
							clusterStatus = types.ClusterStatusAvailable
							clusterVersion = cv
						}
					}
				}

				ctx.ClusterStatus = clusterStatus
				ctx.ClusterVersion = clusterVersion
			}(i, ctx)
		}
		wg.Wait()
	}

	t.AppendHeader(row)

	for _, ctx := range ctxs {
		appendRow(t, ctx)
	}
	t.SetStyle(tableStyle)
	t.Render()
}

func appendRow(t table.Writer, ctx *types.ContextProfile) {
	if ctx.Current {
		ctx.Name = color.CyanString(ctx.Name)
		ctx.Namespace = color.CyanString(ctx.Namespace)
		ctx.Server = color.CyanString(ctx.Server)
	}
	row := table.Row{ctx.Emoji, ctx.Name, ctx.Namespace, ctx.Server}
	if listFlag.clusterInfo {
		row = append(row, string(ctx.ClusterStatus.ColorString()), color.CyanString(ctx.ClusterVersion))
	}
	t.AppendRow(row, table.RowConfig{
		AutoMerge: true,
	})
}
