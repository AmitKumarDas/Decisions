### Motivation
This is a POC implementation of ideas mentioned at following link:
https://github.com/AmitKumarDas/Decisions/blob/master/blogs/declarative_with_nouns_%26_verbs.md

### UseCases
#### UseCase 1
- I want to get a list of namespaces of running pods from a given pod list
- I want to save this list as `myNS`
```yaml
kind: RunTask
spec:
  meta:
  task:
  post:
  - run: GetNamespaceList
    for:
      - --kind=PodList
      - --object-path=.taskresult.id101.pods
    withFilter:
      - --check=IsRunning
    as: myNS
```

#### UseCase 2
- I want to get a list of name, namespace, uuid of running pods from a given pod list
- I want to also filter this for pods with label app=jiva
- I want to save this list as myPodInfo
```yaml
kind: RunTask
spec:
  meta:
  task:
  post:
  - run: GetTupleList
    for:
      - --kind=PodList
      - --object-path=.taskresult.id101.pods
    withFilter:
      - --check=IsRunning
      - --check=IsLabel --label=app:jiva
    withOutput:
      - --name
      - --namespace
      - --uuid
    as: myPodInfo
```

#### UseCase 3
- I want to get a map of disk name & disk status from a given list of CSPs
- I want to filter this for CSPs with following labels:
  - openebs.io/version=0.9.0
  - app=jiva
- I want to additionally filter CSPs with Offline state
- I want to save this list as myDiskStatusInfo
```yaml
kind: RunTask
spec:
  meta:
  task:
  post:
  - run: GetDiskStatusMap
    for:
      - --kind=CStorPoolList
      - --json-path=.taskresult.id101.csps
    withFilter:
      - --check=IsLabel --label=openebs.io/version:0.9.0
      - --check=IsLabel --label=app:jiva
      - --check=IsOffline
    as: myDiskStatusInfo
```

#### UseCase 4
- I want following details from a PVC
  - hostname from annotations
  - replica anti affinity from labels
  - preferred replica anti affinity from labels
  - target affinity
  - sts target affinity
- I want to save this as myPVCInfo
```yaml
kind: RunTask
spec:
  meta:
  task:
  post:
  - run: GetAnnotations
    for:
      - --kind=PersistentVolumeClaim
      - --json-path=.taskresult.id101.pvc
    withOutput:
      - --annotation --key=volume.kubernetes.io/selected-node
    as: myPVCInfo.anns
  - run: GetLabels
    for:
      - --kind=PersistentVolumeClaim
      - --json-path=.taskresult.id101.pvc
    withOutput:
      - --label --key=openebs.io/replica-anti-affinity
      - --label --key=openebs.io/preferred-replica-anti-affinity
      - --label --key=openebs.io/target-affinity
      - --label --key=openebs.io/sts-target-affinity
    as: myPVCInfo.labels
```

#### UseCase 5
- I want to create a pod with following details:
  - name of the pod is `mypod`
  - namespace of the pod is `default`
  - pod should have label `app=jiva`
  - pod image should be `openebs.io/m-apiserver:1.0.0`
- Save this instance of pod as myPod

```yaml
kind: RunTask
spec:
  meta:
  task:
  post:
  - run: Build
    for: 
      - --kind=Pod
    withOptions: 
      - --label=app:jiva
      - --name=mypod
      - --namespace=default
      - --image=openebs.io/m-apiserver:1.0.0
    as: myPod
```

#### UseCase 6
- I should be able to specify a pod yaml
- I woul like to save this as myPodie

```yaml
kind: RunTask
spec:
  meta:
  task:
  post:
  - it: should save pod yaml as pod instance
    run: Build
    for:
      - --kind=Pod
      - |
        --yaml=
        kind: Pod
        apiVersion: v1
        metadata:
          name: myPoddy
          namespace: default
        spec:
          image: openebs.io/m-apiserver:ci
    as: myPodie
```

### Standards & Best Practices
- Some of the values used in post field are keywords
- A keyword is better expressed as CapitalCased
- However, logic should be able to handle keywords that are not specified as CapitalCased


