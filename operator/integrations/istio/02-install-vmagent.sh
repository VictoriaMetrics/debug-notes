#!/usr/bin/env bash

set -e
set -x

cat <<'EOF' > vmagent-demo.yaml
apiVersion: operator.victoriametrics.com/v1beta1
kind: VMAgent
metadata:
  name: demo
  namespace: vm
spec:
  selectAllByDefault: true
  remoteWrite:
    - url: "http://vmsingle-demo.vm.svc:8429/api/v1/write"
EOF

kubectl -n vm apply -f vmagent-demo.yaml;
kubectl -n vm wait --for=jsonpath='{.status.updateStatus}'=operational vmagent/demo;
kubectl -n vm rollout status deployment vmagent-demo  --watch=true;
