### Why controller runtime to build e2e?
If we peel off the jargons of end-to-end testing, it boils down to observing if given entities are in proper states or not.
In addition, to observation, testing can induce failures and then repeat the same observations of those entities.

If we can relate this power of observation with watching of resources and if these resources are kubernetes based, I guess we
can gain by building e2e test cases via reconciliation loops.

#### Approaches
- 1 - Entire test logic inside the reconcile
  - Make use of test spec i.e. a k8s custom resource
  - Write entire test logic in the reconcile
  - Update the test spec inside the reconcile
  - [optional] May make use of Ginkgo to be run & be managed via it's `Done` channel. Control will be with Ginkgo.
- Advantages:
  - Kind of forever monitor operator

- 2 - Use reconcile to watch while verify using Ginkgo
  - No need of any k8s custom resource
  - Set the watch logic via builder
  - Implement the reconcile function to make use of channels that just send the reconcile.Request objects
  - Start the controller
  - Use a forever loop in ginkgo body to verify data received from channels
  - Use ginkgo's Done channel to stop the test
- Advantages:
  - No need of complex polling logic

#### Scenario 1
- Apply a e2e specification to set up one or more resources
- Set the reconciler set up the resource(s) based on the spec's status
- Run the reconciler to observe resource(s)' state based on the spec's status
- Possible state transitions:
  - Init -> Ready -- this will set up the resources
  - Ready -> Complete -- this will observe the resources and all resources states matched the expectation
  - Complete -> Error -- this will observe the resources and some / all resources states did not match the expectation
  - Complete -> Complete --this will observe the resources and all resources states matched the expectation
  - Error -> Complete --this will observe the resources and all resources states matched the expectation

#### Scenario 2
- Apply a e2e specification to set up one or more resources & then induce failures
- Set the reconciler set up the resource(s) based on the spec's status
- Set the reconciler to induce failures based on the spec's status
- Run the reconciler to observe resource(s)' state based on the spec's status
- Possible state transitions:
  - Init -> Ready -- this will set up the resources once per lifetime
  - Ready -> Chaos_Complete -- this will inject failures once per lifetime
  - Chaos_Complete -> Complete -- this will observe the resources and all resources states matched the expectation
  - Chaos_Complete -> Error -- this will observe the resources and some / all resources states did not match the expectation
  - Complete -> Complete -- this will observe the resources and all resources states matched the expectation
  - Error -> Complete -- this will observe the resources and all resources states matched the expectation
  - Complete -> Error -- this will observe the resources and some / all resources states did not match the expectation
  - Error -> Error -- this will observe the resources and some / all resources states did not match the expectation


### NOTES
- Gingkgo has `Done` channel to trigger stopping a test
- You might need to use custom channels extensively to move data from a reconcile function to outside i.e. ginkgo body
- The things that you would want to verify will be inside ginkgo body


### References
- https://github.com/kubernetes-sigs/controller-runtime/blob/master/pkg/builder/build_test.go
