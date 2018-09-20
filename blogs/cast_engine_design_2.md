### CASTemplate Engine Design - Chapter 2
Exposing CASTemplate along with RunTasks enables exposing logic in a declarative fashion. It helped the team in numerous ways. One of the typical use of these templates was the ability to embed logic into RunTasks and solve the problem statement
externally without getting into the project's in-tree repo. 

### Things that went good
```yaml
- Template Functions
- Pipes
```

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
The design should be accomodating to accept specs from users and convert these specs to logic. The specs will not be logic
but will be pseudo-logic i.e. some english like notations. These specs should be provided to feature specific engine. The 
engine will decide the provider and runner based on these specs.

High Level Design:
- UseCase: Volume Provisioning
  - spec - pkg/apis/openebs.io/v1alpha1/casvolume.go may have `type CASVolume struct`
  - engine - pkg/engine/v1alpha1/cas_template.go may have `type CASTemplateEngine struct`
  - engine - pkg/engine/v1alpha1/run_task.go may have `type RunTaskEngine struct`
  - runner - pkg/task/v1alpha1/txt_template.go may have `type TxtTemplateRunner struct`
  - runner - pkg/task/v1alpha1/run_command.go may have `type RunCommandRunner struct`

### Rough Work
This rough work lists down all sorts of templating possibilities. However, only few have been selected by me. This will get
refined further based on feedbacks, experiences & my brain's biasedness.

#### Select Clause
- [ ] `select all | create kubernetes service | spec $yaml | totemplate .Values .Config | run`
- [ ] `select all | text template .Values $doc | tounstruct | run | saveas "123" .Values`

#### Where Clause

#### Join multiple queries
- [ ] `select name, ip | get k8s svc | where label eq abc | join list k8s pod | where 'label' 'eq' 'all'`

#### Error Handling

#### Predicates
- [ ] `create k8s service | spec $yaml | totemplate .Values | set "namespace" "value" isnamespace | run`
- [ ] `create k8s service | spec $yaml | totemplate .Values | set "namespace" "value" | run`

#### Text Template as a template function
- [ ] `$doc | exec template . Values | run`
- [ ] `$doc | exec text template | data . Values | run`
- [x] `$doc | text template | data . Values | run`
- [ ] `create kubernetes service | specs $doc | txttemplate . Values | run`
- [ ] `create k8s svc | spec $doc | txttemplate .Volume .Config  | run`
- [x] `select name, ip | create k8s service | spec $doc | totemplate .Volume .Config | run`

#### Spec
```yaml
- select:
    name namespace 
  get:
    kubernetes service 
  where:
    name eq John
```