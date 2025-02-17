/*
Copyright 2020 The Nocalhost Authors.
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

package profile

import (
	"strconv"
)

// Deprecated
type NocalHostAppConfig struct {
	PreInstall   []*PreInstallItem    `json:"onPreInstall" yaml:"onPreInstall"`
	SvcConfigs   []*ServiceDevOptions `json:"services" yaml:"services"`
	Name         string               `json:"name" yaml:"name"`
	Type         string               `json:"manifestType" yaml:"manifestType"`
	ResourcePath []string             `json:"resourcePath" yaml:"resourcePath"`
	IgnoredPath  []string             `json:"ignoredPath" yaml:"ignoredPath"`
}

type PersistentVolumeDir struct {
	Path     string `json:"path" yaml:"path"`
	Capacity string `json:"capacity,omitempty" yaml:"capacity,omitempty"`
}

type ResourceQuota struct {
	Limits   *QuotaList `json:"limits" yaml:"limits"`
	Requests *QuotaList `json:"requests" yaml:"requests"`
}

type QuotaList struct {
	Memory string `json:"memory" yaml:"memory"`
	Cpu    string `json:"cpu" yaml:"cpu"`
}

type ServiceDevOptions struct {
	Name                  string                 `json:"name" yaml:"name"`
	Type                  string                 `json:"serviceType" yaml:"serviceType"`
	GitUrl                string                 `json:"gitUrl" yaml:"gitUrl"`
	DevImage              string                 `json:"devContainerImage" yaml:"devContainerImage"`
	WorkDir               string                 `json:"workDir" yaml:"workDir"`
	Sync                  []string               `json:"syncDirs" yaml:"syncDirs,omitempty"` // dev start -s
	PriorityClass         string                 `json:"priorityClass,omitempty" yaml:"priorityClass,omitempty"`
	PersistentVolumeDirs  []*PersistentVolumeDir `json:"persistentVolumeDirs" yaml:"persistentVolumeDirs"`
	BuildCommand          []string               `json:"buildCommand,omitempty" yaml:"buildCommand,omitempty"`
	RunCommand            []string               `json:"runCommand,omitempty" yaml:"runCommand,omitempty"`
	DebugCommand          []string               `json:"debugCommand,omitempty" yaml:"debugCommand,omitempty"`
	HotReloadRunCommand   []string               `json:"hotReloadRunCommand,omitempty" yaml:"hotReloadRunCommand,omitempty"`
	HotReloadDebugCommand []string               `json:"hotReloadDebugCommand,omitempty" yaml:"hotReloadDebugCommand,omitempty"`
	DevContainerShell     string                 `json:"devContainerShell" yaml:"devContainerShell"`
	DevContainerResources *ResourceQuota         `json:"devContainerResources" yaml:"devContainerResources"`
	DevPort               []string               `json:"devPorts" yaml:"devPorts"`
	Jobs                  []string               `json:"dependJobsLabelSelector" yaml:"dependJobsLabelSelector,omitempty"`
	Pods                  []string               `json:"dependPodsLabelSelector" yaml:"dependPodsLabelSelector,omitempty"`
	SyncedPattern         []string               `json:"syncFilePattern" yaml:"syncFilePattern"`
	IgnoredPattern        []string               `json:"ignoreFilePattern" yaml:"ignoreFilePattern"`
}

type ComparableItems []*PreInstallItem

func (a ComparableItems) Len() int      { return len(a) }
func (a ComparableItems) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ComparableItems) Less(i, j int) bool {
	iW, err := strconv.Atoi(a[i].Weight)
	if err != nil {
		iW = 0
	}

	jW, err := strconv.Atoi(a[j].Weight)
	if err != nil {
		jW = 0
	}
	return iW < jW
}

func (n *NocalHostAppConfig) GetSvcConfig(name string) *ServiceDevOptions {
	if n.SvcConfigs == nil {
		return nil
	}
	for _, svc := range n.SvcConfigs {
		if svc.Name == name {
			return svc
		}
	}
	return nil
}
