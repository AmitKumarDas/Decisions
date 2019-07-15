```go
// k8s.io/apimachinery/pkg/runtime/scheme_builder.go
```

### Notes
- How it moves between type to pointer & vice versa
- AddToScheme might perhaps have been better named as `Build`?

```go
func (sb *SchemeBuilder) Build(s *Scheme) error {}
```

- It same thing just different technique
  - This does not focus on creating a new entity
  - Entity is passed & can be mutated based on the registered functional options

```go
package runtime

// SchemeBuilder collects functions that add things to a scheme. It's to allow
// code to compile without explicitly referencing generated types. You should
// declare one in each package that will have generated deep copy or conversion
// functions.
type SchemeBuilder []func(*Scheme) error

// AddToScheme applies all the stored functions to the scheme. A non-nil error
// indicates that one function failed and the attempt was abandoned.
func (sb *SchemeBuilder) AddToScheme(s *Scheme) error {
	for _, f := range *sb {
		if err := f(s); err != nil {
			return err
		}
	}
	return nil
}

// Register adds a scheme setup function to the list.
func (sb *SchemeBuilder) Register(funcs ...func(*Scheme) error) {
	for _, f := range funcs {
		*sb = append(*sb, f)
	}
}

// NewSchemeBuilder calls Register for you.
func NewSchemeBuilder(funcs ...func(*Scheme) error) SchemeBuilder {
	var sb SchemeBuilder
	sb.Register(funcs...)
	return sb
}
```