#### Sample RunTask with new `Post` schema
```yaml
apiVersion: openebs.io/v1alpha1
kind: RunTask
metadata:
  name: list-ctrl-deployment
  namespace: default
spec:
  meta: |
    id: listctrldeployment
    apiVersion: extensions/v1beta1
    kind: Deployment
    action: list
    runNamespace: default
  post: |-
    operations:
    - run: GetTupleList # GetTupleList is a keyword - hence CapitalCased
      for:
      - --kind=DeploymentList # DeploymentList is a keyword - hence CapitalCased
      - --object-path=.task.runtimeObj # object-path is right - objectPath is wrong
      withFilter:
      - --check=IsLabel --label=app:nginx # IsLabel is a keyword - hence CapitalCased
      withOutput:
      - --name # name is right; while Name is wrong; adhere to flag arg naming standards
      - --namespace
      as: taskResult.tupleList
```


### Low Level Impl
```go
// pkg/apis/openebs.io/runtask/v1beta2/post.go

// Post represents the desired specifications for the
// post field of runtask. It specifies the action or the
// commands that need to be executed post executing a
// runtask i.e. save the status or name of the current
// resource so that the next runtask can make use of it.
type Post struct {
	Operations []Operation `json:"operations"`
}
```

```go
// Operation defines a particular operation to be
// executed against a particular resource
type Operation struct {
	// Run declares the operation name
	Run        string   `json:"run"`
	
	// For declares the resource against
	// whom this operation will get
	// executed
	For        []ForOption `json:"for"`
	
	// 
	WithFilter []string `json:"withFilter"`
	WithOutput []string `json:"withOutput"`
	As         string   `json:"as"`
}

// TopLevelProperty represents the top level property that
// is a starting point to represent a hierarchical chain of
// properties.
//
// e.g.
// Config.prop1.subprop1 = val1
// Config.prop1.subprop2 = val2
// In above example Config is a top level object
//
// NOTE:
//  The value of any hierarchical chain of properties
// can be parsed via dot notation
type TopLevelProperty string

const (
	// CurrentRuntimeObjectTLP is a top level property supported by CAS template engine
	// The runtime object of the current task's execution is stored in this top
	// level property.
	CurrentRuntimeObjectTLP TopLevelProperty = "RuntimeObject"
)
```

