#!/usr/bin/env bash

set -eux
time go get ./...
time go run ./gen --exit-code
time go build -v -a -ldflags '-s -w' .
ls -l super-potato
