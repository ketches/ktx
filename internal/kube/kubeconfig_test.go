/*
Copyright 2025 The Ketches Authors.

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

package kube

import (
	"testing"

	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

func TestStandardizeConfig(t *testing.T) {
	testdata := []struct {
		config   *clientcmdapi.Config
		expected *clientcmdapi.Config
	}{
		{
			config: &clientcmdapi.Config{
				Clusters:  map[string]*clientcmdapi.Cluster{"cluster1": {}, "cluster2": {}},
				AuthInfos: map[string]*clientcmdapi.AuthInfo{"user1": {}, "user2": {}},
				Contexts: map[string]*clientcmdapi.Context{
					"context1": {Cluster: "cluster1", AuthInfo: "user1"},
					"context2": {Cluster: "cluster2", AuthInfo: "user2"},
				},
			},
			expected: &clientcmdapi.Config{
				Clusters:  map[string]*clientcmdapi.Cluster{"cluster-context1": {}, "cluster-context2": {}},
				AuthInfos: map[string]*clientcmdapi.AuthInfo{"user-context1": {}, "user-context2": {}},
				Contexts: map[string]*clientcmdapi.Context{
					"context1": {Cluster: "cluster-context1", AuthInfo: "user-context1"},
					"context2": {Cluster: "cluster-context2", AuthInfo: "user-context2"},
				},
			},
		},
		{
			config: &clientcmdapi.Config{
				Clusters:  map[string]*clientcmdapi.Cluster{"cluster1": {}},
				AuthInfos: map[string]*clientcmdapi.AuthInfo{"user1": {}},
				Contexts: map[string]*clientcmdapi.Context{
					"context1": {Cluster: "cluster1", AuthInfo: "user1"},
					"context2": {Cluster: "cluster1", AuthInfo: "user1"},
				},
			},
			expected: &clientcmdapi.Config{
				Clusters:  map[string]*clientcmdapi.Cluster{"cluster-context1": {}, "cluster-context2": {}},
				AuthInfos: map[string]*clientcmdapi.AuthInfo{"user-context1": {}, "user-context2": {}},
				Contexts: map[string]*clientcmdapi.Context{
					"context1": {Cluster: "cluster-context1", AuthInfo: "user-context1"},
					"context2": {Cluster: "cluster-context2", AuthInfo: "user-context2"},
				},
			},
		},
	}

	for _, test := range testdata {
		StandardizeConfig(test.config)
		for ctxName, ctx := range test.config.Contexts {
			if ctx.Cluster != "cluster-"+ctxName {
				t.Errorf("StandardizeConfig() failed, context: %s, expected cluster name: %s, got: %s", ctxName, "cluster-"+ctxName, ctx.Cluster)
			}
			if ctx.AuthInfo != "user-"+ctxName {
				t.Errorf("StandardizeConfig() failed, context: %s, expected user name: %s, got: %s", ctxName, "user-"+ctxName, ctx.AuthInfo)
			}
		}
	}
}
