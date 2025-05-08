#!/usr/bin/env bash

set -e
set -x

cat <<EOF > vmalert-demo.yaml
apiVersion: operator.victoriametrics.com/v1beta1
kind: VMAlert
metadata:
  name: demo
  namespace: vm
spec:
  # Prometheus compatible storage to query metrics from.
  datasource:
    url: "http://vmsingle-demo.vm.svc:8489"
  # Prometheus remote storage to send alert state to.
  remoteWrite:
    url: "http://vmsingle-demo.vm.svc:8489"
  # Prometheus remote storage to restore alert state from.
  remoteRead:
    url: "http://vmsingle-demo.vm.svc:8489"
  notifier:
    url: "http://vmalertmanager-demo.vm.svc:9093"
  # alerts to be evaluated at 30s interval by default
  evaluationInterval: "30s"
  # Watch VMRule CRDs in all namespaces across the cluster.
  selectAllByDefault: true
EOF

kubectl apply -f vmalert-demo.yaml