#### Updates
- Version: 10
- LastUpdatedOn: 21-Feb-2019

#### Low Level Implementation
##### [cstor only] - find eligible node(s) to place the replica	
  - `pkg/cstorpool/v1alpha1/cstorpool.go`
```go
import (
  apis "github.com/openebs/maya/pkg/apis/openebs.io/v1alpha1"
)

type csp struct {
  // actual cstor pool object
  object *apis.CstorPool
}

// predicate defines an abstraction 
// to determine conditional checks
// against the provided csp instance
type predicate func(*csp) bool

type predicateList []predicate

// all returns true if all the predicates
// succeed against the provided csp
// instance
func (l predicateList) all(c *csp) bool {
  for _, pred := range l {
    if !pred(c) {
      return false
    }
  }
  return true
}

// IsNotUID returns true if provided csp 
// instance's UID does not match with any
// of the provided UIDs
func IsNotUID(uids ...string) predicate {
  return func(c *csp) bool {
    for _, uid := range uids {
      if uid == c.object.UID {
        return false
      }
    }
    return true
  }
}

type cspList struct {
  // list of cstor pools
  items   []*csp
}

// FilterUIDs will filter the csp instances 
// if all the predicates succeed against that
// csp. The filtered csp instances' UIDs will
// be returned
func (l *cspList) FilterUIDs(p ...predicate) []string {
  var (
    filtered []string
    plist    predicateList
  )
  plist = append(plist, p...)
  for _, csp := range l {
    if plist.all(csp) {
      filtered = append(filtered, csp.object.UID)
    }
  }
  return filtered
}

// listBuilder enables building a
// list of csp instances
type listBuilder struct {
  list *cspList
}

// ListBuilder returns a new instance of
// listBuilder object
func ListBuilder() *listBuilder {
  return &listBuilder{list: &cspList{}}
}

// WithUIDs builds a list of cstor pools
// based on the provided pool UIDs
func (b *listBuilder) WithUIDs(poolUIDs ...string) *listBuilder {
  for _, uid := range poolUIDs {
    item := &csp{object: &apis.CstorPool{UID: uid}}
    b.list.items = append(b.list.items, item)
  }
  return b
}

// List returns the list of csp
// instances that were built by 
// this builder
func (b *listBuilder) List() *cspList {
  return b.list
}
```

  - `pkg/cstorvolumereplica/v1alpha1/kubernetes.go`
```go
import (
  metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

  kclient "github.com/openebs/maya/pkg/client/k8s/v1alpha1"  
  apis "github.com/openebs/maya/pkg/apis/openebs.io/v1alpha1"
  clientset "github.com/openebs/maya/pkg/client/generated/clientset/internalclientset"
)

// getClientsetFn is a typed function that
// abstracts fetching of internal clientset
type getClientsetFn func() (clientset *clientset.Clientset, err error)

// listFn is a typed function that abstracts
// listing of cstor volume replica instances
type listFn func(cli *clientset.Clientset, namespace string, opts metav1.ListOptions) (*apis.CStorVolumeReplicaList, error)

// kubeclient enables kubernetes API operations
// on cstor volume replica instance
type kubeclient struct {
  // clientset refers to cstor volume replica's
  // clientset that will be responsible to
  // make kubernetes API calls
  clientset     clientset.Interface

  // functions useful during mocking
  getClientset  getClientsetFn
  list          listFn
}

// kubeclientBuildOption defines the abstraction
// to build a kubeclient instance
type kubeclientBuildOption func(*kubeclient)

// withDefaults sets the default options
// of kubeclient instance
func (k *kubeclient) withDefaults() {
	if k.getClientset == nil {
    k.getClientset = func() (clientset *clientset.Clientset, err error) {
      config, err := kclient.Config().Get()
      if err != nil {
        return nil, err
      }
      return clientset.NewForConfig(config)
    }
  }
  if k.list == nil {
    k.list = func(cli *clientset.Clientset, namespace string, opts metav1.ListOptions) (*apis.CStorVolumeReplicaList, error) {
      return cli.OpenebsV1alpha1().CStorVolumeReplicas(namespace).List(opts)
    }
  }
}

// WithKubeClient sets the kubernetes client against
// the kubeclient instance
func WithKubeClient(c clientset.Interface) kubeclientBuildOption {
	return func(k *kubeclient) {
		k.clientset = c
	}
}

// KubeClient returns a new instance of kubeclient meant for
// cstor volume replica operations
func KubeClient(opts ...kubeclientBuildOption) *kubeclient {
	k := &kubeclient{}
	for _, o := range opts {
		o(k)
	}
	k.withDefaults()
	return k
}

// getClientOrCached returns either a new instance
// of kubernetes client or its cached copy
func (k *kubeclient) getClientOrCached() (client.Client, error) {
	if k.clientset != nil {
		return k.clientset, nil
	}
	c, err := k.getClientset()
	if err != nil {
		return nil, err
	}
	k.clientset = c
	return k.clientset, nil
}

// List returns a list of cstor volume replica
// instances present in kubernetes cluster
func (k *kubeclient) List(namespace string, opts metav1.ListOptions) (*apis.CStorVolumeReplicaList, error) {
	cli, err := k.getClientOrCached()
	if err != nil {
		return nil, err
	}
  return k.list(cli, namespace, opts)
}
```

  - `pkg/cstorvolumereplica/v1alpha1/cstorvolumereplica.go`
