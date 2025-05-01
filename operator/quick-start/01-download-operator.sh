#!/usr/bin/env bash

set -e
set -x

export VM_OPERATOR_VERSION=`basename $(curl -fs -o/dev/null -w %{redirect_url} https://github.com/VictoriaMetrics/operator/releases/latest)`
wget -O operator-and-crds.yaml https://github.com/VictoriaMetrics/operator/releases/download/$VM_OPERATOR_VERSION/install-no-webhook.yaml