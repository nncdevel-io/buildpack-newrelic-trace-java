package newrelic

import (
	"fmt"

	"github.com/buildpacks/libcnb"
	"github.com/paketo-buildpacks/libpak"
	"github.com/paketo-buildpacks/libpak/bard"
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
	jk := NewNewRelicAgentJava(depJVMKill, dc)
	result.Layers = append(result.Layers, jk)
	return result, nil
}
