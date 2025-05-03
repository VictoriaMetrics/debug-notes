#!/usr/bin/env bash

set -e
set -x

kubectl -n vm apply -f vmagent-demo.yaml
