#!/usr/bin/env bash

set -e
set -x

kubectl -n vm apply -f vmsingle-demo.yaml