```go
// DeployList is the list of deployments
type DeployList struct {
	items []*Deploy
}

// ListBuilder enables building an instance of
// deploymentList
type ListBuilder struct {
	list *DeployList
	//output []map[string]interface{}
	output  outputMap
	filters filterList
	errors  []error
}

// filterList is a list of filters
type filterList []Filter

// outputMap is a map of desired output function
// and their reference key
type outputMap map[*Output]string

// Filter abstracts filtering logic w.r.t the
// deployment list instance
type Filter func(*Deploy) bool

// Output represents the desired output to
// be given against a deployment instance
type Output func(*Deploy) interface{}

// NewListBuilder returns an instance of
// list builder
func NewListBuilder() *ListBuilder {
	return &ListBuilder{
		list:   &DeployList{},
		output: outputMap{},
	}
}

// ListBuilderForRuntimeObject returns a list builder instance
// for deployment list
func ListBuilderForRuntimeObject(obj runtime.Object) *ListBuilder {
	var (
		dl *extn_v1_beta1.DeploymentList
		ok bool
	)
	lb := NewListBuilder()
	if obj == nil {
		lb.errors = append(lb.errors, errors.New("failed to build instance: nil runtime.Object given"))
		return lb
	}
	// Convert the runtime.Object to its desired type i.e.
	// DeploymentList here
	if dl, ok = obj.(*extn_v1_beta1.DeploymentList); !ok {
		lb.errors = append(lb.errors, errors.New(
			"failed to build instance: unable to typecast given object to deployment list"))
		return lb
	}
	// Iterate over deployment list objects and
	// insert it into the ListBuilder instance
	for _, d := range dl.Items {
		d := d
		deploy := &Deploy{}
		deploy.object = &d
		lb.list.items = append(lb.list.items, deploy)
	}
	return lb
}

// List returns a list of deployments after doing
// all the filtering and validations
func (lb *ListBuilder) List() (*DeployList, error) {
	if len(lb.errors) != 0 {
		return nil, errors.Errorf("%v", lb.errors)
	}
	if lb.filters == nil || len(lb.filters) == 0 {
		return lb.list, nil
	}
	filtered := &DeployList{}
	for _, d := range lb.list.items {
		d := d
		if lb.filters.all(d) {
			filtered.items = append(filtered.items, d)
		}
	}
	return filtered, nil
}

// all returns true if all the predicates
// succeed against the provided deployment
// instance
func (f filterList) all(d *Deploy) bool {
	for _, filter := range f {
		filter := filter
		if !filter(d) {
			return false
		}
	}
	return true
}

// AddFilter adds the filter to be applied against the
// deployment instance
func (lb *ListBuilder) AddFilter(f Filter) *ListBuilder {
	lb.filters = append(lb.filters, f)
	return lb
}

// AddFilters adds the provided filters to be applied against
// the deployment instance
func (lb *ListBuilder) AddFilters(filters ...Filter) *ListBuilder {
	for _, filter := range filters {
		filter := filter
		lb.AddFilter(filter)
	}
	return lb
}

// HasLabel returns HasLabel filter for
// the given label
func HasLabel(label string) Filter {
	return func(d *Deploy) bool {
		return d.HasLabel(label)
	}
}

// HasLabel checks if the given label is
// present or not for a particular deployment
// object
func (d *Deploy) HasLabel(label string) bool {
	labels := d.object.GetLabels()
	if _, exist := labels[label]; exist {
		return true
	}
	return false
}

// HasLabels returns IsLabel filter for
// the given label
func HasLabels(labels ...string) Filter {
	return func(d *Deploy) bool {
		return d.HasLabels(labels...)
	}
}

// HasLabels checks if the given labels are
// present or not for a particular deployment
// object
func (d *Deploy) HasLabels(labels ...string) bool {
	const (
		keyIndex   = 0
		valueIndex = 1
	)
	gotLabels := d.object.GetLabels()
	for _, label := range labels {
		label := label
		var labelDetails []string
		// get the label key and value by splitting
		// it based on delimiter '=' or ':'
		if strings.Contains(label, "=") {
			labelDetails = strings.Split(label, "=")
		} else {
			labelDetails = strings.Split(label, ":")
		}
		if lValue, exist := gotLabels[labelDetails[keyIndex]]; exist {
			if lValue != labelDetails[valueIndex] {
				return false
			}
		} else {
			return false
		}
	}
	return true
}

// Name returns an output instance
// for getting name of a deployment
func Name() Output {
	return func(d *Deploy) interface{} {
		return d.Name()
	}
}

// Name returns the name of the given deployment
func (d *Deploy) Name() interface{} {
	return d.object.GetName()
}

// Namespace returns an output instance
// for getting namespace of a deployment
func Namespace() Output {
	return func(d *Deploy) interface{} {
		return d.Namespace()
	}
}

// Namespace returns the namespace of the given deployment
func (d *Deploy) Namespace() interface{} {
	return d.object.GetNamespace()
}

// Labels returns an output instance
// for getting labels of a deployment
func Labels() Output {
	return func(d *Deploy) interface{} {
		return d.Labels()
	}
}

// Labels returns the labels of the given deployment
func (d *Deploy) Labels() interface{} {
	return d.object.GetLabels()
}

// WithOutput returns a listBuilder instance having
// the desired output key added
func (lb *ListBuilder) WithOutput(o Output, referenceKey string) *ListBuilder {
	if o == nil || referenceKey == "" {
		lb.errors = append(lb.errors, errors.Errorf(
			"nil reference key given for output %v", o))
		return lb
	}
	lb.output[&o] = referenceKey
	return lb
}

// TupleList enables building a tuple list instance
// against a deployment list
type TupleList []map[string]interface{}

// GetTupleList returns a tuple list based on the desired
// outputs provided
func (lb *ListBuilder) GetTupleList() (TupleList, error) {
	tList := TupleList{}
	dList, err := lb.List()
	if err != nil {
		return nil, errors.Errorf(
			"failed to get tuple list: error: %v", err)
	}
	for _, d := range dList.items {
		d := d
		deployDetails := d.getDesiredDetails(lb.output)
		tList = append(tList, deployDetails)
	}
	return tList, nil
}

func (d *Deploy) getDesiredDetails(desiredDetails outputMap) map[string]interface{} {
	deployDetails := make(map[string]interface{})
	for out, refKey := range desiredDetails {
		out := out
		refKey := refKey
		i := (*out)(d)
		deployDetails[refKey] = i
	}
	return deployDetails
}
```

