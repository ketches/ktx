package kubeconfig

import (
	"os"
	"sort"

	"github.com/poneding/ktx/internal/output"
	"github.com/poneding/ktx/internal/types"
	"github.com/poneding/ktx/internal/util"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
	"k8s.io/client-go/util/homedir"
)

func Dir() string {
	return homedir.HomeDir() + "/.kube"
}

func File() string {
	return homedir.HomeDir() + "/.kube/config"
}

func Load() *clientcmdapi.Config {
	CheckOrInit()

	return LoadFromFile(File())
}

func LoadFromFile(file string) *clientcmdapi.Config {
	config, err := clientcmd.LoadFromFile(file)
	if err != nil {
		output.Fatal("Failed to load kubeconfig from file: %s", err)
	}

	return config
}

func Save(config *clientcmdapi.Config) {
	SaveToFile(config, File())
}

func SaveToFile(config *clientcmdapi.Config, file string) {
	if err := clientcmd.WriteToFile(*config, file); err != nil {
		output.Fatal("Failed to save kubeconfig to file: %s", err)
	}
}

func CheckOrInit() {
	kubeconfigDir := Dir()
	kubeconfigFile := File()

	if _, err := os.Stat(kubeconfigFile); os.IsNotExist(err) {
		if err := os.MkdirAll(kubeconfigDir, 0700); err != nil {
			output.Fatal("Failed to create %s directory: %s", kubeconfigDir, err)
		}

		f, err := os.OpenFile(kubeconfigFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
		if err != nil {
			output.Fatal("Failed to create %s: %s", kubeconfigFile, err)
		}
		defer f.Close()
	}
}

func Contexts(config *clientcmdapi.Config) []*types.ContextProfile {
	var contexts []*types.ContextProfile
	for contextName, context := range config.Contexts {
		item := &types.ContextProfile{
			Current:   contextName == config.CurrentContext,
			Name:      contextName,
			Cluster:   context.Cluster,
			User:      context.AuthInfo,
			Namespace: util.If(context.Namespace != "", context.Namespace, "default"),
			Server:    config.Clusters[context.Cluster].Server,
		}
		item.Emoji = util.If(item.Current, "✦", " ")
		contexts = append(contexts, item)
	}

	sort.Slice(contexts, func(i, j int) bool {
		return contexts[i].Name < contexts[j].Name
	})

	return contexts
}
