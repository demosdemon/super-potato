#!/usr/bin/env bash

set -eux
time go build .
ls -l super-potato
