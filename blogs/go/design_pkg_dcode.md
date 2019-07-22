```go
// https://github.com/go-functional/dcode
```

### Review
- bool.go
  - func Bool() Decoder
  - Creates a new instance of Decoder
  - It does a json.Unmarshal
  - Makes use of encoding/json
- int.go, string.go, float64.go
  - Similar to bool.go
- NOTE on PERFORMANCE:
  - If one needs to parse multiple fileds from a JSON, then it will lead to **multiple** json.Unmarshal calls
  - Below example leads to
    - 4 json.Unmarshal calls & 
    - 3 json.Marshal calls
  ```go
    dcoder := Field("field1", Field("field2", Field("field3", Int())))
    var i int
    Decode(dcoder, jsonBytes, &i)
  ```
- NOTE on readability
  - Pure functional way vs. loop
  ```go
  dcoder := Field("first", Field("second", Field("third", Int())))
  ```
  - Make it readable via fluent API i.e. builder pattern
  ```go
  dcoder := First("first").Then("second").Then("third").Into(Int())
  ```
