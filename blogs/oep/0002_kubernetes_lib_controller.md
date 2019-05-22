### Motivation
- A library for writing controllers
- Remove the learning curve and effort required to create kubernetes controllers
- Reduce boilerplate involved to write and maintain kubernetes controllers
- Can be used as a library i.e. via go import & go vendoring
- Can be copied via scripts/tools to repositories to avoid vendoring issues

### High Level Design

```go
// pkg/controller/lib/v1alpha1/controller.go

// Base is the common kubernetes based
// controller structure that can be composed
// by specific controllers
type Base struct {
  client Kubernetes.Interface
  
  // name of this controller
  name string
  
  // resyncPeriod is how often the controller relists
  // desired resources. OnUpdate will be called even
  // if nothing has changed, meaning failed operations
  // may be retried on the watched resource every
  // resyncPeriod regardless of whether it changed. 
  // Defaults to 15 minutes.
  resyncPeriod time.Duration
  
  // threadiness is the number of workers each to 
  // reconcile the desired resource i.e. state into
  // actual state.
  //
  // Defaults to 4.
  threadiness               int
  
  rateLimiter               workqueue.RateLimiter
  exponentialBackOffOnError bool
  threadiness               int
  
  // isLeaderElection flags if kubernetes leader
  // election will be used. It should be enabled
  // to avoid duplicate reconcile attempts
  isLeaderElection bool
  
  // leaderElectionNamespace is the kubernetes 
  // namespace in which to create the leader 
  // election object. Defaults to the same 
  // namespace in which the the controller runs
  leaderElectionNamespace string

  // LeaseDuration is the duration that non-leader
  // candidates will wait to force acquire leadership.
  // This is measured against time of last observed ack.
  // Defaults to 15 seconds.
  leaseDuration time.Duration

  // RenewDeadline is the duration that the acting
  // master will retry refreshing leadership before
  // giving up. Defaults to 10 seconds.
  renewDeadline time.Duration

  // RetryPeriod is the duration the LeaderElector 
  // clients should wait between tries of actions.
  // Defaults to 2 seconds.
  retryPeriod time.Duration

  hasRun     bool
  hasRunLock *sync.Mutex

  // TODO -- write appropriate comments
  //
  // Map UID -> *PVC with all claims that may be
 	// provisioned in the background.
  claimsInProgress sync.Map
}

type baseOption func(*Base)
type defaultOption func(*Base)

var defaults = []defaultOption{
  withDefaultResyncPeriod(),
	withDefaultThreadiness(),
}

func setDefaults(b *Base) {
  for _, option := range defaults {
    option(b)
  }
}

func withDefaultResyncPeriod() defaultOption {
  return func(b *Base) {
    if b.resyncPeriod == nil {
      b.resyncPeriod = 15min
    }
  }
}

func withDefaultThreadiness() defaultOption {
  return func(b *Base) {
    if b.threadiness == 0 {
      b.threadiness = 4
    }
  }
}

// New returns an instance of 
// base controller
type New() *Base {
  return &Base{}
}

type FromDefault() *Base {
  b := New()
  setDefaults(b)
  return b
}
```

```go
// pkg/controller/lib/v1alpha1/builder.go

// Builder helps in building an
// instance of base controller
type Builder struct {
  base *Base
}

// NewBuilder returns a new instance
// of builder
func NewBuilder() *Builder {
  return &Builder{base: New()}
}

// WithResyncPeriod is how often the controller relists
// watched resources. OnUpdate will be called even if 
// nothing has changed, meaning failed operations may be
// retried on every resyncPeriod regardless of whether it
// changed. Defaults to 15 minutes.
func (b *Builder) WithResyncPeriod(resyncPeriod time.Duration) *Builder {
  b.base.resyncPeriod = resyncPeriod
  return b
}

func (b *Builder) Build() *Base {
  setDefaults(b.base)
	return b.base
}
```

```go
// pkg/controller/pvc/v1alpha1/controller.go

import (
  controller "github.com/openebs/maya/pkg/controller/lib/v1alpha1"
)

// Controller watches and acts on 
// pvc resources
type Controller struct {
  *controller.Base
}

```
