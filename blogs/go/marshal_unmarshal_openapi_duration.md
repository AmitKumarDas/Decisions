```go
// https://github.com/kubernetes/apimachinery/blob/master/pkg/apis/meta/v1/duration.go
```

```
I have a weird issue, I don't know if its a bug or just my code doing the wrong thing. I'm using `scheme.Codecs.UniversalDeserializer()` to parse some YAML to an object struct and we found it 
doesn't understand time.Duration fields. It wants the field to be an int, not the usual string 
parsing behavior. How do I get the same behavior as places like kubectl?


chancez [4:32 AM]
@coderanger are you using `time.Duration` or `meta/v1.Duration`?
the latter should work when i look at it's UnmarshalJSON method
https://github.com/kubernetes/apimachinery/blob/master/pkg/apis/meta/v1/duration.go (edited) 
```