```go
const (
  cstorPoolUIDLabelKey string = "cstorpool.openebs.io/uid"
)

type cvr struct {
  // actual cstor volume replica
  // object
  object *apis.CstorVolumeReplica
}

type cvrList struct {
  // list of cstor volume replicas
  items []cvr
}

// GetPoolUIDs returns a list of cstor pool
// UIDs corresponding to cstor volume replica
// instances
func (l *cvrList) GetPoolUIDs() []string {
  var uids []string
  for _, cvr := range l.items {
    uid := cvr.Labels[cstorPoolUIDLabelKey]
    uids = append(uids, uid)
  }
  return uids
}

// listBuilder enables building
// an instance of cvrList
type listBuilder struct {
  list *cvrList
}

// ListBuilder returns a new instance
// of listBuilder
func ListBuilder() *listBuilder {
  return &listBuilder{list: &cvrList{}}
}

// WithListObject builds the list of cvr
// instances based on the provided
// cvr api instances
func (b *listBuilder) WithListObject(list *apis.CstorVolumeReplicaList) *listBuilder {
  if list == nil {
    return b
  }
  for _, cvr := range list.Items {
    b.list.items = append(b.list.items, cvr{object: &cvr})
  }
  return b
}

// List returns the list of cvr
// instances that was built by this
// builder
func (b *listBuilder) List() *cvrList{
  return b.list
}
```

  - `pkg/algorithm/cstorpoolselect/v1alpha1/doc.go`
```go
// This package caters to cstorpool's selection
// related operations
```

  - `pkg/algorithm/cstorpoolselect/v1alpha1/select.go`
