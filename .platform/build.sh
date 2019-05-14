#!/usr/bin/env bash

set -eux
time go get .
time go build -v -a -ldflags '-s -w' .
ls -l super-potato
