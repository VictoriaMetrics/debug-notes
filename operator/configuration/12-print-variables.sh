#!/usr/bin/env bash

set -e
set -x

kubectl get deployment -n vm vm-operator \
    -o jsonpath='{range .spec.template.spec.containers[?(@.name=="manager")].env[*]}{.name}{"\n"}{end}';
