#!/usr/bin/env bash

set -e
set -x

kubectl exec -n default  demo-app -- curl -i http://127.0.0.1:8080/metrics
