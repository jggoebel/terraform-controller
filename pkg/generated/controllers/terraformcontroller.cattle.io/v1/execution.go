/*
Copyright 2019 Rancher Labs, Inc.

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

	v1 "github.com/rancher/terraform-controller/pkg/apis/terraformcontroller.cattle.io/v1"
	clientset "github.com/rancher/terraform-controller/pkg/generated/clientset/versioned/typed/terraformcontroller.cattle.io/v1"
	informers "github.com/rancher/terraform-controller/pkg/generated/informers/externalversions/terraformcontroller.cattle.io/v1"
	listers "github.com/rancher/terraform-controller/pkg/generated/listers/terraformcontroller.cattle.io/v1"
	"github.com/rancher/wrangler/pkg/generic"
	"k8s.io/apimachinery/pkg/api/equality"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/tools/cache"
)

type ExecutionHandler func(string, *v1.Execution) (*v1.Execution, error)

type ExecutionController interface {
	ExecutionClient

	OnChange(ctx context.Context, name string, sync ExecutionHandler)
	OnRemove(ctx context.Context, name string, sync ExecutionHandler)
	Enqueue(namespace, name string)

	Cache() ExecutionCache

	Informer() cache.SharedIndexInformer
	GroupVersionKind() schema.GroupVersionKind

	AddGenericHandler(ctx context.Context, name string, handler generic.Handler)
	AddGenericRemoveHandler(ctx context.Context, name string, handler generic.Handler)
	Updater() generic.Updater
}

type ExecutionClient interface {
	Create(*v1.Execution) (*v1.Execution, error)
	Update(*v1.Execution) (*v1.Execution, error)
	UpdateStatus(*v1.Execution) (*v1.Execution, error)
	Delete(namespace, name string, options *metav1.DeleteOptions) error
	Get(namespace, name string, options metav1.GetOptions) (*v1.Execution, error)
	List(namespace string, opts metav1.ListOptions) (*v1.ExecutionList, error)
	Watch(namespace string, opts metav1.ListOptions) (watch.Interface, error)
	Patch(namespace, name string, pt types.PatchType, data []byte, subresources ...string) (result *v1.Execution, err error)
}

type ExecutionCache interface {
	Get(namespace, name string) (*v1.Execution, error)
	List(namespace string, selector labels.Selector) ([]*v1.Execution, error)

	AddIndexer(indexName string, indexer ExecutionIndexer)
	GetByIndex(indexName, key string) ([]*v1.Execution, error)
}

type ExecutionIndexer func(obj *v1.Execution) ([]string, error)

type executionController struct {
	controllerManager *generic.ControllerManager
	clientGetter      clientset.ExecutionsGetter
	informer          informers.ExecutionInformer
	gvk               schema.GroupVersionKind
}

func NewExecutionController(gvk schema.GroupVersionKind, controllerManager *generic.ControllerManager, clientGetter clientset.ExecutionsGetter, informer informers.ExecutionInformer) ExecutionController {
	return &executionController{
		controllerManager: controllerManager,
		clientGetter:      clientGetter,
		informer:          informer,
		gvk:               gvk,
	}
}

func FromExecutionHandlerToHandler(sync ExecutionHandler) generic.Handler {
	return func(key string, obj runtime.Object) (ret runtime.Object, err error) {
		var v *v1.Execution
		if obj == nil {
			v, err = sync(key, nil)
		} else {
			v, err = sync(key, obj.(*v1.Execution))
		}
		if v == nil {
			return nil, err
		}
		return v, err
	}
}

func (c *executionController) Updater() generic.Updater {
	return func(obj runtime.Object) (runtime.Object, error) {
		newObj, err := c.Update(obj.(*v1.Execution))
		if newObj == nil {
			return nil, err
		}
		return newObj, err
	}
}

func UpdateExecutionOnChange(updater generic.Updater, handler ExecutionHandler) ExecutionHandler {
	return func(key string, obj *v1.Execution) (*v1.Execution, error) {
		if obj == nil {
			return handler(key, nil)
		}

		copyObj := obj.DeepCopy()
		newObj, err := handler(key, copyObj)
		if newObj != nil {
			copyObj = newObj
		}
		if obj.ResourceVersion == copyObj.ResourceVersion && !equality.Semantic.DeepEqual(obj, copyObj) {
			newObj, err := updater(copyObj)
			if newObj != nil && err == nil {
				copyObj = newObj.(*v1.Execution)
			}
		}

		return copyObj, err
	}
}

func (c *executionController) AddGenericHandler(ctx context.Context, name string, handler generic.Handler) {
	c.controllerManager.AddHandler(ctx, c.gvk, c.informer.Informer(), name, handler)
}

func (c *executionController) AddGenericRemoveHandler(ctx context.Context, name string, handler generic.Handler) {
	removeHandler := generic.NewRemoveHandler(name, c.Updater(), handler)
	c.controllerManager.AddHandler(ctx, c.gvk, c.informer.Informer(), name, removeHandler)
}

func (c *executionController) OnChange(ctx context.Context, name string, sync ExecutionHandler) {
	c.AddGenericHandler(ctx, name, FromExecutionHandlerToHandler(sync))
}

func (c *executionController) OnRemove(ctx context.Context, name string, sync ExecutionHandler) {
	removeHandler := generic.NewRemoveHandler(name, c.Updater(), FromExecutionHandlerToHandler(sync))
	c.AddGenericHandler(ctx, name, removeHandler)
}

func (c *executionController) Enqueue(namespace, name string) {
	c.controllerManager.Enqueue(c.gvk, namespace, name)
}

func (c *executionController) Informer() cache.SharedIndexInformer {
	return c.informer.Informer()
}

func (c *executionController) GroupVersionKind() schema.GroupVersionKind {
	return c.gvk
}

func (c *executionController) Cache() ExecutionCache {
	return &executionCache{
		lister:  c.informer.Lister(),
		indexer: c.informer.Informer().GetIndexer(),
	}
}

func (c *executionController) Create(obj *v1.Execution) (*v1.Execution, error) {
	return c.clientGetter.Executions(obj.Namespace).Create(obj)
}

func (c *executionController) Update(obj *v1.Execution) (*v1.Execution, error) {
	return c.clientGetter.Executions(obj.Namespace).Update(obj)
}

func (c *executionController) UpdateStatus(obj *v1.Execution) (*v1.Execution, error) {
	return c.clientGetter.Executions(obj.Namespace).UpdateStatus(obj)
}

func (c *executionController) Delete(namespace, name string, options *metav1.DeleteOptions) error {
	return c.clientGetter.Executions(namespace).Delete(name, options)
}

func (c *executionController) Get(namespace, name string, options metav1.GetOptions) (*v1.Execution, error) {
	return c.clientGetter.Executions(namespace).Get(name, options)
}

func (c *executionController) List(namespace string, opts metav1.ListOptions) (*v1.ExecutionList, error) {
	return c.clientGetter.Executions(namespace).List(opts)
}

func (c *executionController) Watch(namespace string, opts metav1.ListOptions) (watch.Interface, error) {
	return c.clientGetter.Executions(namespace).Watch(opts)
}

func (c *executionController) Patch(namespace, name string, pt types.PatchType, data []byte, subresources ...string) (result *v1.Execution, err error) {
	return c.clientGetter.Executions(namespace).Patch(name, pt, data, subresources...)
}

type executionCache struct {
	lister  listers.ExecutionLister
	indexer cache.Indexer
}

func (c *executionCache) Get(namespace, name string) (*v1.Execution, error) {
	return c.lister.Executions(namespace).Get(name)
}

func (c *executionCache) List(namespace string, selector labels.Selector) ([]*v1.Execution, error) {
	return c.lister.Executions(namespace).List(selector)
}

func (c *executionCache) AddIndexer(indexName string, indexer ExecutionIndexer) {
	utilruntime.Must(c.indexer.AddIndexers(map[string]cache.IndexFunc{
		indexName: func(obj interface{}) (strings []string, e error) {
			return indexer(obj.(*v1.Execution))
		},
	}))
}

func (c *executionCache) GetByIndex(indexName, key string) (result []*v1.Execution, err error) {
	objs, err := c.indexer.ByIndex(indexName, key)
	if err != nil {
		return nil, err
	}
	for _, obj := range objs {
		result = append(result, obj.(*v1.Execution))
	}
	return result, nil
}