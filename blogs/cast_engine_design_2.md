### CASTemplate Engine Design - Chapter 2
Exposing CASTemplate along with RunTasks enables exposing logic in a declarative fashion. It helped the team in numerous ways. One of the typical use of these templates was the ability to embed logic into RunTasks and solve the problem statement
externally withoyt getting into the project's in-tree repo. 

However, there are still numerous pain points that need to be addressed in cast. In this design, which will form the basis of 
next CAST engine, I will list down my thoughts of what should go in to make **cast** better.

### Things to improve
```yaml
- Avoid the need to write specific CAST engines
- Ability to Unit Test a RunTask
- Reduce the need for programming in templates
- Avoid verbosity in RunTasks
- Make the template intuitive & less learning curve
- Errors are hard to debug
  - Improvements:
    - Infos as a list of successful runtasks
    - Warns as a list of warning runtasks
    - Skips as a list of skipped runtasks
    - Errors as a list of errored runtasks
- Avoid the need for multiple RunTasks
- Ability to execute RunTasks conditionally
- Streaming based template functions
  - Map invocation against a list
  - Filter invocation against a list
  - Use of Predicate
- Template Functions without the use of method calls
- All validations should be triggered before invoking all/any runtasks
  - data related validations
  - template related validations
  - old vs. new templating style
```

### Things that went good
- Template Functions
- Pipes

### Things to Add
- Generic HTTP API
- Below is a sample yaml snippet:
```yaml
post:
  - {{- $url = http://{{.Volume.ControllerIP}}:5104/volume/delete/{{.Volume.owner}} -}}
  - {{- $url | hpost | struct | saveAs .TaskResult.jivadel.resp -}}
```

### New Design - Day 1 Thinking
- Values as feeds from user, runtime, engine, etc
  - Can be .Volume, .Config or both or more based on what is fed to CAS engine
- Stores as storage during execution of one or more runtasks
- Docs stores the yamls as unstruct instances indexed with name of unstruct
  - i.e. `[]*unstruct`
- Values, Stores & Docs are stored in a map[string]interface{} & is provided to template
- Keep the template function generic to start with
```yaml
yamls:
  - yaml: |
      kind: cool
      apiVersion: v1
      metadata:
        name: abc-123
runs:
  - if: 
    run: |
      - {{- .Docs.abc-123 | New | noop -}}
      - {{- $isFalse := fn eq "false" -}}
      - {{- $isLblPresent := .Stores.abc-123 | $isFalse ".spec.metadata.labels" -}}
      - {{- $myNewRTObjList := $myRTObjList | Map | Filter p | Any p -}}
      - {{- $myRTObjList | Map | Filter p | Any p | Store id123 -}}
      - {{- $aMap := .Stores.id123 | Items -}}
      - {{- $myRTObjList | Select ".spec.abc" ".spec.def" ".spec.xyz" -}}
      - {{- $myRTObj | Select ".spec.abc" ".spec.def" ".spec.xyz" -}}
    onErr: 
onErr:
```

### New Design - Day 2 Thinking
- Ability to execute a runtask programatically
  - Decouples from dependency on yaml
  - Decouples from k8s api server, custom resource, etcd
  - More abstracted runtask runner engine
  - Ability to UnitTest on kubernetes cluster
- Ability to execute a castemplate programatically
  - Same as above advantages
- Better handling of flow of runtask steps
  - Aborts the runtask automatically based on errors or injected errors
  - Run conditionally a runtask step based on condition supplied
  - Retry a runtask step based on particular error type
