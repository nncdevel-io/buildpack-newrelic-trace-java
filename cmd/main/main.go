package main

import (
	"os"

	"github.com/nncdevel-io/paketo-newrelic-trace-java/newrelic"
	"github.com/paketo-buildpacks/libpak"
	"github.com/paketo-buildpacks/libpak/bard"
)

func main() {
	libpak.Main(
		datadog.Detect{},
		datadog.Build{Logger: bard.NewLogger(os.Stdout)},
	)
}
