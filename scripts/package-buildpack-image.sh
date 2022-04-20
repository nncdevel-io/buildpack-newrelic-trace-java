#!/usr/bin/env bash

set -euo pipefail

pack package-buildpack paketo-newrelic-java-agent --config ./package.toml --format image
