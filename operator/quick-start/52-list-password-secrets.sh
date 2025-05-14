#!/usr/bin/env bash

set -e
set -x

kubectl get secret -n vm -l "app.kubernetes.io/instance=demo" -l "app.kubernetes.io/name=vmuser"