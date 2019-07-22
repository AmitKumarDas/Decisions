```go
// Let's try to understand how this package is designed

// https://github.com/cuelang/cue/tree/master/cue
```

### Understand
- Q: Does doc.go offer any help?
  - Package cue creates, evaluates and manipulates CUE configurations.
  - Seems to be the core
- Q: What is cue/types.go?
  - Kind as enum
    - Lot of Kinds e.g. StringKind, IntKind, BottomKind, NullKind, ...
  - A marshalError struct
    - composes its nested errors package
    - can hopefully provide detailed errors with path, position etc details
  - A Iterator to iterate over values
    - This can be marshaled to a json output
    - Uses encoding/json
  - A Value struct that holds any value i.e. string, int, struct, bool, etc
    - Lot of conversion methods against this Value
    - Value to Kind also possible
    - Value is a chain of values
- Q: What is cue/build.go?
