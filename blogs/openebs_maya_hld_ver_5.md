### OpenEBS Maya - Version 5
Exposing CASTemplate along with RunTasks enables exposing **logic in a declarative manner**. It helped the team in numerous
ways. One of the typical use of these templates was the ability to embed logic into RunTasks and solve the problem statement
externally without writing code.

### Is Templating & Template Functions Good?
```yaml
- Template functions can be made similar to invoking CLI commands
  - The template will look like a scripted copy of bunch of CLI commands on terminals
- Template functions can be much more powerful when
  - We can join multiple operations
  - We can operate these functions selectively
  - We can keep piping the results to new functions with ease
- Demerit
  - The execution control is entirely with template runner
  - Erroring out of this template runner is not elegant
  - Debugging is pathetic
```

### The way forward
```yaml
- Avoid the need to write specific engines that wrap template runner
- Ability to Unit Test a RunTask
- Reduce the need for programming in templates
  - In other words convert the template into a spec
- Avoid verbosity in RunTasks
- Make the template intuitive & simple to learn
- Improve the debuggability aspect
- Avoid the need for multiple RunTasks
- Ability to execute RunTasks conditionally
- Ability to add skip, validate, error injection & non-functional logic
```

### Next Steps
Maya design will expose two new kinds of config in addition to template based config. 
They are explained below:
- human config - human understood and are inline with the feature
- raw config - program understood that implements a feature

A raw config looks like a pseudo-code. A human based config is typically implemented (or transformed & then run) as a raw
config.

Package Design:
```
  - spec - pkg/apis/openebs.io/v1alpha1/
  - engine - pkg/engine/v1alpha1/
  - runner - pkg/task/v1alpha1/ vs. pkg/runner/v1alpha1
  - resource - pkg/resource/v1alpha1/
  - client - pkg/client/k8s/v1alpha1/
  - http client - pkg/client/http/v1alpha1/
  - template - pkg/template/v1apha1/
  - install - pkg/install/v1alpha1/
  - volume - pkg/volume/v1alpha1/
  - snapshot - pkg/snapshot/v1alpha1/
  - clone - pkg/clone/v1alpha1/
  - pool - pkg/pool/v1alpha1/
```

Low Level Design:
- engine
```go
type Engine interface {}

type CASTemplate struct {
  data        interface{}
  provider    CASTemplateProvider
  rtEngine    RunTaskEngine
}

type RunTask struct {
  config    map[string]interface{}
  data      interface{}
  provider  RunTaskProvider
  runner    CommandRunner
}
```
- runner
```go
type CommandRunner interface {}

type RunCommand struct {}
type TxtTemplateRunner struct {}
type RunCommandRunner struct {}
```
- resource
```go
type ResourceProvider interface {}
type ResourceGetter interface {}
type ResourceUpdater interface{}

type K8sResource struct {}
type CASTemplate struct {}
type RunTask struct {}
type ConfigMap struct{}
```
- engine helpers
```go
type ResourceNameFetcher interface {}

type ResourceName struct {}
type ResourceNameBySC struct {}
type ResourceNameByENV struct {}
```

### Template Config - Futuristic Version i.e. similar to **MySQL Query Language**
This rough work lists down all sorts of templating possibilities for a `template config`. 

_NOTE: This will get refined further based on feedbacks, experiences & my brain's biasedness._
_NOTE: The core structures/code will be reused to implement **Human Config** & **Raw Config**._

#### Select Clause
- [ ] `select all | create kubernetes service | spec $yaml | totemplate .Values .Config | run`
- [ ] `select all | text template .Values $doc | tounstruct | run | saveas "123" .Values`

#### Where Clause
- [ ] `select "name" "ip" | get k8s service | where "name" $name | run`
- [ ] `select "name" "ip" | get k8s service | where "name" "ne" $name | run`
- [ ] `select "name" "ip" | delete k8s service | where "name" $name | run`
- [ ] `select "name" "ip" | list k8s service | where "name" "ne" $name | run`

#### Whereop Clause
- [ ] `select "name" "ip" | list k8s service | where "name" $name | where "label" $app | whereop "any" | run`
- [ ] `select "name" "ip" | list k8s service | where "name" $name | where "label" $app | whereop "all" | run`
- [ ] `select "name" "ip" | list k8s service | where "name" $name | where "label" $app | whereop "atleasttwo" | run`

#### When Clause
- [ ] `select "name" "ip" | get k8s service | where "name" "ne" $name | when eq $willrun "true" | run`

#### Join multiple queries
- [ ] `select name, ip | get k8s svc | where label eq abc | join list k8s pod | where "label" "all"`

#### Error Handling Clause

#### Set Clause
- [ ] `create k8s service | spec $yaml | totemplate .Values | set "namespace" $ns | run`
- [ ] `create jiva snapshot | set "spec" $yaml | txttemplate .Values | run`
- [ ] `create jiva snapshot | set "name" $name | set "capacity" "2G" | run`
- [ ] `create k8s pod | set spec $spec | run`
- [ ] `creates k8s pod | set image $image | run`
- [ ] `create k8s jiva replica | set image $image | run`
- [ ] `delete jiva data | run`
- [ ] `create jiva target | set image $image | run`
- [ ] `create jiva target | set "container.path" $container | run`
- [ ] `create jiva target | set "spec" $spec | run`
- [ ] `create jiva target | set "doc" $doc | run`
- [ ] `create jiva target | set "yaml" $yaml | run`
- [ ] `create k8s deploy | set "yaml" $yaml | txttemplate .Values .Config | run`
- [ ] `create k8s deploy | set "yaml" $yaml | set "image" $image | txttemplate .Values .Config | run`

#### Set Conditionally
- [ ] `create k8s service | spec $yaml | totemplate .Values | setif isnamespace "namespace" "value" | run`


#### Text Template as a template function
- [ ] `$doc | exec template . Values | run`
- [ ] `$doc | exec text template | data . Values | run`
- [x] `$doc | text template | data . Values | run`
- [ ] `create kubernetes service | specs $doc | txttemplate . Values | run`
- [ ] `create k8s svc | spec $doc | txttemplate .Volume .Config  | run`
- [x] `select name, ip | create k8s service | spec $doc | totemplate .Volume .Config | run`


#### Sample Template Config
```yaml
apiVersion: openebs.io/v1alpha1
kind: RunTask
metadata:
  name: jiva-volume-delete-deletedata-default-0.7.0
spec:
  meta: |
    id: deljivadata
    kind: Command
  post: |
    {{- $base := print "http://" .TaskResult.deletelistsvc.ips ":9501" -}}
    {{- delete jiva volume | url $base | name .Volume.owner | run | saveas "deljivadata" .TaskResult -}}
    {{- $err := toString .TaskResult.deljivadata.error -}}
    {{- $err | empty | not | verifyErr $err | saveIf "deljivadata.verifyErr" .TaskResult | noop -}}
```

### Raw Config
```yaml
- select:
    name namespace 
  get:
    kubernetes service 
  where:
    name eq John
```

### Human Config
```yaml
```
