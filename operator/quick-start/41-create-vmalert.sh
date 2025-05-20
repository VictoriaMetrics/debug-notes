#!/usr/bin/env bash

set -e
set -x

cat <<'EOF' > vmalert-demo.yaml
apiVersion: operator.victoriametrics.com/v1beta1
kind: VMAlert
metadata:
  name: demo
  namespace: vm
spec:
  # Metrics source (VMCluster/VMSingle)
  datasource:
    url: "http://vmsingle-demo.vm.svc:8429"

  # Where alert state and recording rules are stored
  remoteWrite:
    url: "http://vmsingle-demo.vm.svc:8429"

  # Where the previous alert state is loaded from. Optional
  remoteRead:
    url: "http://vmsingle-demo.vm.svc:8429"

  # Alertmanager URL for sending alerts
  notifier:
    url: "http://vmalertmanager-demo.vm.svc:9093"

  # How often the rules are evaluated
  evaluationInterval: "10s"

  # Watch VMRule resources in all namespaces
  selectAllByDefault: true
EOF

kubectl apply -f vmalert-demo.yaml;
kubectl -n vm wait --for=jsonpath='{.status.updateStatus}'=operational vmalert/demo;
kubectl -n vm rollout status deployment vmalert-demo  --watch=true;