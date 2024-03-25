/*
Copyright The Kubernetes Authors.

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

// Code generated by main. DO NOT EDIT.

package v1

import (
	"github.com/kubernot/wrangler/pkg/generic"
	v1 "k8s.io/api/core/v1"
)

// ConfigMapController interface for managing ConfigMap resources.
type ConfigMapController interface {
	generic.ControllerInterface[*v1.ConfigMap, *v1.ConfigMapList]
}

// ConfigMapClient interface for managing ConfigMap resources in Kubernetes.
type ConfigMapClient interface {
	generic.ClientInterface[*v1.ConfigMap, *v1.ConfigMapList]
}

// ConfigMapCache interface for retrieving ConfigMap resources in memory.
type ConfigMapCache interface {
	generic.CacheInterface[*v1.ConfigMap]
}
