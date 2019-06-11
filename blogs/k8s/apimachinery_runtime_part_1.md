```go
// k8s.io/apimachinery/pkg/runtime/codec_test.go

import (
  "k8s.io/apimachinery/pkg/runtime/schema"
)

func gv(group, version string) schema.GroupVersion {
  return schema.GroupVersion{Group: group, Version: version}
}

func gvk(group, version, kind string) schema.GroupVersionKind {
  return schema.GroupVersionKind{Group: group, Version: version, Kind: kind}
}

func gk(group, kind string) schema.GroupKind {
  return schema.GroupKind{Group: group, Kind: kind}
}

// usage
//
// gv("mygroup", "__internal")
// gk("mygroup", "Foo")
// gvk("yetanother", "v1", "Baz")
// gvk("mygroup", "__internal", "Bar")
```

```go
import (
  "k8s.io/apimachinery/pkg/runtime"
  "k8s.io/apimachinery/pkg/runtime/schema"
  runtimetesting "k8s.io/apimachinery/pkg/runtime/testing"
)

  internalGV := schema.GroupVersion{Group: "test.group", Version: runtime.APIVersionInternal}
  externalGV := schema.GroupVersion{Group: "test.group", Version: "external"}
  
  scheme := runtime.NewScheme()
  
  // add these schemas to the runtime scheme
  scheme.AddKnownTypeWithName(internalGV.WithKind("Complex"), &runtimetesting.InternalComplex{})
  scheme.AddKnownTypeWithName(externalGV.WithKind("Complex"), &runtimetesting.ExternalComplex{})
```

```go
// k8s.io/apimachinery/pkg/runtime/converter.go

// UnstructuredConverter is an interface for converting between interface{}
// and map[string]interface representation.
type UnstructuredConverter interface {
  ToUnstructured(obj interface{}) (map[string]interface{}, error)
  FromUnstructured(u map[string]interface{}, obj interface{}) error
}
```

```go
// k8s.io/apimachinery/pkg/runtime/embedded_test.go
import (
  "k8s.io/apimachinery/pkg/runtime"
  "k8s.io/apimachinery/pkg/runtime/schema"
  "k8s.io/apimachinery/pkg/runtime/serializer"
)

  internalGV := schema.GroupVersion{Group: "test.group", Version: runtime.APIVersionInternal}
  externalGV := schema.GroupVersion{Group: "test.group", Version: "v1test"}
  externalGVK := externalGV.WithKind("ObjectTest")
	
  s := runtime.NewScheme()
  s.AddKnownTypes(internalGV, &runtimetesting.ObjectTest{})
  s.AddKnownTypeWithName(externalGVK, &runtimetesting.ObjectTestExternal{})
	
  codec := serializer.NewCodecFactory(s).LegacyCodec(externalGV)
	
  obj, gvk, err := codec.Decode([]byte(
    `{"kind":"`+externalGVK.Kind+`","apiVersion":"`+externalGV.String()+`","items":[{}]}`
    ), 
    nil, 
    nil,
  )

  if *gvk != externalGVK {
    t.Fatalf("unexpected kind: %#v", gvk)
  }
```

```go
import (
  "k8s.io/apimachinery/pkg/runtime"
  "k8s.io/apimachinery/pkg/runtime/schema"
  "k8s.io/apimachinery/pkg/runtime/serializer"
)

  internalGV := schema.GroupVersion{Group: "test.group", Version: runtime.APIVersionInternal}
  externalGV := schema.GroupVersion{Group: "test.group", Version: "v1test"}
  embeddedTestExternalGVK := externalGV.WithKind("EmbeddedTest")

  s := runtime.NewScheme()
  s.AddKnownTypes(internalGV, &runtimetesting.EmbeddedTest{})
  s.AddKnownTypeWithName(embeddedTestExternalGVK, &runtimetesting.EmbeddedTestExternal{})

  codec := serializer.NewCodecFactory(s).LegacyCodec(externalGV)

  inner := &runtimetesting.EmbeddedTest{
    ID: "inner",
  }
  outer := &runtimetesting.EmbeddedTest{
    ID:     "outer",
    Object: runtime.NewEncodable(codec, inner),
  }

  wire, err := runtime.Encode(codec, outer)
  if err != nil {
    t.Fatalf("Unexpected encode error '%v'", err)
  }

  decoded, err := runtime.Decode(codec, wire)
  if err != nil {
    t.Fatalf("Unexpected decode error %v", err)
  }
```
