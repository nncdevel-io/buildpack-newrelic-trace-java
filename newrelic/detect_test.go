package newrelic

import (
	"testing"

	"github.com/buildpacks/libcnb"
	. "github.com/onsi/gomega"
	"github.com/sclevine/spec"
)

func testDetect(t *testing.T, context spec.G, it spec.S) {
	var (
		Expect = NewWithT(t).Expect
		ctx    libcnb.DetectContext
	)

	it.Before(func() {

		ctx = libcnb.DetectContext{
			Buildpack: libcnb.Buildpack{
				Metadata: map[string]interface{}{},
			},
		}
	})

	it("BP_USE_NEWRELIC unset", func() {

		detect := Detect{}

		res, err := detect.Detect(ctx)
		Expect(err).NotTo(HaveOccurred())
		Expect(res).NotTo(BeNil())

		Expect(res.Plans).To(Equal([]libcnb.BuildPlan{
			{
				Provides: []libcnb.BuildPlanProvide{},
				Requires: []libcnb.BuildPlanRequire{},
			},
		}))

	})

	it("BP_USE_NEWRELIC is false", func() {

		t.Setenv("BP_USE_NEWRELIC", "false")

		detect := Detect{}

		res, err := detect.Detect(ctx)
		Expect(err).NotTo(HaveOccurred())
		Expect(res).NotTo(BeNil())

		Expect(res.Plans).To(Equal([]libcnb.BuildPlan{
			{
				Provides: []libcnb.BuildPlanProvide{},
				Requires: []libcnb.BuildPlanRequire{},
			},
		}))

	})

	it("BP_USE_NEWRELIC is true", func() {

		t.Setenv("BP_USE_NEWRELIC", "true")

		detect := Detect{}

		res, err := detect.Detect(ctx)
		Expect(err).NotTo(HaveOccurred())
		Expect(res).NotTo(BeNil())

		Expect(res.Plans).To(Equal([]libcnb.BuildPlan{
			{
				Provides: []libcnb.BuildPlanProvide{
					{Name: "newrelic-java-agent"},
				},
				Requires: []libcnb.BuildPlanRequire{
					{Name: "newrelic-java-agent"},
					{Name: "jvm-application"},
				},
			},
		}))

	})
}
