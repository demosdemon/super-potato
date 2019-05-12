#!/usr/bin/env bash

time ./super-potato dump -f environ > tmp/.post_deploy_env
time ./super-potato dump -f shell > tmp/post_deploy_env.sh
