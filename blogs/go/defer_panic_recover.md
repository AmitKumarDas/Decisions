### Refer
- https://blog.golang.org/defer-panic-and-recover

### Panic
- Panic is a built-in function that stops the ordinary flow of control and begins panicking. 
- When the function F calls panic, execution of F stops, any deferred functions in F are executed normally, 
  - and then F returns to its caller. 
- To the caller, F then behaves like a call to panic. 
- The process continues up the stack until all functions in the current goroutine have returned, 
  - at which point the program crashes. 
- Panics can be initiated by invoking panic directly. 
- They can also be caused by runtime errors, such as out-of-bounds array accesses.

### Recover
- Recover is a built-in function that regains control of a panicking goroutine. 
  - _Used inside go routine only_
- Recover is only useful inside deferred functions. 
  - _Used only whithin a goroutine + defer function_
- During normal execution, a call to recover will return nil and have no other effect. 
- If the current goroutine is panicking, a call to recover will capture the value given to panic and resume normal execution.

### K8s Refer
```go
func TestConflictingAddKnownTypes(t *testing.T) {
  s := runtime.NewScheme()
  gv := schema.GroupVersion{Group: "foo", Version: "v1"}

  panicked := make(chan bool)
  go func() {
  	defer func() {
  		if recover() != nil {
  			panicked <- true
  		}
  	}()
    // Below calls panic in its logic
	  s.AddKnownTypeWithName(gv.WithKind("InternalSimple"), &runtimetesting.InternalSimple{})
	  s.AddKnownTypeWithName(gv.WithKind("InternalSimple"), &runtimetesting.ExternalSimple{})
	  panicked <- false
  }()
  if !<-panicked {
	  t.Errorf("Expected AddKnownTypesWithName to panic with conflicting type registrations")
  }

  go func() {
	  defer func() {
		  if recover() != nil {
			  panicked <- true
		  }
	  }()
    // Below calls panic in its logic
	  s.AddUnversionedTypes(gv, &runtimetesting.InternalSimple{})
	  s.AddUnversionedTypes(gv, &InternalSimple{})
	  panicked <- false
  }()
  if !<-panicked {
  	t.Errorf("Expected AddUnversionedTypes to panic with conflicting type registrations")
  }
}
```
