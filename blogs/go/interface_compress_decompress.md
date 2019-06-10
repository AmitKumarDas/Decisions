## Reference
- http://jmoiron.net/blog/built-in-interfaces

### Problem Statement
Let's take a look at a more complex example now, which uses both interfaces. Suppose we wanted to have a field stored as gzipped compressed text in the database. We'll want to implement the Valuer interface, as before, as a way of automatically compressing data on its way to the database:

```go
type Valuer interface {
  // Value returns a driver Value.
  Value() (Value, error)
}
```

```go
type Scanner interface {
  // Scan assigns a value from a database driver.
  //
  // The src value will be of one of the following restricted
  // set of types:
  //
  //    int64, float64, bool, []byte, string, time.Time
  //    nil - for NULL values
  //
  // An error should be returned if the value can not be stored
  // without loss of information.
  Scan(src interface{}) error
}
```

```go
type GzippedText []byte

func (g GzippedText) Value() (driver.Value, error) {
  b := make([]byte, 0, len(g))
  buf := bytes.NewBuffer(b)
  w := gzip.NewWriter(buf)
  w.Write(g)
  w.Close()
  return buf.Bytes(), nil
}
```

This time, we'll also implement Scanner, to decompress data coming from the database:
```go
func (g *GzippedText) Scan(src interface{}) error {
  var source []byte

  switch src.(type) {
  case string:
    source = []byte(src.(string))
  case []byte:
    source = src.([]byte)
  default:
    return errors.New("incompatible type for GzippedText")
  }

  reader, _ := gzip.NewReader(bytes.NewReader(source))
  defer reader.Close()

  b, err := ioutil.ReadAll(reader)
  if err != nil {
    return err
  }

  *g = GzippedText(b)
  return nil
}
```
