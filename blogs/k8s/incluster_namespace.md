```go
// inClusterNamespace returns the namespace in which the pod is running in by checking
// the env var POD_NAMESPACE, then the file /var/run/secrets/kubernetes.io/serviceaccount/namespace.
// if neither returns a valid namespace, the "default" namespace is returned
func inClusterNamespace() string {
  if ns := os.Getenv("POD_NAMESPACE"); ns != "" {
    return ns
  }

  if data, err := ioutil.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/namespace"); err == nil {
    if ns := strings.TrimSpace(string(data)); len(ns) > 0 {
      return ns
    }
  }

  return "default"
}
```
