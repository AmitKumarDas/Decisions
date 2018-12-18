## CASTLang

CASTLang is the _code name_ for next version of CASTemplate. This should provide all the benefits that CASTemplate provides 
as well help building **kubernetes workflows easier**. This version concentrates on RunTask. A RunTask can get executed via 
CASTemplate runner or via a new kubernetes controller. Care has been taken **NOT to modify CASTemplate** as it is user 
facing and hence is prone to breaking its users. This attempt to improvise CASTemplate code named **CASTLang** tries to make 
RunTask as independent and devops friendly than its earlier version.

#### Low Level Design - v2.0
- The design boils down to:
  - Variable Declaration, Definition, Options, Predicates Actions & AutoSave

- Main structure
```yaml
kind: RunTask
spec:
  config:
  runs:
status:
```

- Config - Example 1
```yaml
kind: RunTask
spec:
  config: # map[string]interface{}
  - volume:
    - name: myvol
      namespace: default
      taints:
      - key: type
        operator: In
        value: node1
  - pool:
    - name: mypool
```

- Config - Example 2
```yaml
kind: RunTask
spec:
  config:
  - spec: |
      kind: Pod
      apiVersion: v1
      metadata:
        name: Hulla
```

- Config - Example 3
```yaml
kind: RunTask
spec:
  config:
  - values: 
    - name: poddy
      namespace: default
  - data:
    - name: herdy
      namespace: openebs
  - podtpl: |
    {{- $name := .config.values.name | default "cool" -}}
    {{- $ns := .config.values.namespace | default "default" -}}
    kind: Pod
    apiVersion: v1
    metadata:
      name: $name
      namespace: $ns
```

- Run - Example 1
```yaml
kind: RunTask
spec:
  runs:
  - id: 101
    name: # optional; set to id value if not set
    action: list
    kind: PodList
    options:
    - func: labelSelector app=jiva
    - func: namespaceSelector default,openebs
```

- Run - Example 2
```yaml
kind: RunTask
spec:
  runs:
  - id: 101
    name: # optional; set to id value if not set
    action: create
    kind: Pod
    options:
    - func: template ${@.config.podtpl} ${@.config.values}
```

- Run - Example 3
```yaml
kind: RunTask
spec:
  runs:
  - id: 101
    name: # optional; set to id value if not set
    action: create
    kind: Pod
    options:
    - func: spec ${@.config.spec}
```

- Run - Example 4
```yaml
kind: RunTask
spec:
  config:
  - tolerations:
    - key: "key"
      operator: "Equal"
      value: "value"
      effect: "NoSchedule"
  - ns: openebs
  - k8s11: v1.11.0
  - spec: abc
  runs:
  - id: 101
    name: # optional; set to id value if not set
    action: create
    kind: Pod
    options:
    - func: spec ${@.config.spec}
    - func: addTolerations ${@.config.tolerations}
      conditions:
      - isCAST
      - isNamespace ${@.config.ns}
    conditions:
    - isVersion ${@.config.k8s11}
```

### SpinOffs - CASTemplate
- CASTemplate can be a pure templating solution that generates one or more yamls(s)
```yaml
kind: CASTemplate
spec:
  defaultConfig:
  - values: 
    - name: poddy
      namespace: default
  - podTemplate: |
    {{- $name := .config.values.name | default "cool" -}}
    {{- $ns := .config.values.namespace | default "default" -}}
    kind: Pod
    apiVersion: v1
    metadata:
      name: $name
      namespace: $ns
  templates:
  - id: 101
    name: # optional; set to id value if not set; used for descriptive stuff
    result: # hidden when viewed as yaml; set after executing this action
    kind: Pod
    options:
    # make use of .config & not .defaultConfig
    - func: template ${@.config.podTemplate} ${@.config.values}
status: # hidden while viewing the yaml; used in code
```

### SpinOffs - TestTask
- One can extend RunTask to meet their specific requirement
- I shall explain how `TestTask` extends from RunTask
- Most of stuff remains same baring `expect` which gets introduced as a new field
- Assumptions:
  - Depends on consuming Kubernetes APIs
  - Based on workflow or pipelining of tasks
  - Do not name this as `KubeTest` as we want this to be more task/worker oriented.
  - NOTE: There is a separate proposal for [KubeTest](https://github.com/AmitKumarDas/Decisions/blob/master/blogs/openebs_operator_design.md)

#### My Conviction on TestTask
More than the technical aspects that I mentioned above, I view this as a compelling tool to radical customer-centrism. 
This tool lets **dev, test, users, doc authors and so on** as a quick and easy approach to try, test, and verify their 
expectations. For example, an user wants to verify if upgrade from version **X** to version **Y** is successful and then 
test if provisioning works as expected in the new version. With _TestTask_, user needs to articulate the requirements as specs and on execution of TestTask verifies the truth of the same.

```yaml
kind: TestTask
spec:
  given:
  when:
  then:
status:
```

```yaml
kind: TestTask
spec:
  given:
  when:
    - id: 101
      name: # optional; set to id value if not set
      action: list
      kind: PodList
      labelSelector: app=jiva,org=openebs
      expect:
        match:
        - status == Online
        - kind == PodList
        - namespace In default,openebs
        - labels == app=jiva,org=openebs
  then:
    - id: 201
      action: get
      kind: Pod
      name: cstor
      expect:
        retry:
        match:
        - status == Online
        - kind == PodList
        - namespace In default,openebs
        - labels == app=jiva,org=openebs
status:
```
