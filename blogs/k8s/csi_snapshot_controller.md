### references
- https://github.com/kubernetes-csi/external-snapshotter/blob/master/pkg/controller/snapshot_controller.go

```go
// while syncing volume snapshot content

// review: delete is first check
if isContentDeletionCandidate(content) {
  isUsed := ctrl.isSnapshotContentBeingUsed(content)
  if !isUsed {
    klog.V(5).Infof("syncContent: Remove Finalizer for VolumeSnapshotContent[%s]", content.Name)
    return ctrl.removeContentFinalizer(content)
  }
}

// review: safeguard is second check
if needToAddContentFinalizer(content) {
  klog.V(5).Infof("syncContent: Add Finalizer for VolumeSnapshotContent[%s]", content.Name)
  return ctrl.addContentFinalizer(content)
}

// review: event & error
if content.Spec.VolumeSnapshotRef == nil {
  // content is not bound
  klog.V(4).Infof(
    "synchronizing VolumeSnapshotContent[%s]: VolumeSnapshotContent is not bound to any VolumeSnapshot", 
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
    "synchronizing VolumeSnapshotContent[%s]: VolumeSnapshotContent is pre-bound to VolumeSnapshot %s", 
    content.Name, 
    snapshotRefKey(content.Spec.VolumeSnapshotRef),
  )
  return nil
}
```
