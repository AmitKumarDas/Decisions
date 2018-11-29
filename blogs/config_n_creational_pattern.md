
- refer - https://github.com/kubernetes/node-problem-detector/blob/master/pkg/systemlogmonitor/logwatchers/types/log_watcher.go
- refer - https://github.com/kubernetes/node-problem-detector/blob/master/pkg/systemlogmonitor/logwatchers/log_watchers.go
```go
package types

import (
	"k8s.io/node-problem-detector/pkg/systemlogmonitor/types"
)

// LogWatcher is the interface of a log watcher.
type LogWatcher interface {
	// Watch starts watching logs and returns logs via a channel.
  //
  // Comments
  // A channel can also be returned
  // Very good for non-blocking scenario
  // However if channel is returned then hooks should be available to clear it off
  //
  // Comments
  // Log comes from different package
  // In other words, log is treated as an API/types 
  // log watcher is treated as business logic
	Watch() (<-chan *types.Log, error)
  
	// Stop stops the log watcher. Resources open should be closed properly.
	Stop()
}

// WatcherConfig is the configuration of the log watcher.
//
// Comments
// No need to repeat Watcher
//
// Comments
// You can treat Config similar to specification. IMO if these are intended 
// to be provided as a whole by the caller. 
//
// However, there are cases when some fields of these structures can be set 
// with defaults, while setting other fields are dependent on some conditions.
// In such cases, much more thought is required to build the spec itself.
type WatcherConfig struct {
	// Plugin is the name of plugin which is currently used.
	// Currently supported: filelog, journald, kmsg.
	Plugin string `json:"plugin,omitempty"`
	
  // PluginConfig is a key/value configuration of a plugin. Valid configurations
	// are defined in different log watcher plugin.
	PluginConfig map[string]string `json:"pluginConfig,omitempty"`
	
  // LogPath is the path to the log
	LogPath string `json:"logPath,omitempty"`
	
  // Lookback is the time log watcher looks up
	Lookback string `json:"lookback,omitempty"`
	
  // Delay is the time duration log watcher delays after node boot time. This is
	// useful when the log watcher needs to wait for some time until the node
	// becomes stable.
	Delay string `json:"delay,omitempty"`
}

// WatcherCreateFunc is the create function of a log watcher.
//
// Comments
// Often wondered how to initialize the interface instance in a typed way
// You can go for factory interface or just this function type
type WatcherCreateFunc func(WatcherConfig) LogWatcher

// createFuncs is a table of createFuncs for all supported log watchers.
//
// Comments
// I call this as registrar
var createFuncs = map[string]types.WatcherCreateFunc{}

// registerLogWatcher registers a createFunc for a log watcher.
func registerLogWatcher(name string, create types.WatcherCreateFunc) {
	createFuncs[name] = create
}

// GetLogWatcherOrDie get a log watcher based on the passed in configuration.
// The function panics when encounters an error.
//
// Comments
// This naming is good as it says any error will panic
//
// Comments
// Donot get confused with the way config is used to extract plugin as 
// well as create an instance of the same
//
// Comments
// This is kind of singleton, but we really donot need to do init() or 
func GetLogWatcherOrDie(config types.WatcherConfig) types.LogWatcher {
	create, ok := createFuncs[config.Plugin]
	if !ok {
		glog.Fatalf("No create function found for plugin %q", config.Plugin)
	}
	glog.Infof("Use log watcher of plugin %q", config.Plugin)
	return create(config)
}
```
