#!/usr/bin/env bash

set -e
set -x

cat <<EOF >vmagent-demo.yaml
apiVersion: operator.victoriametrics.com/v1beta1
kind: VMAgent
metadata:
  name: demo
  namespace: vm
spec:
  selectAllByDefault: true
  remoteWrite:
    - url: "http://vmsingle-demo.vm.svc:8480/prometheus/api/v1/write"
EOF