```go
import (
  csp "github.com/openebs/maya/pkg/cstorpool/v1alpha1"
  cvr "github.com/openebs/maya/pkg/cstorvolumereplica/v1alpha1"
)

// policyName is a type that caters to 
// naming of various pool selection 
// policies
type policyName string

const (
  // antiAffinityLabelPolicyName is the name of the 
  // policy that applies anti-affinity rule based on
  // label
  antiAffinityLabelPolicyName policyName = "anti-affinity-label"
  
  // preferAntiAffinityLabelPolicyName is the name of 
  // the policy that does a best effort while appling
  // anti-affinity rule based on label
  preferAntiAffinityLabelPolicyName policyName = "prefer-anti-affinity-label"
)

// policy exposes the contracts that need
// to be satisfied by any pool selection
// implementation
type policy interface {
  name() policyName
  filter([]string) ([]string, error)
}

// antiAffinityLabel is a pool selection
// policy implementation
type antiAffinityLabel struct {
  labelSelector string
}

// name returns the name of this
// policy
func (p antiAffinityLabel) name() policyName {
  return antiAffinityLabelPolicyName
}

// filter excludes the pool(s) if they are
// already associated with the label 
// selector. In other words, it applies anti
// affinity rule against the provided list of
// pools.
func (l antiAffinityLabel) filter(poolUIDs []string) ([]string, error) {
  if l.labelSelector == "" {
    return poolUIDs, nil
  }
  // pools that are already associated with 
  // this label should be excluded
  //
  // NOTE: we try without giving any namespace
  // so that it lists from all available 
  // namespaces
  cvrs, err := cvr.KubeClient().List("", l.labelSelector)
  if err != nil {
    return nil, err
  }
  exclude := cvr.ListBuilder().WithListObject(cvrs).List().GetPoolUIDs()
  plist := csp.ListBuilder().WithUIDs(poolUIDs).List()
  return plist.FilterUIDs(csp.IsNotUID(exclude...)), nil
}

// preferAntiAffinityLabel is a pool
// selection policy implementation
type preferAntiAffinityLabel struct {
  antiAffinityLabel
}

// name returns the name of this policy
func (p preferAntiAffinityLabel) name() policyName {
  return preferAntiAffinityLabelPolicyName
}

// filter piggybacks on antiAffinityLabel policy
// with the difference being; this logic returns all
// the provided pools if there are no pools that 
// satisfy antiAffinity rule
func (p preferAntiAffinityLabel) filter(poolUIDs []string) ([]string, error) {
  plist, err := p.antiAffinityLabel.filter(poolUIDs)
  if err != nil {
    return nil, err
  }
  if len(plist) > 0 {
    return plist, nil
  }
  return poolUIDs, nil
}

// selection enables selecting required pools
// based on the registered policies
//
// NOTE:
//  There can be cases where multiple policies
// can be set to determine the required pools
//
// NOTE:
//  This code will evolve as we try implementing
// different set of policies
type selection struct {
  // list of original pools aginst whom 
  // selection will be made
  poolUIDs             []string

  // selection is based on these policies
  policies             []policy
}

// buildOption is a typed function that
// abstracts configuring a selection instance
type buildOption func(*selection)

// newSelection returns a new instance of
// selection
func newSelection(poolUIDs []string, opts ...buildOption) *selection {
  s := &selection{poolUIDs: poolUIDs}
  for _, o := opts {
    o(s)
  }
  return s
}

// isPolicy determines if the provided policy 
// needs to be considered during selection
func (s *selection) isPolicy(p policyName) bool {
  for _, pol := range s.policies {
    if pol.name() == p {
      return true
    }
  }
  return false
}

// isPreferAntiAffinityLabel determines if
// prefer anti affinity label needs to be
// considered during selection
func (s *selection) isPreferAntiAffinityLabel() bool {
  return s.isPolicy(preferAntiAffinityLabelPolicyName)
}

// isAntiAffinityLabel determines if anti affinity
// label needs to be considered during
// selection
func (s *selection) isAntiAffinityLabel() bool {
  return s.isPolicy(antiAffinityLabelPolicyName)
}

// PreferAntiAffinityLabel adds anti affinity label
// as a preferred policy to be used during pool 
// selection
func PreferAntiAffinityLabel(lbl string) buildOption {
  return func(s *selection) {
    p := preferAntiAffinityLabel{labelSelector: lbl}
    s.policies = append(s.policies, p)
  }
}

// AntiAffinityLabel adds anti affinity label
// as a policy to be used during pool selection
func AntiAffinityLabel(lbl string) buildOption {
  return func(s *selection) {
    a := antiAffinityLabel{labelSelector: lbl}
    s.policies = append(s.policies, a)
  }
}

// validate runs some validations/checks
// against this selection instance
func (s *selection) validate() error {
  if s.isAntiAffinityLabel() && s.isPreferAntiAffinityLabel() {
    return errors.New("invalid selection: both antiAffinityLabel and preferAntiAffinityLabel policies can not be together")
  }
  return nil
}

// filter returns the final list of pools that
// gets selected, after passing the original list
// of pools through the registered selection policies
func (s *selection) filter() ([]string, error) {
  var (
    filtered []string
    err error
  )
  if len(s.policies) == 0 {
    return s.poolUIDs, nil
  }
  // make a copy of original pool UIDs
  filtered = append(filtered, s.poolUIDs...)
  for _, policy := range s.policies {
    filtered, err = policy.filter(filtered)
    if err != nil {
      return nil, err
    }
  }
  return filtered, nil
}

// Filter will return filtered pool UIDs
// from the provided list based on pool 
// selection options
func Filter(origPoolUIDs []string, opts ...buildOption) ([]string, error) {
  if len(opts) == 0 {
    return origPoolUIDs, nil
  }
  s := newSelection(origPoolUIDs, opts...)
  err := s.validate()
  if err != nil {
    return nil, err
  }
  return s.filter()
}

// expose below go functions as go template functions
// Filter                      as cspFilter
// AntiAffinityLabel           as cspAntiAffinity
// PreferAntiAffinityLabel     as cspPreferAntiAffinity
```
  - template will look something similar to below:
```yaml
- {{- $poolUids := keys .ListItems.cvolPoolList.pools }}
- {{- $lblSelector := ifNotNil $appUniqVal "openebs.io/replica-anti-affinity: $appUniqVal" }}
- {{- $preferAntiAffinity := cspPreferAntiAffinity $lblSelector }}
- {{- $poolUids = cspFilter $poolUids $preferAntiAffinity | randomize }}
```
  - finally when a CVR is created, the label needs to be applied against the CVR
    - i.e. `openebs.io/replica-anti-affinity: $appUniqVal`
    - where $appUniqVal is the label extracted from the volume's PVC
    - if $appUniqVal is empty then
      - no need to set the label against the CVR
