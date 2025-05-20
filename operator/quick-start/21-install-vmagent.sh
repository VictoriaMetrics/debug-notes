#!/usr/bin/env bash

set -e
set -x

kubectl -n vm apply -f vmagent-demo.yaml;
kubectl -n vm wait --for=jsonpath='{.status.updateStatus}'=operational vmagent/demo;
kubectl -n vm rollout status deployment vmagent-demo  --watch=true;
