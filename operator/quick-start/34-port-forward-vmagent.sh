#!/usr/bin/env bash

set -e
set -x

VMAGENT_POD_NAME=$(kubectl get pod -n vm -l "app.kubernetes.io/name=vmagent" -o jsonpath="{.items[0].metadata.name}")
kubectl port-forward -n vm $VMAGENT_POD_NAME 8439:8429
