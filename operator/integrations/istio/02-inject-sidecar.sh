#!/usr/bin/env bash

set -e
set -x

kubectl label namespace default istio-injection=enabled;

kubectl create namespace vm || true;
kubectl label namespace vm istio-injection=enabled;