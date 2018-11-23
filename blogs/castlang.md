## CASTLang

This is the next version of CASTemplate. This should provide all the benefits that CASTemplate provides as well make scripting
easier.

### CASTemplate - Good Parts - 1
- variable declaration & definition via go template functions before the task yaml
  - e.g. `{{- $name := .TaskResult.name | default "cool" -}}`
  - e.g. `{{- $resourceLimitsVal := fromYaml .Config.PoolResourceLimits.value -}}`
- task yaml makes use of these variables as go templating


#### Low Level Design - Good Parts - 1
- variable declaration, definition, transformation, autosave
```yaml
kind: RunTask
spec:
  template:
  - pod: |
      {{- $name := .TaskResult.name | default "cool" -}}
      kind: Pod
      apiVersion: v1
      metadata:
        name: $name
```

```go
type runtask struct {
  meta ObjectMeta
  spec RunTaskSpec `json:"spec"`
}

type RunTaskSpec struct {
  template map[string]string `json:"template"`
}
```
