### Meta Info
- Version - 2
- Last Updated On - 20-Mar-2019

### Motivation
What things in software are most common? What is the most simplistic version of any software product. IMO it is its CLI.
Somehow I tend to agree that given a CLI that autocompletes & provides me the help at the terminal itself is the coolest
thing and at the same time the most most simplest thing that has ever happened to software. Recently same can be told for
a declarative specification that can be saved, versioned and injected into workflows.

The motivation of this article is to find a balance between action oriented CLI to specification driven yamls.

### Behaviours == Actions == Verbs
I shall take the example of openebs storage as the software. What can be various actions related to openebs? We can think
of below:
- create
- delete
- get
- identify
- list
- resize
- backup
- snapshot
- etc...

### Nouns == Entities
Lets think of various nouns used in openebs:
- cstor volume
- cstor replica
- cstor target
- cstor pool
- jiva volume
- jiva target
- jiva pool
- etc...

### Mix them up == CLI
- create cstor volume
- delete cstor pool
- resize cstor volume
- backup cstor volume
- get cstor volume
- etc...

### Devil == CLI in YAML
This is considered as the devil by most fellow programmers. However, IMO if it helps users to get the requirement done
in a static compiled language and can be saved & versioned then why not. Being declarative is not a option but a need.
In other words, there would not have any `github actions` to begin with

### Prototype
```yaml
spec:
  create this pod: 
    construct:
      cmd: create kube pod
      with:
        namespace: default
        inCluster: false
        output: ".spec.hostName"
        yaml: blabh blah blah
      if:
        isRunning:
        isLabel: "app=openebs"
    options:
      saveas: createit
      audit: true
      runif: noErrors
  check this pod: 
    construct:
      cmd: get kube pod
      with:
        namespace: default
        name: poddy
        inCluster: false
        output: ".spec.hostName"
      if:
        isRunning:
        isLabel: "app=openebs"
    options:
      saveas: poddy
      audit: true
      runif: noErrors
```

### Prototype - 2
- Yaml when unmarshaled gets the job done
- Build and Run
- Build options and run options
- When marshaled returns the result of cmd in JSON bytes or error
- Can be formed as a list of cmds
- Executor engine can unmarshal the runtask
  - Run a go template of first cmd in the array
  - Then unmarshal and save the result of cmd into global values
  - Next cmd will be executed if no error occurred to previous cmd.

```yaml
kind: RunTask
spec:
  desc: do my work please
  runs:
  - desc: get me a list of pods
    id: pod123
    cmd:
      type: podList
      options: 
      - name: withNamespace 
        args: default
      - name: inCluster
      checks:
      - name: isRunning
      - name: isNamespace
        args: default
      output: 
        name: list
  - desc: display the output
    id: display101
    cmd:
      type: display
      options: 
      - name: withYaml 
        args: |
          kind: PodList
          items:
          - id: {{ .pod123.id[0] }}
```
