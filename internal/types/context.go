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

package types

import "github.com/fatih/color"

// ContextProfile represents a context profile
type ContextProfile struct {
	Current        bool
	Name           string
	Cluster        string
	User           string
	Server         string
	Namespace      string
	Emoji          string
	ClusterStatus  ClusterStatus
	ClusterVersion string
}

type ClusterStatus string

const (
	ClusterStatusAvailable   ClusterStatus = "✓ Available"
	ClusterStatusTimeout     ClusterStatus = "✗ Timeout"
	ClusterStatusUnavailable ClusterStatus = "✗ Unavailable"
)

func (cs ClusterStatus) ColorString() string {
	switch cs {
	case ClusterStatusAvailable:
		return color.GreenString(string(cs))
	case ClusterStatusUnavailable:
		return color.RedString(string(cs))
	default:
		return string(cs)
	}
}
