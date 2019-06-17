```go
// inClusterNamespace returns the namespace in which the pod
// is running in by checking the env var POD_NAMESPACE, then
// the below file:
// /var/run/secrets/kubernetes.io/serviceaccount/namespace
//
// NOTE:
//  if neither returns a valid namespace, the "default" 
// namespace is returned
func inClusterNamespace() string {
  ns := os.Getenv("POD_NAMESPACE")
  if ns != "" {
    return ns
  }

  data, err := ioutil.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/namespace")
  if err == nil {
    ns := strings.TrimSpace(string(data))
    if len(ns) > 0 {
      return ns
    }
  }

  return "default"
}
```
