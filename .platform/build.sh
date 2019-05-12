#!/usr/bin/env bash

set -eux
time go get ./...
time go generate ./...
time go build -v -a -ldflags '-s -w -extldflags "-static"' .
ls -l super-potato
