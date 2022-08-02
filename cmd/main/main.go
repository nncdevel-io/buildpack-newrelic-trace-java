package main

import (
	"os"

	"github.com/nncdevel-io/buildpack-newrelic-trace-java/newrelic"
	"github.com/paketo-buildpacks/libpak"
	"github.com/paketo-buildpacks/libpak/bard"
)

func main() {
	logger := bard.NewLogger(os.Stdout)
	libpak.Main(
		newrelic.Detect{Logger: logger},
		newrelic.Build{Logger: logger},
	)
}
