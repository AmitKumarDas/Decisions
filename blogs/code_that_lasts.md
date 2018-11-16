### Code That Lasts
The title seems to go against the very nature of change. Nothing lasts for forever. However, mankind will always be in a 
constant disagreement with nature and will try to build fortress around their idea(s).

This holds true for developers writing code and expect it to last for many many releases to come. Not all of them really 
think about this. However, this does not matter since all the developers must go with the flow of change by churning the
code based on the bugs discovered or based on the newly articulated requirements that seem to go against the very 
assumption that was made in the first place.

Is there any hope then to have some pieces of logic that can last?

### Code that is malleable
Remember the mighty trees that cannot bend usually break or get uprooted during natural calamities. Surprisingly, the young 
plants which seem to be more fragile come out unscathed. The question that we need to ask ourselves is _"is your code 
malleable to the forces of change"_ ?

### Devil is in the details
The stuff that I am going to write below needs to change based on my experiences and observations. These details should be 
a live document. These are the details that I believe will _make code more malleable to changes without cracking or being 
prone to bugs_.

#### dated - 16-Nov-2018
- API:
  - Should have well defined payloads also known as API types
  - These API types should be in a well defined namespace/packaging
  - e.g. `pkg/apis/org.group/resource-name/v1alpha1/`
- Entity
  - Should have resources/entities that closely match (or exactly match) these API types
  - These entities should be in their own defined namespace (different than API namespaces)
  - e.g. `pkg/resource-name/v1alpha1/`
  - EntityList i.e. a collection of these entities should be a well defined type
  - refer - https://medium.com/@amit.das/thinking-in-types-3c234eb17680
- EntityHelpers
  - Above entities may have some or all of these:
  - Builder, ListBuilder
  - RestClient,
  - Template e.g. GoTemplate
  - Predicate, Map, Filter, Middleware, etc functions i.e. streaming APIs
  - All of these helpers can reside in the same namespace that is meant for entity
  - i.e. `pkg/resource-name/v1alpha1`
- Unit Testable Entity
  - The entity along with their helpers must be Unit Testable
- Caller Code / Service Code
  - These is the place which ties up the business logic
  - These are the feature implementors
  - These should reside in a namespace of their own
  - e.g. `pkg/installer/v1alpha1` 
  - or `pkg/controller/openebs/v1alpha1`
  - or `cmd/maya-apiserver/server/resource_endpoint.go`
- Unit Testable Caller Code
  - The caller code should also be unit testable
  - However, unit testable caller code is optional
  - Some may want to test this via integration tests
