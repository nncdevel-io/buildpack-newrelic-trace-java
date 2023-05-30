package newrelic

import (
	"fmt"

	"github.com/buildpacks/libcnb"
	"github.com/paketo-buildpacks/libpak"
	"github.com/paketo-buildpacks/libpak/bard"
)

type Detect struct {
	Logger bard.Logger
}

func (d Detect) Detect(context libcnb.DetectContext) (libcnb.DetectResult, error) {
	d.Logger.Title(context.Buildpack)

	cr, err := libpak.NewConfigurationResolver(context.Buildpack, &d.Logger)
	if err != nil {
		return libcnb.DetectResult{}, fmt.Errorf("unable to create configuration resolver\n%w", err)
	}

	v, _ := cr.Resolve("BP_USE_NEWRELIC")

	if v == "true" {
		d.Logger.Debugf("Newrelic Agent enabled.")
		return libcnb.DetectResult{
			Plans: []libcnb.BuildPlan{
				{
					Pass: true,
					Provides: []libcnb.BuildPlanProvide{
						{Name: "newrelic-java-agent"},
					},
					Requires: []libcnb.BuildPlanRequire{
						{Name: "newrelic-java-agent"},
						{Name: "jvm-application"},
					},
				},
			},
		}, nil
	} else {
		return libcnb.DetectResult{
			Plans: []libcnb.BuildPlan{
				{
					Provides: []libcnb.BuildPlanProvide{},
					Requires: []libcnb.BuildPlanRequire{},
				},
			},
		}, nil
	}

}
