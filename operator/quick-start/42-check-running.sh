#!/usr/bin/env bash

set -e
set -x

kubectl get pods -n vm -l "app.kubernetes.io/name=vmalertmanager"
kubectl get pods -n vm -l "app.kubernetes.io/name=vmalert"