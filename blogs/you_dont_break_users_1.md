### No.1 Rule in Software Development
[You Dont Break Users](https://lkml.org/lkml/2018/8/3/621)


### Here is one attempt
The snippet tries to get an object instance called runtask. The runtask was earlier available via an API call against a 
resource called ConfigMap. Later this same runtask object can be fetched via an API call against a new resource called RunTask
custom resource. How does this relate to No. 1 Rule which is _**You don't break users**_? This entire change works only for 
recent versions of the software. In addition to the recent version the environment should have these objects as custom 
resource. In other words there are multiple things that needs to work to get this object via newer way. This is an condition
ripe for failures at customer's end.

Let us get into the solutioning part. What I have done is to have the older way to fetch the object as well as introduce the 
newer approach for the same. The logic will try the newer approach and in case of any failures will revert to use the 
deprecated way.

_NOTE: I have also included a fair amount reviews on my own code. One can review their own code and can come out with better
suggestions. Just try to give a gap of probably a week's time between writing the code and reviewing the same_

```go
// runTaskGetter abstracts fetching of runtask instance
//
// Notes:
// - The output is a struct which should ideally be an interface
// - The error output is an interface. Good.
// - There are no input arguments. This ensures conformance to many specific structs
// - Above implies specific structs to have as many properties to make logic adhere to interface
type runTaskGetter interface {
	Get() (runtask *v1alpha1.RunTask, err error)
}

// getRunTaskSpec composes common properties required to get a run task
// instance
//
// Notes:
// - Naming is apt:
//  - Single responsibility
//  - Spec gives a touch of struct name
//  - k8sClient is a struct. Can it be an interface?
type getRunTaskSpec struct {
	taskName  string
	k8sClient *m_k8s_client.K8sClient
}

// runTaskGetterFn abstracts fetching of runtask instance based on the provided
// runtask specifications
//
// Notes:
// - Naming is apt
// - This models the struct to be used in a higher order function
// - The signature is similar to earlier interface
type runTaskGetterFn func(getSpec getRunTaskSpec) (runtask *v1alpha1.RunTask, err error)

// getRunTaskFromConfigMap fetches runtask instance from a config map instance
//
// Notes:
// - Naming is good
// - This is the specific implementation of above high order function
// - The input arguments are all composed inside one i.e. spec struct
// - However, unit testing will be problem as structs as composed in struct
func getRunTaskFromConfigMap(g getRunTaskSpec) (runtask *v1alpha1.RunTask, err error) {
	if len(strings.TrimSpace(g.taskName)) == 0 {
		err = fmt.Errorf("missing run task name: failed to get runtask from config map")
		return
	}

	if g.k8sClient == nil {
		err = fmt.Errorf("nil kubernetes client found: failed to get runtask '%s' from config map", g.taskName)
		return
	}

	cm, err := g.k8sClient.GetConfigMap(g.taskName, mach_apis_meta_v1.GetOptions{})
	if err != nil {
		err = errors.Wrap(err, fmt.Sprintf("failed to get run task '%s' from config map", g.taskName))
		return
	}

	runtask = &v1alpha1.RunTask{}
	runtask.Name = g.taskName
	runtask.Spec.Meta = cm.Data["meta"]
	runtask.Spec.Task = cm.Data["task"]
	runtask.Spec.PostRun = cm.Data["post"]

	return
}

// getRunTaskFromCustomResource fetches runtask instance from its custom
// resource instance
// 
// Notes:
// - Naming is good
// - Input arguments are all composed in spec struct
// - Issues during Unit Testing as structs are composed into this spec struct
func getRunTaskFromCustomResource(g getRunTaskSpec) (runtask *v1alpha1.RunTask, err error) {
	if len(strings.TrimSpace(g.taskName)) == 0 {
		err = fmt.Errorf("missing run task name: failed to get runtask from custom resource")
		return
	}

	if g.k8sClient == nil {
		err = fmt.Errorf("nil kubernetes client: failed to get runtask '%s' from custom resource", g.taskName)
		return
	}

	runtask, err = g.k8sClient.GetOEV1alpha1RunTask(g.taskName, mach_apis_meta_v1.GetOptions{})
	if err != nil {
		err = errors.Wrap(err, fmt.Sprintf("failed to get runtask '%s' from custom resource", g.taskName))
	}

	return
}

// getRunTask enables fetching a run task instance based on various run task
// getter strategies
//
// NOTE:
//  This is an implementation of runTaskGetter
//
// Notes:
// - Composes spec and high order functions 
// - These higher order functions make use of this spec
type getRunTask struct {
	// getRunTaskSpec is the specifications required to fetch runtask instance
	getRunTaskSpec
	// currentStrategy is the latest strategy to fetch runtask instance
	currentStrategy runTaskGetterFn
	// oldStrategies are the older strategies to fetch runtask instance
	oldStrategies []runTaskGetterFn
}

// Get returns an instance of runtask
//
// Notes:
// - Clear and concise logic implementation
// - Makes proper use of higher order function
// - Unit testing is possible due to these higher order function
func (g *getRunTask) Get() (runtask *v1alpha1.RunTask, err error) {
	var allStrategies []runTaskGetterFn
	if g.currentStrategy != nil {
		allStrategies = append(allStrategies, g.currentStrategy)
	}

	if len(g.oldStrategies) != 0 {
		allStrategies = append(allStrategies, g.oldStrategies...)
	}

	if len(allStrategies) == 0 {
		err = fmt.Errorf("no strategies to get runtask: failed to get runtask '%s'", g.taskName)
		return
	}

	for _, s := range allStrategies {
		runtask, err = s(g.getRunTaskSpec)
		if err == nil {
			return
		}

		err = errors.Wrap(err, fmt.Sprintf("failed to get runtask '%s'", g.taskName))
		glog.Warningf("%s", err)
	}

	// at this point, we have a real error we can not recover from
	err = fmt.Errorf("exhausted all strategies to get runtask: failed to get runtask '%s'", g.taskName)
	return
}

// Notes:
// - Naming might be improved
// - Input argument i.e. spec is delegated to caller which might be tough for later
// - It will be ideal to have only the taskName and probably namespace
// - If and only if number of arguments are more than 2 then spec can be the input argument
func defaultRunTaskGetter(getSpec getRunTaskSpec) *getRunTask {
	// current strategy
	current := getRunTaskFromCustomResource
	// older strategies
	old := []runTaskGetterFn{getRunTaskFromConfigMap}

	return &getRunTask{
		getRunTaskSpec:  getSpec,
		currentStrategy: current,
		oldStrategies:   old,
	}
}
```
