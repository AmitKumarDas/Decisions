- How do you think this neat code is even possible?
- How can one abstract http client into such domain specific bits?
- Multiple ways while still being idiomatic

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

```go
// DeleteBatchV1Job deletes a K8s job
func (k *K8sClient) DeleteBatchV1Job(name string) error {
	deletePropagation := mach_apis_meta_v1.DeletePropagationForeground
	return k.cs.BatchV1().Jobs(k.ns).
		Delete(name, &mach_apis_meta_v1.DeleteOptions{PropagationPolicy: &deletePropagation})
}
```
