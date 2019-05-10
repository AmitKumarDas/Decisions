refer - https://github.com/kubermatic/kubeone/blob/master/pkg/task/task.go

```go
package task

import (
	"time"

	"github.com/kubermatic/kubeone/pkg/util"

	"k8s.io/apimachinery/pkg/util/wait"
)

// defaultRetryBackoff is backoff with with duration of 5 seconds and factor of 2.0
func defaultRetryBackoff(retries int) wait.Backoff {
	return wait.Backoff{
		Steps:    retries,
		Duration: 5 * time.Second,
		Factor:   2.0,
	}
}

// Task is a runnable task
type Task struct {
	Fn      func(*util.Context) error
	ErrMsg  string
	Retries int
}

// Run runs a task
func (t *Task) Run(ctx *util.Context) error {
	if t.Retries == 0 {
		t.Retries = 1
	}
	backoff := defaultRetryBackoff(t.Retries)

	var lastError error
	err := wait.ExponentialBackoff(backoff, func() (bool, error) {
		lastError = t.Fn(ctx)
		if lastError != nil {
			ctx.Logger.Warn("Task failed, retrying…")
			return false, nil
		}
		return true, nil
	})
	if err == wait.ErrWaitTimeout {
		err = lastError
	}
	return err
}
```
