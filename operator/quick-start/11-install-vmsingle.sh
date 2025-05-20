#!/usr/bin/env bash

set -e
set -x

kubectl -n vm apply -f vmsingle-demo.yaml;
kubectl -n vm wait --for=jsonpath='{.status.updateStatus}'=operational vmsingle/demo;
kubectl -n vm rollout status deployment vmsingle-demo  --watch=true;

