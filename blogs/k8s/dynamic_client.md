```go
import (
  "encoding/json"
  "fmt"

  "k8s.io/apimachinery/pkg/api/errors"
  "k8s.io/apimachinery/pkg/api/meta"
  metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
  "k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
  "k8s.io/apimachinery/pkg/runtime/schema"
  "k8s.io/apimachinery/pkg/types"
  "k8s.io/client-go/util/retry"
  "k8s.io/klog"
)

// cluster scoped resources don't have namespaces.
// Default to the item's namespace, but clear it 
// for cluster scoped resources
func resourceDefaultNamespace(namespaced bool, defaultNamespace string) string {
  if namespaced {
    return defaultNamespace
  }
  return ""
}

// apiResource consults the REST mapper to translate an 
// <apiVersion, kind, namespace> tuple to a 
// unversioned.APIResource struct.
func (gc *GarbageCollector) apiResource(apiVersion, kind string) (schema.GroupVersionResource, bool, error) {
  fqKind := schema.FromAPIVersionAndKind(apiVersion, kind)
  mapping, err := gc.restMapper.RESTMapping(fqKind.GroupKind(), fqKind.Version)
  if err != nil {
    return schema.GroupVersionResource{}, false, newRESTMappingError(kind, apiVersion)
  }
  return mapping.Resource, mapping.Scope == meta.RESTScopeNamespace, nil
}

func (gc *GarbageCollector) deleteObject(item objectReference, policy *metav1.DeletionPropagation) error {
  resource, namespaced, err := gc.apiResource(item.APIVersion, item.Kind)
  if err != nil {
    return err
  }
  uid := item.UID
  preconditions := metav1.Preconditions{UID: &uid}
  deleteOptions := metav1.DeleteOptions{Preconditions: &preconditions, PropagationPolicy: policy}
  return gc.dynamicClient.
    Resource(resource).
    Namespace(resourceDefaultNamespace(namespaced, item.Namespace)).
    Delete(item.Name, &deleteOptions)
}

func (gc *GarbageCollector) getObject(item objectReference) (*unstructured.Unstructured, error) {
  resource, namespaced, err := gc.apiResource(item.APIVersion, item.Kind)
  if err != nil {
    return nil, err
  }
  return gc.dynamicClient.
    Resource(resource).
    Namespace(resourceDefaultNamespace(namespaced, item.Namespace)).
    Get(item.Name, metav1.GetOptions{})
}

func (gc *GarbageCollector) patchObject(item objectReference, patch []byte, pt types.PatchType) (*unstructured.Unstructured, error) {
  resource, namespaced, err := gc.apiResource(item.APIVersion, item.Kind)
  if err != nil {
    return nil, err
  }
  return gc.dynamicClient.
    Resource(resource).
    Namespace(resourceDefaultNamespace(namespaced, item.Namespace)).
    Patch(item.Name, pt, patch, metav1.PatchOptions{})
}
```
