### Installer Design
At my current day job, I have been tasked to write an installer program. It came in for a shock to me that installer is 
considered so important at this stage of our product which is still at 0.6 release version. Thinking in terms of ROE, I 
would have assumed this to be a low returns on investment.

### Disconfirming Evidence
Now lets try to think otherwise before getting biased with my first conclusion. There have been a few niggling problems 
that I have been noticing with our existing installation practices. We depend on a tool called _**kubectl**_ and some 
yaml files to deploy our components. We have also been making use of another tool called _**helm**_ to do the same but 
by delegating the defaults and other specifics to the end user. In other words, helm enables our users to specify 
exactly how our product should get installed. However, there is no integration between the kubectl way & the helm way. 
We have concluded ourselves that kubectl way is a simple way that a team either internal to our company or external 
i.e. our users, can use without bothering about the product internals.

This was all fine till our teams wanted to install the product in a way that suits their purpose. With above logic, one
would advice this team to make use of helm. This is where things start to get murky. We need to update the helm charts & are
dependant on core maintainers for such activities. Some of the component installation require yet to be released helm binary
bringing our entire company to a standstill till above process is sorted out.

The reader of this article would have probably concluded by now that we need a installer of our own. Some of the assumptions
that might be doing rounds in the reader's mind are:
- Decouple the dependency from helm
- Using similar tools might lead to above problems
- A home grown solution might fit our requirements

### There is lot more
It is not just about decoupling from above tools to get the install done. There is lot more to installation than just 
installing the solution on any given environment. On one hand, installation need to bother about cloud versus on-premise
environments. At the same time, it must abstract away the product internals from the end user. Expecting these, which also
sounds to be contradictory statements, to be available in community tools does not seem convincing. In fact this is one of
the reasons, we find a lot of home grown as well as community developed  operators whose main job is to provide a lifecycle
to install and ensure a smooth day two related operations. Scale up, scale  down, upgrade, remove, etc are other factors 
that are very much a part of _**install lifecycle**_. Install is no longer a one time activity but a reconcile loop that 
runs throughout the lifecycle of the product.

On a concluding note, this move from being an activity to a process should not even be a thing of concern to the end users.

### High Level Design Decisions
I will try to write down my thoughts with respect to high level design. This enables me to remain un-biased with first
thoughts that hit my brain. I try to write down several thoughts that come in later and either disconfirm or confirm with my earlier thoughts. 

_NOTE: There will be a lot of initial thoughts, which I have written down but may not (will not) actually implement based on
the thoughts that come in later due to better evidence, interactions, findings, real needs versus speculations and so on._

#### First Time Conclusions
- Install should solve the problem of end users trying to deploy massive amounts of yamls
- Make use of ENV to install or un-install
- Install the current release only
- Unstall the current release only
- Code the entire yaml content into specific strings
- Make use of http client to install these yaml strings
- Make use of yaml as a install config to install other yaml(s)
- Install based on the install yaml
- Un-install based on the un-install yaml

#### Second Time Conclusions
- Make use of Integration Tests to test install & un-install
- Make use of go templates to implement install & un-install

#### Third Time Conclusions
- Install should solve problems faced by different teams i.e. internal & external
  - Install can solve above by exposing helm like templating
- Install should be able to accept values from helm
- Install should be able to accept overlay yamls from kustomize
- Make use of interface oriented programming

#### Fourth Time Conclusions
- Install should not bother dealing with helm, kustomize or something else
- Install should be able to work with any tool without being a provider of that tool
- Make use of ConfigMap as a config to install and un-install
- There will be a single ConfigMap that manages the install and un-install
- ConfigMap will be referred via an ENV variable
- ConfigMap will be fetched via kubernetes client
- ConfigMap can be populated with default values
- ConfigMap should be easy to deploy by helm, kustomize or team's desired tool
- Avoid interfaces by replacing them with high order functions
- Make use of pkg/errors while throwing errors
- Accumulate error(s) during install and un-install
- Do not panic due to install or un-install errors
- An install resource can be a native kubernetes resource or a custom resource
- Install will install all the runtasks as-is
- Install can install multiple versions of a resource
- Install can un-install multiple versions of a resource
- Install can install multiple flavours of a resource
- Install & Un-install will be coupled to resource
- User is free to not use this install & un-install feature
- Make use of versioning in namespaces to manage the technical debt without impact
- Install will default to re-install i.e. remove if already installed & re-install once again

#### Fifth Time Conclusions
- Install will default to install if not installed
- Install will `update` or `patch` or `do nothing` if already installed (RnD)
- Install will not un-install during a panic of installer
- Install will not un-install during a shutdown of installer

#### Sixth Time Conclusions
- Install will install the resource with a custom label e.g. `installer.openebs.io/release: 0.7.0`
- Install will re-install i.e. remove & add if a resource already exists
- Above will continue till kubectl like _**apply**_ feature is exposed via client-go
- Install will manage very specific resources like CASTemplate & RunTasks in its first release

#### Seventh Time Conclusions
- Install & UnInstall will deal with resources defined in installer
- It is not a generic installer & un-installer

### Final Install Config:
```yaml
spec:
  install:
    # release version to install
    - version: 0.7.0
      # will make use of this namespace for all the appropriate resources
      withNamespace:
      # will make use of this service account for all the appropriate resources
      withServiceAccount:
      # will add these label(s) to all the resources
      withLabels:
      # will add these annotation(s) to all the resources
      withAnnotations:
  uninstall:
    # remove all resources from test namespace with 0.7.0 as the version
    - version: 0.7.0
      namespace: test
    # remove all the resources from openebs namespace with 0.6.5 as the version
    - version: 0.6.5
      namespace: openebs
    # remove all the resources from all namespaces with 0.6.6 as the version
    - version: 0.6.6
```

### Low Level Design Decisions
I follow the same approach with low level design thoughts as well. However, since this is pretty much a coding activity, it
depends on how I start implementing the code and obviously depends a lot on my previous coding skills. Having said this, 
code is also a derivative of brain's biasedness towards a particular style of implementation. The important point to 
consider here is _"It is best to embrace simplicity than be stupid to attempt perfection from the word go"_.

#### First Time Conclusions
- Packaging structure consisting of versioned folders:
  - pkg/install/v1alpha1
  - pkg/k8s/v1alpha1
  - pkg/client/k8s/v1alpha1
  - install/v1alpha1-0.7.0/cstor-pool/
  - install/v1alpha1-0.7.0/cstor-volume/
  - install/v1alpha1-0.7.0/jiva-volume/

#### Second Time Conclusions
- cmd/maya-apiserver/app/command/start.go will trigger install

#### Third Time Conclusions
- Build high order function based strategies only if they are required
  - One can always design a function based strategy & compose it inside the service based struct **later**
  - Above refactoring will have little impact
- High order functions are fine as long as they are independent
  - In other words, if functions do not depend on any other packages, then being a function suffices
- If high order functions need to use other packages, or other functions, or strategies, or depend on multiple arguments:
  - Then convert this high order function into a single method interface
  - Design a struct that implements above interface
  - Ensure this struct is composed of other functions, strategies, arguments, interfaces
- Start with higher order function & refactor it to an interface if required
- When the service structure that is composed of various interfaces, or functions needs to be set dynamically:
  - Then create a builder of above structure
  - In other words a fluent API might serve the requirements at the caller side
  - For Example: 
  ```go
  ServiceBuilder()
    .Config(xyzConfigFn)
    .Executor(newExecutor())
    .Build()
  ```
  - NOTE - A typical instantiation function like `NewService()` may not work out in these cases
