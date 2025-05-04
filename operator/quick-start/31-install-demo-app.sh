#!/usr/bin/env bash

set -e
set -x

kubectl -n default apply -f demo-app.yaml
