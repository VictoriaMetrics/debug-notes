#!/usr/bin/env bash

set -e
set -x

OPERATOR_POD_NAME=$(kubectl get pod -l "app.kubernetes.io/name=victoria-metrics-operator"  -n vm -o jsonpath="{.items[0].metadata.name}");
kubectl exec -n vm "$OPERATOR_POD_NAME" -- /app --printDefaults 2>&1 | grep VMSINGLEDEFAULT_RESOURCE;