### Reference
- https://github.com/kubernetes-sigs/kubefed/blob/master/pkg/controller/dnsendpoint/controller.go

```go
  apierrors "k8s.io/apimachinery/pkg/api/errors"
  metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
  pkgruntime "k8s.io/apimachinery/pkg/runtime"
  "k8s.io/apimachinery/pkg/util/runtime"
  "k8s.io/apimachinery/pkg/util/wait"
  "k8s.io/client-go/tools/cache"
  "k8s.io/client-go/util/workqueue"

  feddnsv1a1 "sigs.k8s.io/kubefed/pkg/apis/multiclusterdns/v1alpha1"
  genericclient "sigs.k8s.io/kubefed/pkg/client/generic"
  "sigs.k8s.io/kubefed/pkg/controller/util"
```

```go
type GetEndpointsFunc func(interface{}) ([]*feddnsv1a1.Endpoint, error)

type controller struct {
	client genericclient.Client

	// Informer Store for DNS objects
	dnsObjectStore cache.Store

	// Informer controller for DNS objects
	dnsObjectController cache.Controller

	dnsObjectKind string
	getEndpoints  GetEndpointsFunc

	queue         workqueue.RateLimitingInterface
	minRetryDelay time.Duration
	maxRetryDelay time.Duration
}
```

```go
// ----------------------------------------
// instantiating the controller...
// ----------------------------------------
```

```go
// TODO: Change this to **SHARED INFORMER**
// Start informer for DNS objects
//
// NOTE:
//  NewGenericInformer accepts a few functions as well
//
//  It provides the informer store and informer controller
d.dnsObjectStore, d.dnsObjectController, err = util.NewGenericInformer(
  config.KubeConfig,
  config.TargetNamespace,
  objectType,
  util.NoResyncPeriod,
  d.enqueueObject,
)

if minimizeLatency {
	d.minimizeLatency()
}

d.queue = workqueue.NewNamedRateLimitingQueue(
  workqueue.NewItemExponentialFailureRateLimiter(
	  d.minRetryDelay, 
	  d.maxRetryDelay,
	), 
	objectKind+"DNSEndpoint",
)
```
```go
// ---------------------------------------------
// Run
// ---------------------------------------------
```
```go
func (d *controller) Run(stopCh <-chan struct{}) {
  defer runtime.HandleCrash()
  defer d.queue.ShutDown()

  klog.Infof("Starting %q DNSEndpoint controller", d.dnsObjectKind)
  defer klog.Infof("Shutting down %q DNSEndpoint controller", d.dnsObjectKind)

  // remember !! this is the informer controller
  //
  // What does this do? 
  //  Is informer itself a controller whose job is to look out 
  // for changes to particular kind?
  go d.dnsObjectController.Run(stopCh)

  // wait for the caches to synchronize before starting the worker
  if !cache.WaitForCacheSync(stopCh, d.dnsObjectController.HasSynced) {
	  runtime.HandleError(errors.New("Timed out waiting for caches to sync"))
	  return
  }

  klog.Infof("%q DNSEndpoint controller synced and ready", d.dnsObjectKind)

  // It starts off a number of workers
  // Each worker do their job in their individual goroutine
  for i := 0; i < numWorkers; i++ {
    // here d.worker is a function that
    // satisfies the signature of wait.Until
	  go wait.Until(d.worker, time.Second, stopCh)
  }

<-stopCh
}
```

```go
// --------------------------------------------
// workers
// --------------------------------------------
```
```go
func (d *controller) worker() {
  // processNextWorkItem will automatically wait until there's work available
  for d.processNextItem() {
	  // continue looping
  }
}

func (d *controller) processNextItem() bool {
  key, quit := d.queue.Get()
  if quit {
    // this is the only instance when
    // false is returned
	  return false
  }
  defer d.queue.Done(key)

  // actual reconcile happens here !!!
  err := d.processItem(key.(string))

  if err == nil {
	  // No error, tell the queue to stop tracking history
	  d.queue.Forget(key)

  } else if d.queue.NumRequeues(key) < maxRetries {
	  klog.Errorf("Error processing %s (will retry): %v", key, err)
	  // requeue the item to work on later
	  d.queue.AddRateLimited(key)

  } else {
	  // err != nil and too many retries
	  klog.Errorf("Error processing %s (giving up): %v", key, err)
	  d.queue.Forget(key)
	  runtime.HandleError(err)
  }

  return true
}
```

```go
// ---------------------------------
// other functions/methods
// ---------------------------------
```

```go
// enqueueObject is fed as a function to NewGenericInformer
//
// Does this generic informer have the workqueue?
//  I guess generic informer provides a hook only. The workqueue
// is part of this controller instance
func (d *controller) enqueueObject(obj pkgruntime.Object) {
  key, err := cache.DeletionHandlingMetaNamespaceKeyFunc(obj)
  if err != nil {
	  klog.Errorf("Couldn't get key for object %#v: %v", obj, err)
	  return
  }
	
  d.queue.Add(key)
}

// Nice way of logging with timetaken !!!
func (d *controller) processItem(key string) error {
  startTime := time.Now()
  klog.V(4).Infof("Processing change to %q DNSEndpoint %s", d.dnsObjectKind, key)
  defer func() {
	  klog.V(4).Infof("Finished processing %q DNSEndpoint %q (%v)", d.dnsObjectKind, key, time.Since(startTime))
  }()

	...
}

// minimizeLatency reduces delays and timeouts to make the controller 
// more responsive (useful for testing).
func (d *controller) minimizeLatency() {
	d.minRetryDelay = 50 * time.Millisecond
	d.maxRetryDelay = 2 * time.Second
}
```
