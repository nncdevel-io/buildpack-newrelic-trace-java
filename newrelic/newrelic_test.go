package newrelic

import (
	"os"
	"testing"

	"github.com/buildpacks/libcnb"
	. "github.com/onsi/gomega"
	"github.com/paketo-buildpacks/libpak"
	"github.com/sclevine/spec"
)

func testNewrelic(t *testing.T, context spec.G, it spec.S) {

	var (
		Expect = NewWithT(t).Expect
		ctx    libcnb.BuildContext
	)

	it("NewNewrelicAgent", func() {

		dep := libpak.BuildpackDependency{
			Version: "11.0.0",
			URI:     "https://localhost/stub-jre-11.tar.gz",
			SHA256:  "3aa01010c0d3592ea248c8353d60b361231fa9bf9a7479b4f06451fef3e64524",
		}

		dc := libpak.DependencyCache{CachePath: "testdata"}

		dd := NewNewRelicAgentJava(dep, dc)

		Expect(dd).NotTo(BeNil())

	})

	it("LayerName", func() {
		dep := libpak.BuildpackDependency{
			Version: "11.0.0",
			URI:     "https://localhost/stub-jre-11.tar.gz",
			SHA256:  "3aa01010c0d3592ea248c8353d60b361231fa9bf9a7479b4f06451fef3e64524",
		}

		dc := libpak.DependencyCache{CachePath: "testdata"}

		dd := NewNewRelicAgentJava(dep, dc)

		Expect(dd.Name()).To(Equal("newrelic-trace-java"))
	})

	it("dependencyLayerFunc", func() {

		depFile, _ := os.CreateTemp("", "dependency")

		dep := libpak.BuildpackDependency{
			Version: "11.0.0",
			URI:     "https://localhost/stub-jre-11.tar.gz",
			SHA256:  "3aa01010c0d3592ea248c8353d60b361231fa9bf9a7479b4f06451fef3e64524",
		}

		dc := libpak.DependencyCache{CachePath: "testdata"}

		nr := NewNewRelicAgentJava(dep, dc)

		layer, layerErr := ctx.Layers.Layer("dd-trace-java")
		Expect(layerErr).NotTo(HaveOccurred())

		f := nr.dependencyLayerFunc(layer)

		layer, layerErr = f(depFile)

		Expect(layer.LaunchEnvironment["JAVA_OPTS.append"]).To(HavePrefix(" -javaagent:dd-trace-java/dependency"))

		Expect(layerErr).NotTo(HaveOccurred())
		Expect(layer.Build).To(BeFalse())
		Expect(layer.Cache).To(BeFalse())
		Expect(layer.Launch).To(BeTrue())

		Expect(nr.Name()).To(Equal("newrelic-trace-java"))
	})
}
