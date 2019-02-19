#### Low Level Implementation
##### [cstor only] - find eligible node(s) to place the replica
  - `pkg/volume/cstorpool/v1alpha1/doc.go`
```go
// This namespace caters to cstorpool related operations that are a part of
// volume related provisioning.
```

  - `pkg/volume/cstorpool/v1alpha1/cstorpool.go`
```go
import (
  csp "github.com/openebs/maya/pkg/cstorpool/v1alpha1"
  cvr "github.com/openebs/maya/pkg/cstorvolumereplica/v1alpha1"
)

type policy int

const (
  // PreferAntiAffinityLabelPolicy when specified will prefer distributing
  // data across pools if possible; else will default to any pool
  PreferAntiAffinityLabelPolicy policy = iota + 1

  // AntiAffinityLabelPolicy when specified will distribute
  // data across pools that are not associated with
  // the label
  AntiAffinityLabelPolicy
)

type selection struct {
  policies              []policy
  antiAffinityLabel string
}

type selectionBuildOption func(*selection)

func newSelection(opts ...selectionBuildOption) *selection {
  s := &selection{}
  for _, o := opts {
    o(s)
  }
  return s
}

func (s *selection) isPolicy(p policy) bool {
  if len(s.policies) == 0 {
    return false
  }
  for _, pol := range s.policies {
    if pol == p {
      return true
    }
  }
  return false
}

func (s *selection) isPreferAntiAffinityLabel() bool {
  return s.isPolicy(PreferAntiAffinityLabelPolicy)
}

func (s *selection) isAntiAffinityLabel() bool {
  return s.isPolicy(AntiAffinityLabelPolicy)
}

// PreferAntiAffinityLabel adds anti affinity label
// as a preferred policy to be used during pool 
// selection
func PreferAntiAffinityLabel(lbl string) selectionBuildOption {
  return func(s *selection) {
    s.antiAffinityLabel = lbl
    s.policies = append(s.policies, PreferAntiAffinityLabelPolicy)
  }
}

// AntiAffinityLabel adds anti affinity as a policy
// to be used during pool selection
func AntiAffinityLabel(lbl string) selectionBuildOption {
  return func(s *selection) {
    s.antiAffinityLabel = lbl
    s.policies = append(s.policies, AntiAffinityLabelPolicy)
  }
}

// Filter will filter the given pools based on the selection
// options
func Filter(origPools []string, opts ...selectionBuildOption) []string {
  var pools []string
  if len(opts) == 0 {
    return origPools
  }
  s := newSelection(opts...)
  if s.isAntiAffinityLabel() || s.isPreferAntiAffinityLabel() {
    pools := ExcludePoolWithLabel(origPools, s.antiAffinityLabel)
  }
  if len(pools) > 0 {
    return pools
  }
  if s.isPreferAntiAffinityLabel() {
    return origPools
  }
  return pools
}

func ExcludePoolWithLabel(origPools []string, lblSelector string) []string {
  if lblSelector == "" {
    return origPools
  }
  pools := csp.ListBuilder().WithNames(origPools).List()
  exclude := cvr.ListBuilder().WithLabel(lblSelector).List().GetPoolNames()
  return pools.Filter(csp.IsNotName(exclude...))
}

// expose below go functions as go template functions
// Filter                             as cspFilter
// AntiAffinityLabel           as cspAntiAffinity
// PreferAntiAffinityLabel as cspPreferAntiAffinity
```
  - `pkg/cstorpool/v1alpha1`
    - define structs & predicates & other functions
  - `pkg/cstorvolumereplica/v1alpha1`
    - define structs & predicates & other functions
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
