#!/usr/bin/env bash

set -eux
time go get ./...
time go build -v -a -ldflags '-s -w -extldflags "-static"' ./cmd/gen
# fail to build if gen writes a file
time ./gen --exit-code enums ./data/enums.yaml ./pkg/platformsh/enums_gen.go
time ./gen --exit-code variables ./data/variables.yaml ./pkg/platformsh/environment_gen.go
time ./gen --exit-code api ./data/variables.yaml ./cmd/serve/generated.go
time go build -v -a -ldflags '-s -w -extldflags "-static"' .
time ./super-potato dump -f environ > data/.build_env
time ./super-potato dump -f shell > data/build_env.sh
ls -l super-potato gen
