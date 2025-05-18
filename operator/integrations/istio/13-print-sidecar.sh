#!/usr/bin/env bash

set -e
set -x

kubectl get pods --all-namespaces  -o jsonpath='{range .items[*]}{.metadata.name}{"\n"}{range .spec.containers[*]}- {.name}{"\n"}{end}{"\n"}{end}'