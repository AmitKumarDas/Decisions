### Meta Info
- Version - 1
- Last Updated On - 14-Mar-2019

### Avoid these stuff as seen in RunTasks
```yaml
post: |-
  '{{- jsonpath .JsonResult "{.Stdout}" | trim | saveAs "remoteexec.Stdout" .TaskResult | noop -}}'
  '{{- jsonpath .JsonResult "{.Stderr}" | trim | saveAs "remoteexec.Stderr" .TaskResult | noop -}}'
```

### Solution
