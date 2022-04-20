#!/usr/bin/env bash

set -euo pipefail

if [[ -d ../go-cache ]]; then
  GOPATH=$(realpath ../go-cache)
  export GOPATH
fi

GOOS="linux" GOARCH="amd64" go build -ldflags='-s -w' -o bin/main github.com/nncdevel-io/paketo-newrelic-trace-java/cmd/main
ln -fs main bin/build
#ln -fs main bin/detect

touch bin/detect

echo build succeeded.
