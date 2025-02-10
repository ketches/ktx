/*
Copyright © 2025 Pone Ding <poneding@gmail.com>
*/
package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/poneding/ktx/internal/kubeconfig"
	"github.com/poneding/ktx/internal/output"
	"github.com/poneding/ktx/internal/prompt"
	"github.com/poneding/ktx/internal/util"
	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/util/rand"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

var (
	addFile string
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "add context from kubeconfig file to ~/.kube/config",
	Long:  `add context from kubeconfig file to ~/.kube/config`,
	Run: func(cmd *cobra.Command, args []string) {
		runAdd()
	},
}

func init() {
	rootCmd.AddCommand(addCmd)

	addCmd.Flags().StringVarP(&addFile, "file", "f", "", "kubeconfig file")

	addCmd.MarkFlagRequired("file")
}

func runAdd() {
	config := kubeconfig.Load()

	if !util.IsFileExist(addFile) {
		output.Fatal("File %s not found.", addFile)
	}

	new := kubeconfig.LoadFromFile(addFile)

	merge(config, new)
}

func merge(curr, new *clientcmdapi.Config) {
	for newCtxName, newCtx := range new.Contexts {
		newCluster, ok := new.Clusters[newCtx.Cluster]
		if !ok {
			output.Note("Cluster not found for context <%s> in file %s, skipped.", newCtxName, addFile)
			continue
		}

		newUser, ok := new.AuthInfos[newCtx.AuthInfo]
		if !ok {
			output.Note("User not found for context <%s> in file %s, skipped.", newCtxName, addFile)
			continue
		}

		// 如果 context 名称已经存在，要求用户输入新的 context 名称
		var quitWithConflict bool
		for contextNameConflict(newCtxName, curr) {
			if prompt.YesNo(fmt.Sprintf("Context name <%s> already exists, rename it", newCtxName)) == "Yes" {
				newCtxName = prompt.TextInput("Enter a new name")
			} else {
				quitWithConflict = true
				break
			}
		}

		if quitWithConflict {
			continue
		}

		mf := &mergeFrom{
			contextName:  newCtxName,
			clusterName:  newCtx.Cluster,
			userName:     newCtx.AuthInfo,
			context:      newCtx,
			cluster:      newCluster,
			user:         newUser,
			uniqueSuffix: rand.String(5),
		}

		handleMerge(curr, mf)
	}

	kubeconfig.Save(curr)

	// 如果当前没有 context，那么提示用户选择一个 context
	if curr.CurrentContext == "" {
		curr.CurrentContext = prompt.ContextSelection("Select context to use", curr)
		kubeconfig.Save(curr)
	}
}

type mergeFrom struct {
	contextName, clusterName, userName string
	context                            *clientcmdapi.Context
	cluster                            *clientcmdapi.Cluster
	user                               *clientcmdapi.AuthInfo
	uniqueSuffix                       string
}

func handleMerge(current *clientcmdapi.Config, mf *mergeFrom) {
	// merge cluster
	var clusterExist, userExist bool
	for clusterName, cluster := range current.Clusters {
		if marshalEqual(cluster, mf.cluster) {
			clusterExist = true
			// 存在相同的 cluster，不需要添加
			output.Note("Cluster %s already exists, skipped.", mf.clusterName)
			mf.clusterName = clusterName
			mf.context.Cluster = clusterName
			break
		}
	}
	if !clusterExist {
		if _, ok := current.Clusters[mf.clusterName]; ok {
			// cluster 名称冲突，需要重新命名
			mf.clusterName = fmt.Sprintf("%s-%s", mf.clusterName, mf.uniqueSuffix)
		}
		mf.context.Cluster = mf.clusterName
		current.Clusters[mf.clusterName] = mf.cluster
	}

	if _, ok := current.Clusters[mf.clusterName]; !ok {
		current.Clusters[mf.clusterName] = mf.cluster
	}

	// merge user
	for userName, user := range current.AuthInfos {
		if marshalEqual(user, mf.user) {
			userExist = true
			// 存在相同的 user，不需要添加
			output.Note("User %s already exists, skipped.", mf.userName)
			mf.userName = userName
			mf.context.AuthInfo = userName
			break
		}
	}
	if !userExist {
		if _, ok := current.AuthInfos[mf.userName]; ok {
			// auth info 名称冲突，需要重新命名
			mf.userName = fmt.Sprintf("%s-%s", mf.userName, mf.uniqueSuffix)
		}
		mf.context.AuthInfo = mf.userName
		current.AuthInfos[mf.userName] = mf.user
	}
	if _, ok := current.AuthInfos[mf.userName]; !ok {
		current.AuthInfos[mf.userName] = mf.user
	}

	// merge context
	current.Contexts[mf.contextName] = mf.context

	output.Done("Context <%s> added.", mf.contextName)
}

func contextNameConflict(name string, config *clientcmdapi.Config) bool {
	_, ok := config.Contexts[name]
	return ok
}

func marshalEqual(a, b any) bool {
	ab, _ := json.Marshal(a)
	bb, _ := json.Marshal(b)
	return bytes.Equal(ab, bb)
}
