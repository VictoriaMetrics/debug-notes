#!/usr/bin/env bash

set -e
set -x

kubectl apply -f operator-and-crds.yaml;
kubectl -n vm rollout status deployment vm-operator --watch=true;