- **Self manageable functions without need of engine**
- Sample Runtask yaml:
```yaml
metadata:
  name: cstor-volume-create
spec:
  yamls:
  runs:
  - {{- cast "step1" | ymlmap "name eq my-pod" | reqdata | kcreate | run -}}
  - {{- cast "step2" | kget "pod" .Config.name | ns .Volume.runNS | select ".spec.ip" ".spec.stat" | runas ".poddetails" -}}
  - {{- cast "step3" | respmap "name eq poddetails" | respdata | select ".spec.ip" ".spec.uid" ".spec.name" | run -}}
  - {{- cast "step4" | klist "pods" | ns "abc" "def" "def" | select ".spec.name" | whereall ".spec.status eq running" ".spec.label haskey abc" | runas "podlistdetails" -}}
  - {{- cast "step5" | hdelete | url $url | select all | run -}}
  onErr:
  - {{- cast "rbstep1" | kdelete "pod" .Config.name | ns .Volume.runNS | run -}}
```
- Accessing cast related _Template Values_:
```yaml
# rollbacks, rollbackresps, fallbacks, errors, infos, warns are debugging purposes
values:
  - .cast.yamls.
  - .cast.reqs.
  - .cast.resps.[name==poddetails].
  - .cast.resps.[name==podlistdetails].
  - .cast.rollbacks.
  - .cast.rollbackresps.
  - .cast.fallbacks.
```
- Sample go code:
```go
type CastAction string

const (
  KLoad   CastAction = "kload"
  HDelete CastAction = "hdelete"
  KGet    CastAction = "kget"
)

// cast represents a runtask step request
//
// A request can be cast to anything
type cast struct {
  StepID    string
  Action    CastAction
  Namespace string
  Reqdata   *unstructured.Unstructured
  Resdata   []*unstructured.Unstructured
  Select    []string
  Where     []CastCondition
  WhereOp   CastOperator
  Errors    []error
  Warns     []string
}

// castresp is the successful response of executing a cast request
type castresp  *unstructured.Unstructured

// kcast is kubernetes based cast structure
type kcast *cast
// hcast is http based cast structure
type hcast *cast
// gcast is grpc based cast structure
type gcast *cast
// castvalues are used during templating
type castvalues map[string]interface{}
```
- Sample go code for engine:
```go
type CasTemplateRunner struct {}
type RunTaskRunner struct {}
```
- Sample go code for message:
```go
// pkg/msg/v1alpha1/msg.go
type msg struct {
  type msgType
  desc string
  err  error
}
func (m msg) String() {}

type msgs []*msg

type (m msgs) Log(l func(...interface{})){
  for _, msg := range m {
    l(msg.String())
  }
}

type (m msgs) LogAny(l func(...interface{}), pl []msgPredicate)
type (m msgs) LogAll(l func(...interface{}), pl []msgPredicate)

func (m msgs) AddInfo() {}
func (m msgs) AddWarn() {}
func (m msgs) AddSkip() {}
func (m msgs) AddError() {}
func (m msgs) Filter(p msgPredicate) (f msgs){}
func (m msgs) Infos() (f msgs){
  return m.Filter(infoPredicate)
}
```

### New Design - Day 3 Thinking
```
rt 123 | get volume jiva | where name == abc
rt 123 | get jiva volume | where name == abc
rt 123 | select name namespace | from jiva volume | where namespace == default | where label == app=jiva
rt 123 | delete jiva volume | where name == default | where ip == $ip (edited)
rt 123 | delete jiva volume | whereany 'name == abc' 'namespace == default'
```

```
rt 123 | select all | delete jiva volume | where 'name' 'eq' $name | where 'namespace' 'eq' $namespace | run
rt 123 | select 'all' | delete jiva volume | where 'name' 'eq' $name | where 'namespace' 'eq' $namespace | runif $count 'eq' 0
rt 123 | select 'all' | delete jiva volume | where 'name' 'eq' $name | runif eq $count 0
```

```
rt 123 | select 'all' | delete jiva volume | where 'name' '=' $name | runif eq $count 0
rt 123 | select 'all' | from jiva volumes  | where 'name' '=' $name | runif eq $count 0
rt 123 | select 'all' | from jiva volumes | where 'namespace' 'in' $names | runif eq $ciunt 0
```

```
rt 123 | select 'serviceip' | create k8s service | where 'spec' '=' $yaml | runif eq $exists false
rt 123 | select 'serviceip' | get k8s service | where 'namespace' '=' 'default' | runif eq $err false
rt 123 | select 'serviceip' | from k8s service | where 'name' '=' 'okie' | runif eq $err false
```

