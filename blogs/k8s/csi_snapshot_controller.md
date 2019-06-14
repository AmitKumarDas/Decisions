### references
- https://github.com/kubernetes-csi/external-snapshotter/blob/master/pkg/controller/snapshot_controller.go

```go
// while syncing volume snapshot content

// review: delete is first check
if isContentDeletionCandidate(content) {
  isUsed := ctrl.isSnapshotContentBeingUsed(content)
  if !isUsed {
    klog.V(5).Infof("Remove Finalizer for VolumeSnapshotContent[%s]", content.Name)
    return ctrl.removeContentFinalizer(content)
  }
}

// review: safeguard is second check
if needToAddContentFinalizer(content) {
  klog.V(5).Infof("Add Finalizer for VolumeSnapshotContent[%s]", content.Name)
  return ctrl.addContentFinalizer(content)
}

// review: event & error
if content.Spec.VolumeSnapshotRef == nil {
  // content is not bound
  klog.V(4).Infof(
    "VolumeSnapshotContent[%s] is not bound to any VolumeSnapshot", 
    content.Name,
  )
	ctrl.eventRecorder.Event(
	  content, 
	  v1.EventTypeWarning, 
	  "SnapshotContentNotBound", 
	  "VolumeSnapshotContent is not bound to any VolumeSnapshot",
	)
	
  return fmt.Errorf(
    "volumeSnapshotContent %s is not bound to any VolumeSnapshot", content.Name,
  )
}

// review: log & exit with no error
if content.Spec.VolumeSnapshotRef.UID == "" {
  klog.V(4).Infof(
    "VolumeSnapshotContent[%s] is pre-bound to VolumeSnapshot %s", 
    content.Name, 
    snapshotRefKey(content.Spec.VolumeSnapshotRef),
  )
  return nil
}
```

- https://github.com/kubernetes-csi/external-snapshotter/blob/master/pkg/controller/snapshot_controller_base.go
```go
type csiSnapshotController struct {
  clientset       clientset.Interface
  client          kubernetes.Interface
  snapshotterName string
  eventRecorder   record.EventRecorder
  snapshotQueue   workqueue.RateLimitingInterface
  contentQueue    workqueue.RateLimitingInterface

  snapshotLister       storagelisters.VolumeSnapshotLister
  snapshotListerSynced cache.InformerSynced
  snapshotStore cache.Store

  contentLister        storagelisters.VolumeSnapshotContentLister
  contentListerSynced  cache.InformerSynced
  contentStore  cache.Store

  classLister          storagelisters.VolumeSnapshotClassLister
  classListerSynced    cache.InformerSynced

  pvcLister            corelisters.PersistentVolumeClaimLister
  pvcListerSynced      cache.InformerSynced

  handler Handler

  // Map of scheduled/running operations.
  runningOperations goroutinemap.GoRoutineMap

  createSnapshotContentRetryCount int
  createSnapshotContentInterval   time.Duration
  resyncPeriod                    time.Duration
}
```

```go
// constructor: event

broadcaster := record.NewBroadcaster().
  StartLogging(klog.Infof).
  StartRecordingToSink(
    &corev1.EventSinkImpl{
      Interface: client.CoreV1().Events(v1.NamespaceAll),
    },
  )

var eventRecorder record.EventRecorder
eventRecorder = broadcaster.NewRecorder(
  scheme.Scheme, 
  v1.EventSource{
    Component: fmt.Sprintf("csi-snapshotter %s", snapshotterName),
  },
)
```

```go
// constructor: store & queue

ctrl.snapshotStore = cache.NewStore(cache.DeletionHandlingMetaNamespaceKeyFunc)
ctrl.contentStore = cache.NewStore(cache.DeletionHandlingMetaNamespaceKeyFunc)

ctrl.snapshotQueue = workqueue.NewNamedRateLimitingQueue(
  workqueue.DefaultControllerRateLimiter(),
  "csi-snapshotter-snapshot"
)

ctrl.contentQueue = workqueue.NewNamedRateLimitingQueue(
  workqueue.DefaultControllerRateLimiter(), 
  "csi-snapshotter-content",
)
```

```go
// constructor: listers & informers

import (
  storageinformers "github.com/kubernetes-csi/external-snapshotter/pkg/client/informers/externalversions/volumesnapshot/v1alpha1"
  storagelisters "github.com/kubernetes-csi/external-snapshotter/pkg/client/listers/volumesnapshot/v1alpha1"

  coreinformers "k8s.io/client-go/informers/core/v1"
  corelisters "k8s.io/client-go/listers/core/v1"
)

// NewCSISnapshotController returns a new *csiSnapshotController
func NewCSISnapshotController(
  volumeSnapshotInformer        storageinformers.VolumeSnapshotInformer,
  volumeSnapshotContentInformer storageinformers.VolumeSnapshotContentInformer,
  volumeSnapshotClassInformer   storageinformers.VolumeSnapshotClassInformer,
  pvcInformer                   coreinformers.PersistentVolumeClaimInformer,
  resyncPeriod                  time.Duration,
) *csiSnapshotController {

ctrl.pvcLister = pvcInformer.Lister()
ctrl.pvcListerSynced = pvcInformer.Informer().HasSynced

ctrl.classLister = volumeSnapshotClassInformer.Lister()
ctrl.classListerSynced = volumeSnapshotClassInformer.Informer().HasSynced

volumeSnapshotInformer.Informer().AddEventHandlerWithResyncPeriod(
  cache.ResourceEventHandlerFuncs{
    AddFunc:    func(obj interface{}) { ctrl.enqueueSnapshotWork(obj) },
    UpdateFunc: func(oldObj, newObj interface{}) { ctrl.enqueueSnapshotWork(newObj) },
    DeleteFunc: func(obj interface{}) { ctrl.enqueueSnapshotWork(obj) },
  },
ctrl.resyncPeriod,
)
ctrl.snapshotLister = volumeSnapshotInformer.Lister()
ctrl.snapshotListerSynced = volumeSnapshotInformer.Informer().HasSynced

volumeSnapshotContentInformer.Informer().AddEventHandlerWithResyncPeriod(
  cache.ResourceEventHandlerFuncs{
    AddFunc:    func(obj interface{}) { ctrl.enqueueContentWork(obj) },
    UpdateFunc: func(oldObj, newObj interface{}) { ctrl.enqueueContentWork(newObj) },
    DeleteFunc: func(obj interface{}) { ctrl.enqueueContentWork(obj) },
  },
  ctrl.resyncPeriod,
)
ctrl.contentLister = volumeSnapshotContentInformer.Lister()
ctrl.contentListerSynced = volumeSnapshotContentInformer.Informer().HasSynced
}
```
