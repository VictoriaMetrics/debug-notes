#!/usr/bin/env bash

set -e
set -x

# pod is running
kubectl get pods -n vm -l "control-plane=vm-operator"

# vm resource installed
kubectl api-resources --api-group=operator.victoriametrics.com

# operator version
kubectl get pods -n vm -l "control-plane=vm-operator" -o jsonpath='{range .items[*]}{range .spec.containers[?(@.name=="manager")]}{.image}{"\n"}{end}{end}'

# env variables
OPERATOR_POD_NAME=$(kubectl get pod -l "control-plane=vm-operator"  -n vm -o jsonpath="{.items[0].metadata.name}")
kubectl exec -n vm "$OPERATOR_POD_NAME" -- /app --printDefaults