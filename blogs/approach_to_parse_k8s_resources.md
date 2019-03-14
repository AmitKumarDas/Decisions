### Meta Info
- Version - 1
- Last Updated On - 14-Mar-2019

### Avoid these stuff as seen in RunTasks
```yaml
post: |-
  '{{- jsonpath .JsonResult "{.Stdout}" | trim | saveAs "remoteexec.Stdout" .TaskResult | noop -}}'
  '{{- jsonpath .JsonResult "{.Stderr}" | trim | saveAs "remoteexec.Stderr" .TaskResult | noop -}}'
```

### Simplify these further
```yaml
post: |
  {{- $cmd := resize cstor volume -}}
  {{- $cmd = $cmd | withoption "ip" .TaskResult.resizelistsvc.clusterIP -}}
  {{- $cmd = $cmd | withoption "volname" .Volume.owner -}}
  {{- $cmd = $cmd | withoption "capacity" .Volume.capacity -}}
  {{- run $cmd | saveas "resizecstorvolume" .TaskResult -}}
  {{- $err := .TaskResult.resizecstorvolume.error | default "" | toString -}}
  {{- $err | empty | not | verifyErr $err | saveIf "resizecstorvolume.verifyErr" .TaskResult | noop -}}
```

```yaml
post: |
  {{- $err := .TaskResult.createpatchcv.error | default "" | toString -}}
  {{- $err | empty | not | verifyErr $err | saveIf "createpatchcv.verifyErr" .TaskResult | noop -}}
```

### Solution
