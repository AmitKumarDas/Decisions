### References
- https://github.com/bitrise-io/stepman/blob/master/stepman/util.go
- https://github.com/bitrise-io/go-utils/blob/master/retry/retry.go

```go
// usage
err := retry.Times(2).Wait(3 * time.Second).Try(func(attempt uint) error {
  return command.DownloadAndUnZIP(
    downloadLocation.Src, 
    stepPth,
  )
})

```

```go
// actual logic:
// - functional
// - fluent API
// - method chaining
// - flatten your code

package retry

import (
  "fmt"
  "time"
)

// Action ...
type Action func(attempt uint) error

// Model ...
type Model struct {
  retry    uint
  waitTime time.Duration
}

// Times ...
func Times(retry uint) *Model {
  Model := Model{}
  return Model.Times(retry)
}

// Times ...
func (Model *Model) Times(retry uint) *Model {
  Model.retry = retry
  return Model
}

// Wait ...
func Wait(waitTime time.Duration) *Model {
  Model := Model{}
  return Model.Wait(waitTime)
}

// Wait ...
func (Model *Model) Wait(waitTime time.Duration) *Model {
  Model.waitTime = waitTime
  return Model
}

// Try ...
func (Model Model) Try(action Action) error {
  if action == nil {
    return fmt.Errorf("no action specified")
  }

  var err error
  for attempt := uint(0); (0 == attempt || nil != err) && attempt <= Model.retry; attempt++ {
    if attempt > 0 && Model.waitTime > 0 {
      time.Sleep(Model.waitTime)
    }

    err = action(attempt)
  }

  return err
}
```
