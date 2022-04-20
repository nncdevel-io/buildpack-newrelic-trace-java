package datadog

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/buildpacks/libcnb"
	"github.com/paketo-buildpacks/libpak"
	"github.com/paketo-buildpacks/libpak/bard"
	"github.com/paketo-buildpacks/libpak/sherpa"
)

// Build foo bar
type Build struct {
	Logger bard.Logger
}

// Build foo bar
func (b Build) Build(context libcnb.BuildContext) (libcnb.BuildResult, error) {
	result := libcnb.NewBuildResult()

	cr, err := libpak.NewConfigurationResolver(context.Buildpack, &b.Logger)
	if err != nil {
		return libcnb.BuildResult{}, fmt.Errorf("unable to create configuration resolver\n%w", err)
	}

	dr, err := libpak.NewDependencyResolver(context)
	if err != nil {
		return libcnb.BuildResult{}, fmt.Errorf("unable to create dependency resolver\n%w", err)
	}

	dc, err := libpak.NewDependencyCache(context)
	if err != nil {

		return libcnb.BuildResult{}, fmt.Errorf("unable to create dependency cache\n%w", err)
	}

	v, _ := cr.Resolve("BP_NEWRELIC_AGENT_VERSION")

	depJVMKill, err := dr.Resolve("newrelic-java-agent", v)
	if err != nil {
		b.Logger.Infof("error: %w", err)
		return libcnb.BuildResult{}, fmt.Errorf("unable to find dependency\n%w", err)
	}
	jk := NewNewRelicAgentJava(depJVMKill, dc, result.Plan)
	result.Layers = append(result.Layers, jk)
	return result, nil
}

type NewRelicAgentJava struct {
	LayerContributor libpak.DependencyLayerContributor
	Logger           bard.Logger
}

func NewNewRelicAgentJava(dependency libpak.BuildpackDependency, cache libpak.DependencyCache, plan *libcnb.BuildpackPlan) NewRelicAgentJava {
	return NewRelicAgentJava{LayerContributor: libpak.NewDependencyLayerContributor(dependency, cache, plan)}
}

func (j NewRelicAgentJava) Contribute(layer libcnb.Layer) (libcnb.Layer, error) {
	j.LayerContributor.Logger = j.Logger

	return j.LayerContributor.Contribute(layer, func(artifact *os.File) (libcnb.Layer, error) {
		j.Logger.Infof("Copying to %s", layer.Path)
		file := filepath.Join(layer.Path, filepath.Base(artifact.Name()))
		if err := sherpa.CopyFile(artifact, file); err != nil {
			return libcnb.Layer{}, fmt.Errorf("unable to copy %s to %s\n%w", artifact.Name(), file, err)
		}

		j.Logger.Infof("Add -javaagent:%s", file)
		layer.LaunchEnvironment.Appendf("JAVA_OPTS", " -javaagent:%s", file)

		layer.Launch = true
		return layer, nil
	})
}

func (NewRelicAgentJava) Name() string {
	return "newrelic-trace-java"
}
