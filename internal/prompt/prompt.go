package prompt

import (
	"fmt"
	"os"
	"strings"

	"github.com/ketches/ktx/internal/kube"
	"github.com/ketches/ktx/internal/output"
	"github.com/ketches/ktx/internal/types"
	"github.com/manifoldco/promptui"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

// YesNo prompts the user to select Yes or No
func YesNo(label string) bool {
	templates := &promptui.SelectTemplates{
		Label:    promptui.Styler(promptui.FGYellow)("❖ {{ . }}?"),
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

	return obj == "Yes"
}

// TextInput prompts the user to input a value
func TextInput(label, def string) string {
	prompt := promptui.Prompt{
		Label: promptui.Styler(promptui.FGYellow)(label),
		Validate: func(input string) error {
			if len(strings.TrimSpace(input)) == 0 {
				return fmt.Errorf("please input a valid value")
			}
			return nil
		},
		Default: def,
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

// ContextSelection prompts the user to select a context
func ContextSelection(label string, config *clientcmdapi.Config) string {
	ctxs := kube.ListContexts(config)
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
		Label:    promptui.Styler(promptui.FGYellow)("❖ {{ . }}:"),
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
