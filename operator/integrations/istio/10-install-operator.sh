#!/usr/bin/env bash

set -e
set -x

# install operator
export VM_OPERATOR_VERSION=$(basename $(curl -fs -o /dev/null -w %{redirect_url} \
  https://github.com/VictoriaMetrics/operator/releases/latest));

# TODO: rm when the latest is not v0.58.0 (it has buggy vmagent).
VM_OPERATOR_VERSION=v0.57.0

wget -O operator-and-crds.yaml \
  "https://github.com/VictoriaMetrics/operator/releases/download/$VM_OPERATOR_VERSION/install-no-webhook.yaml";
kubectl apply -f operator-and-crds.yaml;