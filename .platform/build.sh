#!/usr/bin/env bash

set -eux
time go get ./...
go generate ./... # exits > 0 if it writes a file
time go build -v -a -ldflags '-s -w' .
ls -l super-potato
