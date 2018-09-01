### CASTemplate Engine Design - Chapter 2
Exposing CASTemplate along with RunTasks enables exposing logic in a declarative fashion. It helped the team in numerous ways.
One of the typical use of these templates was the ability to embed logic into RunTasks and solve the problem statement without
taking a dip into compilation and so on. There as other rosy aspects with the use of template. However, other side of the coin
has a lot of pain points. In this design which will form the basis of next CAST engine, I will list down my thoughts of what
should go in to make the engine better.

### Things to improve
- Avoid as much as possible the need to write specific CAST engines
  - Volume, Clone, Pool, Upgrade and so on needs to be written to influence generic CAST engine
- Ability to Unit Test a RunTask
  - This helps in verifying the validity of the RunTask with various data options
- Avoid as much as possible the need to write conditional logic
- Avoid as much as possible the need to write meta info in RunTasks
- Avoid get from specific paths and setting at specific paths
  - User/Dev/Operator does not care about this
- Errors are hard to debug
  - Improvements:
    - Infos as a list of successful runtasks
    - Warns as a list of warning runtasks
    - Skips as a list of skipped runtasks
    - Errors as a list of errored runtasks
- Lots of RunTasks are needed
- Executing RunTasks conditionally
- Streaming based template functions
  - Map invocation against a list
  - Filter invocation against a list
  - Use of Predicate
- Template Functions without the use of method calls
- All validations should be triggered before invoking all/any runtasks
  - data related validations
  - template related validations
  - old vs. new templating style
- Use of Predicates -- same as selectors
  - isPresent,
  - isNotEmpty ".spec.disks",
  - isLbl, isLblAll, isLblAny,
  - isAnn, isAnnAll, isAnnAny,
  - isNode,
  - isTaint,
  - isNamespaced,
- Might need to wrap existing template functions that work for unstruct instances
  - new template function names will need to be thought of
- Ability to curry the existing template functions & give good names to these curried functions

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
  - {{- cast "step1" | kload .cast.yamls[0] | kcreate | select ".spec.ip" ".spec.uid" ".spec.name" | run -}}
  - {{- cast "step2" | kget "pod" .Config.name | ns .Volume.runNS | select ".spec.ip" ".spec.node" ".spec.status" | runas ".poddetails" -}}
  - {{- cast "step3" | klist "pods" | ns "abc" "def" "def" | select ".spec.name" | where ".spec.status" "eq" "running" | where ".spec.label" "haskey" "abc" | where ".spec.label" "hasval" "def" | and | runas "podlistdetails" -}}
  - {{- cast "step4" | klist "deploy" | ns "abc" | select ".spec.name" | where ".spec.labels" "has" "key=val" | where ".spec.labels" "has" "key1=val1" | or | run -}}
  - {{- cast "step5" | hdelete | url $url | select all | run -}}
  onErr:
```
- Accessing cast related _Template Values_:
```yaml
# rollbacks, rollbackresps, fallbacks, errors, infos, warns are debugging purposes
values:
  - .cast.reqs.
  - .cast.resps.[].poddetails.
  - .cast.resps.[].podlistdetails.
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
  Unstruct  *unstructured.Unstructured
  Select    []string
  Where     []string
  Errors    []error
  Warns     []string
}

// castresp is the successful response of executing a cast request
type castresp  *unstructured.Unstructured
type casterror *castresp
type castwarn  *castresp

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
