- refer:
  - https://github.com/kubernetes/kubernetes/blob/master/pkg/controller/garbagecollector/operations.go

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

func (gc *GarbageCollector) removeFinalizer(owner *node, targetFinalizer string) error {
	err := retry.RetryOnConflict(retry.DefaultBackoff, func() error {
		ownerObject, err := gc.getObject(owner.identity)
		if errors.IsNotFound(err) {
			return nil
		}
		if err != nil {
			return fmt.Errorf(
        "cannot finalize owner %s, because cannot get it: %v. The garbage collector will retry later", 
        owner.identity, 
        err,
      )
		}
    
		accessor, err := meta.Accessor(ownerObject)
		if err != nil {
			return fmt.Errorf(
        "cannot access the owner object %v: %v. The garbage collector will retry later", 
        ownerObject, 
        err,
      )
		}

		finalizers := accessor.GetFinalizers()
		var newFinalizers []string
		found := false
		for _, f := range finalizers {
			if f == targetFinalizer {
				found = true
				continue
			}
			newFinalizers = append(newFinalizers, f)
		}
		if !found {
			klog.V(5).Infof(
        "the %s finalizer is already removed from object %s", 
        targetFinalizer, 
        owner.identity,
      )
			return nil
		}

		// remove the owner from dependent's OwnerReferences
		patch, err := json.Marshal(map[string]interface{}{
			"metadata": map[string]interface{}{
				"resourceVersion": accessor.GetResourceVersion(),
				"finalizers":      newFinalizers,
			},
		})
		if err != nil {
			return fmt.Errorf(
        "unable to finalize %s due to an error serializing patch: %v", 
        owner.identity, 
        err,
      )
		}
    
		_, err = gc.patchObject(owner.identity, patch, types.MergePatchType)
		return err
	})
  
	if errors.IsConflict(err) {
		return fmt.Errorf(
      "updateMaxRetries(%d) has reached. The garbage collector will retry later for owner %v", 
      retry.DefaultBackoff.Steps, 
      owner.identity,
    )
	}
	return err
}
```
