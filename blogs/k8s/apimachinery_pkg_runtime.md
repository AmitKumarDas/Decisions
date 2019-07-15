```go
// https://github.com/kubernetes/kubernetes/tree/master/staging/src/k8s.io/apimachinery/pkg/runtime
```

### Notes
- k8s.io/apimachinery/pkg/runtime/scheme.go -- **tough!**
  - Defines methods to serialize and deserialize API objects
  - Registry to convert GVK to and from Go schemas
  - Foundation for versioned API

- k8s.io/apimachinery/pkg/runtime/scheme_builder.go
  - pretty different way to use builder as a slice of functional options
