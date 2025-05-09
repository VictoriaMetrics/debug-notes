#!/usr/bin/env bash

set -e
set -x

VMALERT_POD_NAME=$(kubectl get pod -n vm -l "app.kubernetes.io/name=vmalert" -o jsonpath="{.items[0].metadata.name}")
kubectl port-forward -n vm $VMALERT_POD_NAME 8080:8080
