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
	"context"
	"sync"
	"time"

	"github.com/kubernot/wrangler/pkg/apply"
	"github.com/kubernot/wrangler/pkg/condition"
	"github.com/kubernot/wrangler/pkg/generic"
	"github.com/kubernot/wrangler/pkg/kv"
	"k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	v1 "k8s.io/kube-aggregator/pkg/apis/apiregistration/v1"
)

// APIServiceController interface for managing APIService resources.
type APIServiceController interface {
	generic.NonNamespacedControllerInterface[*v1.APIService, *v1.APIServiceList]
}

// APIServiceClient interface for managing APIService resources in Kubernetes.
type APIServiceClient interface {
	generic.NonNamespacedClientInterface[*v1.APIService, *v1.APIServiceList]
}

// APIServiceCache interface for retrieving APIService resources in memory.
type APIServiceCache interface {
	generic.NonNamespacedCacheInterface[*v1.APIService]
}

// APIServiceStatusHandler is executed for every added or modified APIService. Should return the new status to be updated
type APIServiceStatusHandler func(obj *v1.APIService, status v1.APIServiceStatus) (v1.APIServiceStatus, error)

// APIServiceGeneratingHandler is the top-level handler that is executed for every APIService event. It extends APIServiceStatusHandler by a returning a slice of child objects to be passed to apply.Apply
type APIServiceGeneratingHandler func(obj *v1.APIService, status v1.APIServiceStatus) ([]runtime.Object, v1.APIServiceStatus, error)

// RegisterAPIServiceStatusHandler configures a APIServiceController to execute a APIServiceStatusHandler for every events observed.
// If a non-empty condition is provided, it will be updated in the status conditions for every handler execution
func RegisterAPIServiceStatusHandler(ctx context.Context, controller APIServiceController, condition condition.Cond, name string, handler APIServiceStatusHandler) {
	statusHandler := &aPIServiceStatusHandler{
		client:    controller,
		condition: condition,
		handler:   handler,
	}
	controller.AddGenericHandler(ctx, name, generic.FromObjectHandlerToHandler(statusHandler.sync))
}

// RegisterAPIServiceGeneratingHandler configures a APIServiceController to execute a APIServiceGeneratingHandler for every events observed, passing the returned objects to the provided apply.Apply.
// If a non-empty condition is provided, it will be updated in the status conditions for every handler execution
func RegisterAPIServiceGeneratingHandler(ctx context.Context, controller APIServiceController, apply apply.Apply,
	condition condition.Cond, name string, handler APIServiceGeneratingHandler, opts *generic.GeneratingHandlerOptions) {
	statusHandler := &aPIServiceGeneratingHandler{
		APIServiceGeneratingHandler: handler,
		apply:                       apply,
		name:                        name,
		gvk:                         controller.GroupVersionKind(),
	}
	if opts != nil {
		statusHandler.opts = *opts
	}
	controller.OnChange(ctx, name, statusHandler.Remove)
	RegisterAPIServiceStatusHandler(ctx, controller, condition, name, statusHandler.Handle)
}

type aPIServiceStatusHandler struct {
	client    APIServiceClient
	condition condition.Cond
	handler   APIServiceStatusHandler
}

// sync is executed on every resource addition or modification. Executes the configured handlers and sends the updated status to the Kubernetes API
func (a *aPIServiceStatusHandler) sync(key string, obj *v1.APIService) (*v1.APIService, error) {
	if obj == nil {
		return obj, nil
	}

	origStatus := obj.Status.DeepCopy()
	obj = obj.DeepCopy()
	newStatus, err := a.handler(obj, obj.Status)
	if err != nil {
		// Revert to old status on error
		newStatus = *origStatus.DeepCopy()
	}

	if a.condition != "" {
		if errors.IsConflict(err) {
			a.condition.SetError(&newStatus, "", nil)
		} else {
			a.condition.SetError(&newStatus, "", err)
		}
	}
	if !equality.Semantic.DeepEqual(origStatus, &newStatus) {
		if a.condition != "" {
			// Since status has changed, update the lastUpdatedTime
			a.condition.LastUpdated(&newStatus, time.Now().UTC().Format(time.RFC3339))
		}

		var newErr error
		obj.Status = newStatus
		newObj, newErr := a.client.UpdateStatus(obj)
		if err == nil {
			err = newErr
		}
		if newErr == nil {
			obj = newObj
		}
	}
	return obj, err
}

type aPIServiceGeneratingHandler struct {
	APIServiceGeneratingHandler
	apply apply.Apply
	opts  generic.GeneratingHandlerOptions
	gvk   schema.GroupVersionKind
	name  string
	seen  sync.Map
}

// Remove handles the observed deletion of a resource, cascade deleting every associated resource previously applied
func (a *aPIServiceGeneratingHandler) Remove(key string, obj *v1.APIService) (*v1.APIService, error) {
	if obj != nil {
		return obj, nil
	}

	obj = &v1.APIService{}
	obj.Namespace, obj.Name = kv.RSplit(key, "/")
	obj.SetGroupVersionKind(a.gvk)

	if a.opts.UniqueApplyForResourceVersion {
		a.seen.Delete(key)
	}

	return nil, generic.ConfigureApplyForObject(a.apply, obj, &a.opts).
		WithOwner(obj).
		WithSetID(a.name).
		ApplyObjects()
}

// Handle executes the configured APIServiceGeneratingHandler and pass the resulting objects to apply.Apply, finally returning the new status of the resource
func (a *aPIServiceGeneratingHandler) Handle(obj *v1.APIService, status v1.APIServiceStatus) (v1.APIServiceStatus, error) {
	if !obj.DeletionTimestamp.IsZero() {
		return status, nil
	}

	objs, newStatus, err := a.APIServiceGeneratingHandler(obj, status)
	if err != nil {
		return newStatus, err
	}
	if !a.isNewResourceVersion(obj) {
		return newStatus, nil
	}

	err = generic.ConfigureApplyForObject(a.apply, obj, &a.opts).
		WithOwner(obj).
		WithSetID(a.name).
		ApplyObjects(objs...)
	if err != nil {
		return newStatus, err
	}
	a.storeResourceVersion(obj)
	return newStatus, nil
}

// isNewResourceVersion detects if a specific resource version was already successfully processed.
// Only used if UniqueApplyForResourceVersion is set in generic.GeneratingHandlerOptions
func (a *aPIServiceGeneratingHandler) isNewResourceVersion(obj *v1.APIService) bool {
	if !a.opts.UniqueApplyForResourceVersion {
		return true
	}

	// Apply once per resource version
	key := obj.Namespace + "/" + obj.Name
	previous, ok := a.seen.Load(key)
	return !ok || previous != obj.ResourceVersion
}

// storeResourceVersion keeps track of the latest resource version of an object for which Apply was executed
// Only used if UniqueApplyForResourceVersion is set in generic.GeneratingHandlerOptions
func (a *aPIServiceGeneratingHandler) storeResourceVersion(obj *v1.APIService) {
	if !a.opts.UniqueApplyForResourceVersion {
		return
	}

	key := obj.Namespace + "/" + obj.Name
	a.seen.Store(key, obj.ResourceVersion)
}
