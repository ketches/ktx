package prompt

import (
	"fmt"
	"os"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/poneding/ktx/internal/kubeconfig"
	"github.com/poneding/ktx/internal/output"
	"github.com/poneding/ktx/internal/types"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

func YesNo(label string) string {
	templates := &promptui.SelectTemplates{
		Label:    promptui.Styler(promptui.FGYellow)("# {{ . }}?"),
		Active:   promptui.Styler(promptui.FGCyan, promptui.FGUnderline)("➤ {{ . }}"),
		Inactive: promptui.Styler(promptui.FGFaint)("  {{ . }}"),
	}
	prompt := promptui.Select{
		Label:        label,
		Items:        []string{"No", "Yes"},
		Templates:    templates,
		Size:         4,
		HideSelected: true,
	}
	_, obj, err := prompt.Run()
	if err != nil {
		output.Fatal("Prompt failed %v", err)
	}

	return obj
}

func TextInput(label string) string {
	prompt := promptui.Prompt{
		Label: promptui.Styler(promptui.FGYellow)(label),
		Validate: func(input string) error {
			if len(strings.TrimSpace(input)) == 0 {
				return fmt.Errorf("Please input a valid value")
			}
			return nil
		},
		Templates: &promptui.PromptTemplates{
			Prompt:          promptui.Styler(promptui.FGCyan)("➤ {{ . }} "),
			ValidationError: promptui.Styler(promptui.FGRed)("✗ {{ . }}"),
		},
		HideEntered: true,
	}
	result, err := prompt.Run()
	if err != nil {
		output.Fatal("Prompt failed %v", err)
	}
	result = strings.TrimSpace(result)

	return result
}

func ContextSelection(label string, config *clientcmdapi.Config) string {
	ctxs := kubeconfig.Contexts(config)
	ctxs = append(ctxs, &types.ContextProfile{
		Name:  "Exit",
		Emoji: "✗",
	})
	cursorPos := 0
	for i, ctx := range ctxs {
		if ctx.Current {
			cursorPos = i
		}
	}

	templates := &promptui.SelectTemplates{
		Label:    promptui.Styler(promptui.FGYellow)("# {{ . }}:"),
		Active:   promptui.Styler(promptui.FGCyan, promptui.FGUnderline)("➤ {{ .Emoji }} {{ .Name }}"),
		Inactive: promptui.Styler(promptui.FGFaint)("  {{ .Emoji }} {{ .Name }}"),
		Details: `{{if .Cluster}}
---------- Context ----------
{{ "Name:" | faint }}	{{ .Name }}
{{ "Namespace:" | faint }}	{{ .Namespace }}
{{ "Cluster:" | faint }}	{{ .Cluster }}
{{ "User:" | faint }}	{{ .User }}
{{ "Server:" | faint }}	{{ .Server }}{{end}}`,
	}

	prompt := promptui.Select{
		Label: label,
		Items: ctxs,
		Searcher: func(input string, index int) bool {
			if index < 0 || index >= len(ctxs) {
				return false
			}

			current := ctxs[index]
			if strings.Contains(strings.ToLower(current.Name), strings.ToLower(input)) ||
				strings.Contains(strings.ToLower(current.Server), strings.ToLower(input)) {
				return true
			}

			return false
		},
		HideSelected: true,
		CursorPos:    cursorPos,
		Templates:    templates,
	}

	index, _, err := prompt.Run()
	if err != nil {
		output.Fatal("Prompt failed %v", err)
	}

	if index == len(ctxs)-1 {
		os.Exit(0)
	}

	return ctxs[index].Name
}
