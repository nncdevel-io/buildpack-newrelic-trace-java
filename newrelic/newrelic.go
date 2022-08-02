package newrelic

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/buildpacks/libcnb"
	"github.com/paketo-buildpacks/libpak"
	"github.com/paketo-buildpacks/libpak/bard"
	"github.com/paketo-buildpacks/libpak/sherpa"
)

type NewRelicAgentJava struct {
	LayerContributor libpak.DependencyLayerContributor
	Logger           bard.Logger
}

func NewNewRelicAgentJava(dependency libpak.BuildpackDependency, cache libpak.DependencyCache) NewRelicAgentJava {

	types := libcnb.LayerTypes{
		Build:  true,
		Cache:  true,
		Launch: true,
	}

	layerContributor := libpak.NewDependencyLayerContributor(dependency, cache, types)

	return NewRelicAgentJava{LayerContributor: layerContributor}
}

func (j NewRelicAgentJava) Contribute(layer libcnb.Layer) (libcnb.Layer, error) {
	j.LayerContributor.Logger = j.Logger

	return j.LayerContributor.Contribute(layer, j.dependencyLayerFunc(layer))
}

func (j NewRelicAgentJava) dependencyLayerFunc(layer libcnb.Layer) libpak.DependencyLayerFunc {
	return func(artifact *os.File) (libcnb.Layer, error) {
		j.Logger.Infof("Copying to %s", layer.Path)
		file := filepath.Join(layer.Path, filepath.Base(artifact.Name()))
		if err := sherpa.CopyFile(artifact, file); err != nil {
			return libcnb.Layer{}, fmt.Errorf("unable to copy %s to %s\n%w", artifact.Name(), file, err)
		}

		j.Logger.Infof("Add -javaagent:%s", file)
		layer.LaunchEnvironment.Appendf("JAVA_OPTS", " ", " -javaagent:%s", file)

		layer.Launch = true
		return layer, nil
	}
}

func (NewRelicAgentJava) Name() string {
	return "newrelic-trace-java"
}
