#!/usr/bin/env bash

set -e
set -x

kubectl get pods -n vm -l "control-plane=vm-operator" -o jsonpath='{range .items[*]}{range .spec.containers[?(@.name=="manager")]}{.image}{"\n"}{end}{end}'