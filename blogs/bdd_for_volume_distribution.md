### Info
- Version: 2
- Last Updated On: 26-Feb-2019

### High Level Design
- maya/integration-tests/sts/
  - sts_test.go
  - sts_suite_test.go

### Low Level Design
```go
var _ = Describe("StatefulSet", func() {
  var (
	  //
	)

  BeforeEach(func() {
	  //
  })

  AfterEach(func() {
	  //
  })

  Describe("deploy a statefulset", func() {
	  BeforeEach(func() {
      //
	  })

	  It("should distribute the statefulset replicas across pools", func() {
      pvcs, err := pvc.KubeClient().List("", "sts-based-lbl-selectors")
      Ω(err).Should(BeNil())

			upvcs = pvc.ListBuilder().WithListObject(pvcs).AddCheck(pvc.HasName("pvc-name"), pvc.HasName("app-name")).List()
      Ω(upvcs).Should(HaveLen(3))

      cvrs = cvr.KubeClient().List("", "pvc-based-lbl-selectors")
      Ω(cvrs).Should(HaveLen(3))
      
      pools := cvr.ListBuilder().WithListObject(cvrs).RemoveDuplicate().List().ListPoolName()
			Ω(pools).Should(HaveLen(3))      
	  })
  })
})
```
