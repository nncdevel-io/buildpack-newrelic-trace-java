package newrelic

import (
	"testing"

	"github.com/buildpacks/libcnb"
	. "github.com/onsi/gomega"
	"github.com/sclevine/spec"
)

func testBuild(t *testing.T, context spec.G, it spec.S) {
	var (
		Expect = NewWithT(t).Expect
		ctx    libcnb.BuildContext
	)

	it.Before(func() {
		ctx.StackID = "test-stack-id"

		ctx.Buildpack.Info.Name = "nncdevel-io/paketo-newrelic-java-agent"

		ctx.Buildpack.Metadata = map[string]interface{}{
			"configurations": []map[string]interface{}{
				{
					"name":    "BP_NEWRELIC_AGENT_VERSION",
					"default": "0.0.0",
					"launch":  "true",
				},
			},
			"dependencies": []map[string]interface{}{
				{
					"id":      "newrelic-java-agent",
					"version": "0.0.0",
					"stacks":  []interface{}{"test-stack-id"},
					"uri":     "https://repo1.maven.org/maven2/com/newrelic/agent/java/newrelic-agent/7.9.0/newrelic-agent-7.9.0.jar",
				},
			},
		}
	})

	it("Build Success", func() {
		build := Build{}

		buildResult, err := build.Build(ctx)

		Expect(err).NotTo(HaveOccurred())
		Expect(buildResult).NotTo(BeNil())
	})

	it("Failure: Dependency Not Found", func() {

		ctx.Buildpack.Metadata = map[string]interface{}{
			"configurations": []map[string]interface{}{
				{
					"name":    "BP_NEWRELIC_AGENT_VERSION",
					"default": "0.0.0",
					"launch":  "true",
				},
			},
			"dependencies": []map[string]interface{}{},
		}

		build := Build{}

		buildResult, err := build.Build(ctx)

		Expect(err).To(HaveOccurred())
		Expect(buildResult).NotTo(BeNil())
	})
}
