### Info
- Version: 3
- Last Updated On: 27-Feb-2019

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
    // apply the sts
    // verify creation of sts instances
    // verify creation of pvc instances
    // verify creation of pvc instances getting bound
  })

  AfterEach(func() {
    // delete the sts
    // delete the pvc
    // verify deletion of sts instances
    // verify deletion of pvc instances
  })

  Describe("test statefulset application on cstor", func() {
    It("should distribute the cstor volume replicas across pools", func() {
      pvcs, _ := pvc.KubeClient().List("", "sts-based-lbl-selectors")

      pvcList = pvc.ListBuilder().WithListObject(pvcs).AddCheck(pvc.HasName("pvc-name"), pvc.HasName("app-name")).List()
      Ω(pvcList).Should(HaveLen(3), "pvc count should be 3")

      cvrs = cvr.KubeClient().List("", "pvc-based-lbl-selectors")
      Ω(cvrs).Should(HaveLen(3), "cvr count should be 3")
      
      poolNames := cvr.ListBuilder().WithListObject(cvrs).List().ListPoolName()
      Ω(poolNames).Should(HaveLen(3), "pool names count should be 3")
      
      pools, _ := cstorpool.KubeClient().List("", "lbl-selector-based-on-poolnames")
      nodeNames = cstorpool.Listbuilder().WithListObject(pools).List().ListNodeName()
      Ω(nodeNames).Should(haveLen(3), "node names count should be 3")
    })
    
    It("should co-locate the cstor volume targets with application instances", func() {
      // future
    })
  })
})
```
