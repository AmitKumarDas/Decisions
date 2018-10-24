How often we felt the need for a one time registration while writing our logic.

Recently I happened to read code from client-go where the use of sync.Once made this feature implementation really simple.

```go
// Package metrics provides abstractions for registering which metrics
// to record.
package metrics

import (
	"net/url"
	"sync"
	"time"
)

// This has to be package level global variable as we do not want any
// other instances of this object. In other words, as long as this 
// object is not created more than once, the action is triggered once
// in the lifetime of this application.
//
// DOUBT - Why is this not declared as a pointer?
var registerMetrics sync.Once

// LatencyMetric observes client latency partitioned by verb and url.
type LatencyMetric interface {
	Observe(verb string, u url.URL, latency time.Duration)
}

// ResultMetric counts response codes partitioned by method and host.
type ResultMetric interface {
	Increment(code string, method string, host string)
}

var (
	// RequestLatency is the latency metric that rest clients will update.
	RequestLatency LatencyMetric = noopLatency{}
	// RequestResult is the result metric that rest clients will update.
	RequestResult ResultMetric = noopResult{}
)

// Register registers metrics for the rest client to use. This can
// only be called once.
func Register(lm LatencyMetric, rm ResultMetric) {
	registerMetrics.Do(func() {
		RequestLatency = lm
		RequestResult = rm
	})
}

type noopLatency struct{}

func (noopLatency) Observe(string, url.URL, time.Duration) {}

type noopResult struct{}

func (noopResult) Increment(string, string, string) {}
```

reference - https://github.com/kubernetes/sample-controller/blob/master/vendor/k8s.io/client-go/tools/metrics/metrics.go
