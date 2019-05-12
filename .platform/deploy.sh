#!/usr/bin/env bash

time ./super-potato dump -f environ > tmp/.deploy_env
time ./super-potato dump -f shell > tmp/deploy_env.sh
