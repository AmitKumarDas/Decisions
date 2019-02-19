#### Low Level Implementation
##### [cstor only] - find eligible node(s) to place the replica
  - `pkg/cstorpool/v1alpha1`
```go
type csp struct {
  // name of cstor pool
  name string
}

type predicate func(*csp) bool

type predicateList []predicate

func (l predicateList) all(c *csp) bool {
  for _, pred := range l {
    if !pred(c) {
      return false
    }
  }
  return true
}

// IsNotName returns false if csp 
// does not belong to any of the provided
// names
func IsNotName(names ...string) predicate {
  return func(c *csp) bool {
    for _, name := range names {
      if name == c.name {
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

func (l *cspList) FilterNames(p ...predicate) []string {
  var (
    filtered []string
    plist    predicateList
  )
  plist = append(plist, p...)
  for _, csp := range l {
    if plist.all(csp) {
      filtered = append(filtered, csp.name)
    }
  }
  return filtered
}

type listBuilder struct {
  list *cspList
}

func ListBuilder() *listBuilder {
  return &listBuilder{list: &cspList{}}
}

func (b *listBuilder) WithNames(pools ...string) *listBuilder {
  for _, pool := range pools {
    item := &csp{name: pool}
    b.list.items = append(b.list.items, item)
  }
  return b
}

func (b *listBuilder) List() *cspList {
  return b.list
}
```

  - `pkg/cstorvolumereplica/v1alpha1`
```go
// TODO
// exclude := cvr.ListBuilder().WithLabelOption(l.labelSelector).List().GetPoolNames()

type cvr struct {
  // name of the cvr
  name string
}

type cvrList struct {
  // option to use to filter
  // during list execution
  labelOption  string

  // list of cstor volume replicas
  items []cvr
}

func (l *cvrList) GetPoolNames() []string {
  var names []string
  for _, cvr := range l.items {
    names = append(names, cvr.name)
  }
  return names
}

type listBuilder struct {
  list *cvrList
}

func ListBuilder() *listBuilder {
  return &listBuilder{list: &cvrList{}}
}

func (b *listBuilder) WithLabelOption(label string) *listBuilder {
  b.list.labelOption = label
  return b
}

func (b *listBuilder) List() *cvrList{
  if len(b.list.items) != 0 {
    return b.list
  }
  // TODO
  // make a list API call using label as one
  // of the list options
}
```

  - `pkg/volume/cstorpool/v1alpha1/doc.go`
```go
// This namespace caters to cstorpool related operations that
// are a part of overall volume related provisioning.
```

  - `pkg/volume/cstorpool/v1alpha1/cstorpool.go`
```go
import (
  csp "github.com/openebs/maya/pkg/cstorpool/v1alpha1"
  cvr "github.com/openebs/maya/pkg/cstorvolumereplica/v1alpha1"
)

type policyName string

const (
  antiAffinityLabelPolicy policyName = "anti-affinity-label"
  preferAntiAffinityLabelPolicy policyName = "prefer-anti-affinity-label"
)

type policy interface {
  name() policyName
  filter([]string) []string
}

type antiAffinityLabel struct {
  labelSelector string
}

func (p antiAffinityLabel) name() policyName {
  return antiAffinityLabelPolicy
}

func (l antiAffinityLabel) filter(pools []string) []string {
  if l.labelSelector == "" {
    return pools
  }
  exclude := cvr.ListBuilder().WithLabelOption(l.labelSelector).List().GetPoolNames()
  plist := csp.ListBuilder().WithNames(pools).List()
  return plist.FilterNames(csp.IsNotName(exclude...))
}

type preferAntiAffinityLabel struct {
  antiAffinityLabel
}

func (p preferAntiAffinityLabel) name() policyName {
  return preferAntiAffinityLabelPolicy
}

func (p preferAntiAffinityLabel) filter(pools []string) []string {
  plist := p.antiAffinityLabel.Filter(pools)
  if len(plist) > 0 {
    return plist
  }
  return pools
}

type selection struct {
  // list of original pools aginst whom 
  // selection will be made
  pools                []string

  // selection is based on these policies
  policies             []policy
}

// selectionBuildOption is a typed function that
// abstracts configuring a selection instance
type selectionBuildOption func(*selection)

func newSelection(pools []string, opts ...selectionBuildOption) *selection {
  s := &selection{pools: pools}
  for _, o := opts {
    o(s)
  }
  return s
}

// isPolicy determines if the provided policy 
// needs to be considered during selection
func (s *selection) isPolicy(p policyName) bool {
  if len(s.policies) == 0 {
    return false
  }
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
  return s.isPolicy(preferAntiAffinityLabelPolicy)
}

// isAntiAffinityLabel determines if anti affinity
// label needs to be considered during
// selection
func (s *selection) isAntiAffinityLabel() bool {
  return s.isPolicy(antiAffinityLabelPolicy)
}

// PreferAntiAffinityLabel adds anti affinity label
// as a preferred policy to be used during pool 
// selection
func PreferAntiAffinityLabel(lbl string) selectionBuildOption {
  return func(s *selection) {
    p := preferAntiAffinityLabel{labelSelector: lbl}
    s.policies = append(s.policies, p)
  }
}

// AntiAffinityLabel adds anti affinity as a policy
// to be used during pool selection
func AntiAffinityLabel(lbl string) selectionBuildOption {
  return func(s *selection) {
    a := antiAffinityLabel{labelSelector: lbl}
    s.policies = append(s.policies, a)
  }
}

func (s *selection) validate() error {
  if s.isAntiAffinityLabel() && s.isPreferAntiAffinityLabel() {
    return errors.New("invalid selection: antiAffinity and preferAntiAffinity policies can not be together")
  }
  return nil
}

func (s *selection) filter() []string {
  var filtered []string
  if len(s.policies) == 0 {
    return s.pools
  }
  filtered = append(filtered, s.pools...)
  for _, policy := range s.policies {
    filtered = policy.filter(filtered)
  }
}

// Filter will filter the given pools based on the 
// selection options
func Filter(origPools []string, opts ...selectionBuildOption) []string {
  if len(opts) == 0 {
    return origPools
  }
  s := newSelection(origPools, opts...)
  err := s.validate()
  if err != nil {
    // log the error here
    return nil
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
