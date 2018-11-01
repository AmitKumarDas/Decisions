### Why should we use reconciliation to implement e2e?
If we peel off the jargons of end-to-end testing, it boils down to observing if given entities are in proper states or not.
In addition, to observation, testing can induce failures and then repeat the same observations of those entities.

If we can relate this task of observation with watching of kubernetes resources, I guess we can build e2e test cases as 
kubernetes reconciliation loops.

#### Approaches
- 1/ Entire test logic inside the reconcile
  - Make use of test specification i.e. a k8s custom resource
  - Test spec will determine the resource to be set
  - Test spec will determine the errors to be injected
  - Entire test logic will be implemented as part of reconcile
  - Test spec status will be updated based on the current state of resources
  - Optional:
    - Make use of Ginkgo to trigger reconciliation
    - No need of web service
    - Can be run either within or outside container
- Advantages:
  - Kind of forever monitor operator if not managed by ginkgo
  - Spec will reflect the current state of e2e test
  - This seems to be in line with chaos engineering principles

### References
- https://github.com/kubernetes-sigs/controller-runtime/blob/master/pkg/builder/build_test.go
