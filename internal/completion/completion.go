package completion

import (
	"fmt"
	"slices"

	"github.com/poneding/ktx/internal/kube"
	"github.com/spf13/cobra"
)

// None is a shell completion function that does nothing.
func None(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return nil, cobra.ShellCompDirectiveNoFileComp
}

// Context is a shell completion function that completes context names, just one completion.
func Context(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) > 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	current := kube.ListContexts(kube.LoadConfigFromFile(cmd.Flag("kubeconfig").Value.String()))
	var completions []string
	for _, context := range current {
		completions = append(completions, fmt.Sprintf("%s\t[%s] %s - %s", context.Name, context.Emoji, context.Namespace, context.Server))
	}

	return completions, cobra.ShellCompDirectiveNoFileComp
}

// ContextArray is a shell completion function that completes context names, allow multiple completion.
func ContextArray(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	current := kube.ListContexts(kube.LoadConfigFromFile(cmd.Flag("kubeconfig").Value.String()))

	var completions []string
	for _, context := range current {
		if slices.Contains(args, context.Name) {
			continue
		}
		completions = append(completions, fmt.Sprintf("%s\t[%s] %s - %s", context.Name, context.Emoji, context.Namespace, context.Server))
	}

	return completions, cobra.ShellCompDirectiveNoFileComp
}

// Server is a shell completion function that completes server names, just one completion.
func Server(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) > 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	config := kube.LoadConfigFromFile(cmd.Flag("kubeconfig").Value.String())
	ctxName := cmd.Flag("context").Value.String()
	if len(ctxName) == 0 {
		ctxName = config.CurrentContext
	}
	ctx := config.Contexts[ctxName]

	completions := []string{config.Clusters[ctx.Cluster].Server}
	return completions, cobra.ShellCompDirectiveNoFileComp
}

// Namespace is a shell completion function that completes namespace names, just one completion.
func Namespace(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) > 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	kubeClientset := kube.Client(cmd.Flag("kubeconfig").Value.String(), cmd.Flag("context").Value.String())

	return kube.ListNamespaces(kubeClientset), cobra.ShellCompDirectiveNoFileComp
}

// ServiceAccount is a shell completion function that completes service account names, just one completion.
func ServiceAccount(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) > 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	kubeClientset := kube.Client(cmd.Flag("kubeconfig").Value.String(), cmd.Flag("context").Value.String())

	return kube.ListServiceAccounts(kubeClientset, cmd.Flag("namespace").Value.String()), cobra.ShellCompDirectiveNoFileComp
}