```go
// pkg/task/post/operations/v1beta2/deploymentlist.go 

const (
	getTupleList = "gettuplelist"
)

// DeploymentList represents the details for fetching
// the desired values from a deployment list object
type DeploymentList struct {
	Run          string
	DataPath     string
	Values       map[string]interface{}
	WithFilterFs *filter
	WithOutputFs *output
}

// Builder enables building an instance of
// DeploymentList
type Builder struct {
	DeploymentList *DeploymentList
	errors         []error
}

type output struct {
	name      bool
	namespace bool
}

type filter struct {
	isLabel isLabels
}

type isLabels struct {
	labels []string
}

func (labelArr *isLabels) String() string {
	return fmt.Sprint(labelArr.labels)
}

func (labelArr *isLabels) Set(value string) error {
	if value == "" {
		return nil
	}
	labelArr.labels = append(labelArr.labels, value)
	return nil
}

func (labelArr *isLabels) Type() string {
	return "isLabels"
}

// NewBuilder returns a new instance of Builder
func NewBuilder() *Builder {
	return &Builder{
		DeploymentList: &DeploymentList{
			Values:       make(map[string]interface{}),
			WithFilterFs: &filter{},
			WithOutputFs: &output{},
		},
	}
}

// BuilderForPostValues returns a builder instance
// after filling the required post values
func BuilderForPostValues(run string,
	withFilter []string, withOutput []string) *Builder {
	b := NewBuilder()
	f := &filter{}
	o := &output{}

	// withFilter is the flagset having all the flags defined for "withFilter" field
	// of post operations
	withFilterFs := flag.NewFlagSet("withFilter", flag.ExitOnError)
	withFilterFs.Var(&f.isLabel, "isLabel", "checks for given labels against a resource")

	// withOutput is the flagset having all the flags defined for "withOutput" field
	// of post operations
	withOutputFs := flag.NewFlagSet("withOutput", flag.ExitOnError)
	withOutputFs.BoolVar(&o.name, "name", false, "set if output should contain resource name")
	withOutputFs.BoolVar(&o.namespace, "namespace", false, "set if output should contain resource namespace")

	if err := withFilterFs.Parse(withFilter); err != nil {
		b.errors = append(b.errors, errors.Errorf("error parsing withFilter flags, error: %v", err))
	}
	if err := withOutputFs.Parse(withOutput); err != nil {
		b.errors = append(b.errors, errors.Errorf("error parsing withOutput flags, error: %v", err))
	}
	b.DeploymentList.WithFilterFs = f
	b.DeploymentList.WithOutputFs = o
	b.DeploymentList.Run = run
	return b
}

// WithDataPath returns a builder instance with
// objectPath/jsonPath of the saved result
func (b *Builder) WithDataPath(path string) *Builder {
	// Trim unnecessary spaces or dots
	path = strings.Trim(strings.TrimSpace(path), ".")
	b.DeploymentList.DataPath = path
	return b
}

// WithTemplateValues returns a builder instance with
// templateValues map
func (b *Builder) WithTemplateValues(values map[string]interface{}) *Builder {
	b.DeploymentList.Values = values
	return b
}

// Build returns the final instance of post
func (b *Builder) Build() (*DeploymentList, error) {
	if len(b.errors) != 0 {
		return nil, errors.Errorf("%v", b.errors)
	}
	return b.DeploymentList, nil
}

// ExecuteOp executes the post operation on a
// deploymentList instance
func (d *DeploymentList) ExecuteOp() (result interface{}, err error) {
	switch d.Run {
	case getTupleList:
		result, err = d.getTupleList()
	default:
		return result, errors.Errorf(
			"unsupported runtask post operation, `%s` for deploymentList", d.Run)
	}
	if err != nil {
		return result, err
	}
	return result, nil
}

func (d *DeploymentList) getTupleList() (tList []map[string]interface{}, err error) {
	var (
		dListObj runtime.Object
		ok       bool
		lb       *deployment_extnv1beta1.ListBuilder
	)
	fields := strings.Split(d.DataPath, ".")
	// TODO: Need to have a function to get runtime.Object
	// directly instead of interface{}
	data := util.GetNestedField(d.Values, fields...)
	if dListObj, ok = data.(runtime.Object); !ok {
		return nil, errors.New("failed to get tuple list: list given is not runtime.Object type")
	}
	// Form the build instance
	lb = deployment_extnv1beta1.
		ListBuilderForRuntimeObject(dListObj)
	// Call the required filters and output build functions
	// based on provided flags
	lb = d.addFuncForFlags(lb)
	// Call the target function now i.e. tupleList for getting the list
	// of tuples
	tList, err = lb.GetTupleList()
  	if err != nil {
		return nil, err
	}
	return tList, nil
}

func (d *DeploymentList) addFuncForFlags(
	lb *deployment_extnv1beta1.ListBuilder) *deployment_extnv1beta1.ListBuilder {
	// Check for the withFilter flags which has been set and
	// accordingly call the corresponding build functions
	//
	// Check if the isLabel flag is set or not
	if len(d.WithFilterFs.isLabel.labels) != 0 {
		lb = lb.AddFilter(deployment_extnv1beta1.HasLabels(d.WithFilterFs.isLabel.labels...))
	}
	// Check for the withOutput flags that has been set and
	// accordingly call the corresponding build functions
	//
	// Check if the name flag is set or not
	if d.WithOutputFs.name {
		lb = lb.WithOutput(deployment_extnv1beta1.Name(), "name")
	}
	//Check if the namespace flag is set
	if d.WithOutputFs.namespace {
		lb = lb.WithOutput(deployment_extnv1beta1.Namespace(), "namespace")
	}
	return lb
}

// TODO: Need to have these functions at a common place
//
// saveAs stores the provided value at specific hierarchy as mentioned in the
// fields inside the values object.
//
// NOTE:
//  This hierarchy along with the provided value is added or updated
// (i.e. overridden) in the values object.
//
// NOTE:
//  fields is represented as a single string with each field separated by dot
// i.e. '.'
//
// Example:
// {{- "Hi" | saveAs "TaskResult.msg" .Values | noop -}}
// {{- .Values.TaskResult.msg -}}
//
// Above will result in printing 'Hi'
// Assumption here is .Values is of type map[string]interface{}
func saveAs(fields string, destination map[string]interface{}, given interface{}) interface{} {
	fieldsArr := strings.Split(fields, ".")
	// save the run task command result in specific way
	r, ok := given.(v1alpha1.RunCommandResult)
	if ok {
		resultpath := append(fieldsArr, "result")
		util.SetNestedField(destination, r.Result(), resultpath...)
		errpath := append(fieldsArr, "error")
		util.SetNestedField(destination, r.Error(), errpath...)
		debugpath := append(fieldsArr, "debug")
		util.SetNestedField(destination, r.Debug(), debugpath...)
		return given
	}
	util.SetNestedField(destination, given, fieldsArr...)
	return given
}
```
```go
// pkg/task/post/v1beta2/exec.go 

const (
	deploymentList         = "deploymentlist"
	asTemplatedBytesOutput = "''"
)

// Post represents the executor struct for runtask's post operations
type Post struct {
	// postTask holds the task's post information
	PostTask       *runtask.Post
	values         map[string]interface{}
	metaID         string
	metaAPIVersion string
	metaKind       string
	metaAction     string
	forFlag        ForFlag
}

// ForFlag specifies the flags supported by for field
// of runtask post field
type ForFlag struct {
	kind       string
	objectPath string
	jsonPath   string
}

// Builder enables building an instance of
// post
type Builder struct {
	Executor *Post
	errors   []error
}

// NewBuilder returns a new instance of Builder
func NewBuilder() *Builder {
	return &Builder{
		Executor: &Post{
			PostTask: &runtask.Post{},
			values:   make(map[string]interface{}),
		},
	}
}

// WithTemplate returns a builder instance with post and
// template values
func (b *Builder) WithTemplate(context, yamlTemplate string,
	values map[string]interface{}) *Builder {
	p := &runtask.Post{}
	if len(yamlTemplate) == 0 {
		// nothing needs to be done
		b.errors = append(b.errors, errors.Errorf("empty post yaml given: %+v", yamlTemplate))
		return b
	}
	// transform the yaml with provided templateValues
	t, err := template.AsTemplatedBytes(context, yamlTemplate, values)
	if err != nil {
		b.errors = append(b.errors, err)
		return b
	}
	// TODO: Need to have a better approach to handle this case.
	//
	// Check if string of templated bytes is single quotes ('')
	// if yes, then do nothing and return (this is required
	// to handle the case where only go-templating is being used in runtask)
	if string(t) == asTemplatedBytesOutput {
		return b
	}
	// unmarshal the yaml bytes into Post struct
	err = yaml.Unmarshal(t, p)
	if err != nil {
		b.errors = append(b.errors, err)
		return b
	}
	b.Executor.PostTask = p
	b.Executor.values = values
	return b
}

// WithMetaID returns a builder instance after
// setting runtask meta id in builder
func (b *Builder) WithMetaID(id string) *Builder {
	b.Executor.metaID = id
	return b
}

// WithMetaAPIVersion returns a builder instance after
// setting runtask meta APIVersion in builder
func (b *Builder) WithMetaAPIVersion(apiVersion string) *Builder {
	b.Executor.metaAPIVersion = apiVersion
	return b
}

// WithMetaAction returns a builder instance after
// setting runtask meta action in builder
func (b *Builder) WithMetaAction(action string) *Builder {
	b.Executor.metaAction = action
	return b
}

// WithMetaKind returns a builder instance after
// setting runtask meta kind in builder
func (b *Builder) WithMetaKind(kind string) *Builder {
	b.Executor.metaKind = kind
	return b
}

// Build returns the final instance of post
func (b *Builder) Build() (*Post, error) {
	if len(b.errors) != 0 {
		return nil, errors.Errorf("%v", b.errors)
	}
	return b.Executor, nil
}

// Execute will execute all the runtask post operations
func (p *Post) Execute() error {
	pTask := p.PostTask
	for _, operation := range pTask.Operations {
		operation := operation
		err := p.executeOp(&operation)
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *Post) executeOp(o *runtask.Operation) (err error) {
	// register and parse the for flags if defined
	err = registerAndParseForFlags(o.For)
	if err != nil {
		return err
	}
	// var result interface{}
	// set value of kind flag of for if its not set
	err = p.setPostKindIfEmpty()
	if err != nil {
		return err
	}
	// set the dataPath to default dataPath i.e.
	// RuntimeObject (top-level property for saving current
	// runtask data)
	dataPath := p.setDataPathIfEmpty()
	run := strings.ToLower(o.Run)
	switch strings.ToLower(p.forFlag.kind) {
	case deploymentList:
		dlist, err := deployList.
			BuilderForPostValues(run, o.WithFilter, o.WithOutput).
			WithDataPath(dataPath).
			WithTemplateValues(p.values).
			Build()
		if err != nil {
			return err
		}
		_, err = dlist.ExecuteOp()
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("unsupported kind for runtask post operation, %s", p.forFlag.kind)
	}
	// TODO: Save this result at path specified by 'as' field of runtask post
	as := o.As
	fmt.Println(as)
	return nil
}

func (p *Post) setPostKindIfEmpty() error {
	// Check if the kind flag has been set, if not
	// then use the kind set in meta of runtask
	// Note: For list action, kind would be kind + action
	// i.e. action = list and kind = deployment then kind flag
	// will be set as deploymentlist
	if p.forFlag.kind == "" {
		if p.metaAction == "" || p.metaKind == "" {
			return errors.Errorf("Empty value for kind or action found, kind: %s, action: %s", p.metaKind, p.metaAction)
		}
		if p.metaAction == "list" {
			p.forFlag.kind = p.metaKind + p.metaAction
		} else {
			p.forFlag.kind = p.metaKind
		}
	}
	return nil
}

func (p *Post) setDataPathIfEmpty() (dataPath string) {
	// Check if the dataPath is set or not, if not
	// set the dataPath as runtime.Object top-level property
	// i.e. RuntimeObject or json result top-level property i.e.
	// JsonResult
	if p.forFlag.jsonPath != "" {
		dataPath = p.forFlag.jsonPath
	} else if p.forFlag.objectPath != "" {
		dataPath = p.forFlag.objectPath
	} else {
		// If none of the above flags are set then
		// use "RuntimeObject" as the top-level property
		// to get data
		dataPath = string(runtask.CurrentRuntimeObjectTLP)
	}
	return dataPath
}

func registerAndParseForFlags(For []string) error {
	ff := ForFlag{}
	// forFs is the flagset having all the flags defined for "For" field
	// of post operations
	forFs := flag.NewFlagSet("for", flag.ExitOnError)
	forFs.StringVar(&ff.kind, "kind", "", "represents the kind of resource")
	forFs.StringVar(&ff.objectPath, "objectPath", "", "represents the runtime object path for fetching the details")
	forFs.StringVar(&ff.jsonPath, "jsonPath", "", "represents the jsonResult path for fetching the details")
	if err := forFs.Parse(For); err != nil {
		return errors.Errorf("failed to parse post forFlags: given forFlags: %s, error: %v", For, err)
	}
	return nil
}
```
