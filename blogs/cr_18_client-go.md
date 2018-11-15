- How do you think this neat code is even possible?
- How can one abstract http client into such domain specific bits?

```go
func (k *K8sClient) GetBatchV1JobAsRaw(name string) (result []byte, err error) {
	return k.cs.BatchV1().RESTClient().
		Get().
		Resource("jobs").
		Name(name).
		VersionedParams(&mach_apis_meta_v1.GetOptions{}, scheme.ParameterCodec).
		DoRaw()
}
```
