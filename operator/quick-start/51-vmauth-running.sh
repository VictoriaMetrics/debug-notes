#!/usr/bin/env bash

set -e
set -x

kubectl get pod -n vm -l "app.kubernetes.io/name=vmauth" -l "app.kubernetes.io/name=vmauth"

kubectl get ingress -n vm