#!/bin/sh
set -eu

go mod download
docker build . \
  --build-arg VERSION=0.0.1 \
  --build-arg ENV=local \
  --build-arg PORT=22222 \
  --build-arg WORKDIR=/go/src/github.com/warawara28/sample-books \
  -t sample-books:0.0.1

