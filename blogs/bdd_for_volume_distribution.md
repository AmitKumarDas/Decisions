### Info
- Version: 1
- Last Updated On: 25-Feb-2019

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
			upvcs = pvc.ListBuilder().
        WithListObject(pvcs).
        AddCheck(pvc.HasName("pvc-name"), pvc.HasName("app-name")).List()
      Ω(upvcs).Should(HaveLen(3))
      cvrs = cvr.KubeClient().List("", "pvc-based-lbl-selectors")
      Ω(cvrs).Should(HaveLen(3))
      pools := cvr.ListBuilder().WithListObject(cvrs).RemoveDuplicate().List().ListPoolName()
			Ω(pools).Should(HaveLen(3))      
		})
	})

	Describe("when the BeforeSuite fails", func() {
		var beforeSuite *types.SetupSummary

		BeforeEach(func() {
			beforeSuite = &types.SetupSummary{
				State:   types.SpecStateFailed,
				RunTime: 3 * time.Second,
				Failure: types.SpecFailure{
					Message:               "failed to setup",
					ComponentCodeLocation: codelocation.New(0),
					Location:              codelocation.New(2),
				},
			}
			reporter.BeforeSuiteDidRun(beforeSuite)

			reporter.SpecSuiteDidEnd(&types.SuiteSummary{
				NumberOfSpecsThatWillBeRun: 1,
				NumberOfFailedSpecs:        1,
				RunTime:                    testSuiteTime,
			})
		})

		It("should record the test as having failed", func() {
			output := readOutputFile()
			Ω(output.Name).Should(Equal("My test suite"))
			Ω(output.Tests).Should(Equal(1))
			Ω(output.Failures).Should(Equal(1))
			Ω(output.Time).Should(Equal(reportedSuiteTime))
			Ω(output.Errors).Should(Equal(0))
			Ω(output.TestCases[0].Name).Should(Equal("BeforeSuite"))
			Ω(output.TestCases[0].Time).Should(Equal(3.0))
			Ω(output.TestCases[0].ClassName).Should(Equal("My test suite"))
			Ω(output.TestCases[0].FailureMessage.Type).Should(Equal("Failure"))
			Ω(output.TestCases[0].FailureMessage.Message).Should(ContainSubstring("failed to setup"))
			Ω(output.TestCases[0].FailureMessage.Message).Should(ContainSubstring(beforeSuite.Failure.ComponentCodeLocation.String()))
			Ω(output.TestCases[0].FailureMessage.Message).Should(ContainSubstring(beforeSuite.Failure.Location.String()))
			Ω(output.TestCases[0].Skipped).Should(BeNil())
		})
	})

	Describe("when the AfterSuite fails", func() {
		var afterSuite *types.SetupSummary

		BeforeEach(func() {
			afterSuite = &types.SetupSummary{
				State:   types.SpecStateFailed,
				RunTime: 3 * time.Second,
				Failure: types.SpecFailure{
					Message:               "failed to setup",
					ComponentCodeLocation: codelocation.New(0),
					Location:              codelocation.New(2),
				},
			}
			reporter.AfterSuiteDidRun(afterSuite)

			reporter.SpecSuiteDidEnd(&types.SuiteSummary{
				NumberOfSpecsThatWillBeRun: 1,
				NumberOfFailedSpecs:        1,
				RunTime:                    testSuiteTime,
			})
		})

		It("should record the test as having failed", func() {
			output := readOutputFile()
			Ω(output.Name).Should(Equal("My test suite"))
			Ω(output.Tests).Should(Equal(1))
			Ω(output.Failures).Should(Equal(1))
			Ω(output.Time).Should(Equal(reportedSuiteTime))
			Ω(output.Errors).Should(Equal(0))
			Ω(output.TestCases[0].Name).Should(Equal("AfterSuite"))
			Ω(output.TestCases[0].Time).Should(Equal(3.0))
			Ω(output.TestCases[0].ClassName).Should(Equal("My test suite"))
			Ω(output.TestCases[0].FailureMessage.Type).Should(Equal("Failure"))
			Ω(output.TestCases[0].FailureMessage.Message).Should(ContainSubstring("failed to setup"))
			Ω(output.TestCases[0].FailureMessage.Message).Should(ContainSubstring(afterSuite.Failure.ComponentCodeLocation.String()))
			Ω(output.TestCases[0].FailureMessage.Message).Should(ContainSubstring(afterSuite.Failure.Location.String()))
			Ω(output.TestCases[0].Skipped).Should(BeNil())
		})
	})

	specStateCases := []struct {
		state   types.SpecState
		message string
	}{
		{types.SpecStateFailed, "Failure"},
		{types.SpecStateTimedOut, "Timeout"},
		{types.SpecStatePanicked, "Panic"},
	}

	for _, specStateCase := range specStateCases {
		specStateCase := specStateCase
		Describe("a failing test", func() {
			var spec *types.SpecSummary
			BeforeEach(func() {
				spec = &types.SpecSummary{
					ComponentTexts: []string{"[Top Level]", "A", "B", "C"},
					State:          specStateCase.state,
					RunTime:        5 * time.Second,
					Failure: types.SpecFailure{
						ComponentCodeLocation: codelocation.New(0),
						Location:              codelocation.New(2),
						Message:               "I failed",
					},
				}
				reporter.SpecWillRun(spec)
				reporter.SpecDidComplete(spec)

				reporter.SpecSuiteDidEnd(&types.SuiteSummary{
					NumberOfSpecsThatWillBeRun: 1,
					NumberOfFailedSpecs:        1,
					RunTime:                    testSuiteTime,
				})
			})

			It("should record test as failing", func() {
				output := readOutputFile()
				Ω(output.Name).Should(Equal("My test suite"))
				Ω(output.Tests).Should(Equal(1))
				Ω(output.Failures).Should(Equal(1))
				Ω(output.Time).Should(Equal(reportedSuiteTime))
				Ω(output.Errors).Should(Equal(0))
				Ω(output.TestCases[0].Name).Should(Equal("A B C"))
				Ω(output.TestCases[0].ClassName).Should(Equal("My test suite"))
				Ω(output.TestCases[0].FailureMessage.Type).Should(Equal(specStateCase.message))
				Ω(output.TestCases[0].FailureMessage.Message).Should(ContainSubstring("I failed"))
				Ω(output.TestCases[0].FailureMessage.Message).Should(ContainSubstring(spec.Failure.ComponentCodeLocation.String()))
				Ω(output.TestCases[0].FailureMessage.Message).Should(ContainSubstring(spec.Failure.Location.String()))
				Ω(output.TestCases[0].Skipped).Should(BeNil())
			})
		})
	}

	for _, specStateCase := range []types.SpecState{types.SpecStatePending, types.SpecStateSkipped} {
		specStateCase := specStateCase
		Describe("a skipped test", func() {
			var spec *types.SpecSummary
			BeforeEach(func() {
				spec = &types.SpecSummary{
					ComponentTexts: []string{"[Top Level]", "A", "B", "C"},
					State:          specStateCase,
					RunTime:        5 * time.Second,
				}
				reporter.SpecWillRun(spec)
				reporter.SpecDidComplete(spec)

				reporter.SpecSuiteDidEnd(&types.SuiteSummary{
					NumberOfSpecsThatWillBeRun: 1,
					NumberOfFailedSpecs:        0,
					RunTime:                    testSuiteTime,
				})
			})

			It("should record test as failing", func() {
				output := readOutputFile()
				Ω(output.Tests).Should(Equal(1))
				Ω(output.Failures).Should(Equal(0))
				Ω(output.Time).Should(Equal(reportedSuiteTime))
				Ω(output.Errors).Should(Equal(0))
				Ω(output.TestCases[0].Name).Should(Equal("A B C"))
				Ω(output.TestCases[0].Skipped).ShouldNot(BeNil())
			})
		})
	}
})
```
