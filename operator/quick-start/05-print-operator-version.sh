#!/usr/bin/env bash

set -e
set -x

OPERATOR_POD_NAME=$(kubectl get pod -l "control-plane=vm-operator"  -n vm -o jsonpath="{.items[0].metadata.name}")
kubectl exec -n vm "$OPERATOR_POD_NAME" -- /app --printDefaults