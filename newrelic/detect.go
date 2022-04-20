package datadog

import (
	"github.com/buildpacks/libcnb"
)

type Detect struct{}

func (d Detect) Detect(context libcnb.DetectContext) (libcnb.DetectResult, error) {
	return libcnb.DetectResult{
		Plans: []libcnb.BuildPlan{
			{
				Provides: []libcnb.BuildPlanProvide{
					{Name: "newrelic-java-agent"},
				},
				Requires: []libcnb.BuildPlanRequire{
					{Name: "jvm-application"},
					{Name: "newrelic-java-agent"},
				},
			},
		},
	}, nil
}
