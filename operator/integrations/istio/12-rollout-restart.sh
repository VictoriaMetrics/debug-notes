#!/usr/bin/env bash

set -e
set -x

kubectl -n vm rollout restart deployment vm-operator;
kubectl -n vm rollout status deployment vm-operator --watch=true;

kubectl -n vm rollout restart deployment vmsingle-demo;
kubectl -n vm rollout status deployment vmsingle-demo  --watch=true;

kubectl -n vm rollout restart deployment vmagent-demo;
kubectl -n vm rollout status deployment vmagent-demo  --watch=true;

kubectl -n default rollout restart deployment demo-app;
kubectl -n default rollout status deployment demo-app  --watch=true;