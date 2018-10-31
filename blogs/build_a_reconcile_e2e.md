### Why controller runtime to build e2e?
If we peel off the jargons of end-to-end testing, it boils down to observing if given entities are in proper states or not.
In addition, to observation, testing can induce failures and then repeat the same observations of those entities.

If we can relate this task of observation with watching of kubernetes resources, I guess we can build e2e test cases as 
kubernetes reconciliation loops.

#### Approaches
- 1 - Entire test logic inside the reconcile
  - Make use of test spec i.e. a k8s custom resource
  - Write entire test logic in the reconcile
  - Update the test spec inside the reconcile
  - [optional] May make use of Ginkgo to be run & be managed via it's `Done` channel. Control will be with Ginkgo.
- Advantages:
  - Kind of forever monitor operator if not managed by ginkgo
  - Spec will reflect the current state of e2e test

#### Scenario 1
- Apply a e2e specification to set up one or more resources
- Set the reconciler set up the resource(s) based on the spec's status
- Run the reconciler to observe resource(s)' state based on the spec's status

#### Scenario 2
- Apply a e2e specification to set up one or more resources & then induce failures
- Reconciler to set up the resource(s) based on the spec's status
- Reconciler to induce failures based on the spec's status
- Reconciler to observe resource(s)' state based on the spec's status

### References
- https://github.com/kubernetes-sigs/controller-runtime/blob/master/pkg/builder/build_test.go
