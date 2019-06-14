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

// https://github.com/kubernetes-csi/external-snapshotter/blob/master/pkg/controller/snapshot_controller_base.go
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
