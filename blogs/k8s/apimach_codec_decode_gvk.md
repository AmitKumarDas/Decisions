```go
// https://github.com/kubernetes/kubernetes/blob/master/staging/src/k8s.io/apimachinery/pkg/runtime/scheme_test.go
```

```go
func TestUnversionedTypes(t *testing.T) {
	internalGV := schema.GroupVersion{Group: "test.group", Version: runtime.APIVersionInternal}
	internalGVK := internalGV.WithKind("Simple")
	externalGV := schema.GroupVersion{Group: "test.group", Version: "testExternal"}
	externalGVK := externalGV.WithKind("Simple")
	otherGV := schema.GroupVersion{Group: "group", Version: "other"}

	scheme := runtime.NewScheme()
	scheme.AddUnversionedTypes(externalGV, &runtimetesting.InternalSimple{})
	scheme.AddKnownTypeWithName(internalGVK, &runtimetesting.InternalSimple{})
	scheme.AddKnownTypeWithName(externalGVK, &runtimetesting.ExternalSimple{})
	scheme.AddKnownTypeWithName(otherGV.WithKind("Simple"), &runtimetesting.ExternalSimple{})

	codec := serializer.NewCodecFactory(scheme).LegacyCodec(externalGV)

	if unv, ok := scheme.IsUnversioned(&runtimetesting.InternalSimple{}); !unv || !ok {
		t.Fatalf("type not unversioned and in scheme: %t %t", unv, ok)
	}

	kinds, _, err := scheme.ObjectKinds(&runtimetesting.InternalSimple{})
	if err != nil {
		t.Fatal(err)
	}
	kind := kinds[0]
	if kind != externalGV.WithKind("InternalSimple") {
		t.Fatalf("unexpected: %#v", kind)
	}

	test := &runtimetesting.InternalSimple{
		TestString: "I'm the same",
	}
	obj := runtime.Object(test)
	data, err := runtime.Encode(codec, obj)
	if err != nil {
		t.Fatal(err)
	}
	obj2, gvk, err := codec.Decode(data, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := obj2.(*runtimetesting.InternalSimple); !ok {
		t.Fatalf("Got wrong type")
	}
	if !reflect.DeepEqual(obj2, test) {
		t.Errorf("Expected:\n %#v,\n Got:\n %#v", test, obj2)
	}
	// object is serialized as an unversioned object (in the group and version it was defined in)
	if *gvk != externalGV.WithKind("InternalSimple") {
		t.Errorf("unexpected gvk returned by decode: %#v", *gvk)
	}

	// when serialized to a different group, the object is kept in its preferred name
	codec = serializer.NewCodecFactory(scheme).LegacyCodec(otherGV)
	data, err = runtime.Encode(codec, obj)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != `{"apiVersion":"test.group/testExternal","kind":"InternalSimple","testString":"I'm the same"}`+"\n" {
		t.Errorf("unexpected data: %s", data)
	}
}
```
