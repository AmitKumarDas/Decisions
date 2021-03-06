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
- Pure function No Interface
  ```go
  // decoder.go
  
  // note: 
  //  name suffix with er typically points to
  // interface but this is struct
  type Decoder struct {
    call func([]byte) (interface{}, error)
  }
  
  // note:
  //  when pure functional struct like above comes up
  // it is usually complimented by an instantiation
  // function i.e. constructor
  newDecoder(fn func([]byte) (interface{}, error)) Decoder {
    return Decoder{call: fn}
  }
  ```
  ```go
  // decode.go
  
  // Decode decodes b into i using d
  func Decode(b []byte, d Decoder, i interface{}) error {
    ret, err := d.call(b)
    if err != nil {
      return err
    }
    inputType := reflect.TypeOf(i)
    if inputType.Kind() != reflect.Ptr {
      return fmt.Errorf("Given %s type is not a pointer", inputType.Name())
    }
    inputTypeName := inputType.Elem().Name()
    returnedType := reflect.TypeOf(ret)
    if inputTypeName != returnedType.Name() {
      return fmt.Errorf(
        "Got type '%s', but expected '%s'",
        returnedType.Name(),
        inputType.Name(),
      )
    }
    // Need to call Elem to traverse the pointer
    inputVal := reflect.ValueOf(i).Elem()
    retVal := reflect.ValueOf(ret)
    if !inputVal.CanSet() {
      return fmt.Errorf("Couldn't set the decoded value")
    }
    inputVal.Set(retVal)
    return nil
  }
  ```
  ```go
  // int.go
  
  // Int decodes any JSON field into an integer
  func Int() Decoder {
    return newDecoder(func(b []byte) (interface{}, error) {
      var i int
      if err := json.Unmarshal(b, &i); err != nil {
        return 0, err
      }
      return i, nil
    })
  }
  ```
