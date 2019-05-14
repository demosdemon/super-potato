#!/usr/bin/env bash

set -eux
export CGO_ENABLED=0
time go build -v -a -ldflags '-s -w -extldflags "-static"' .
ls -l super-potato
