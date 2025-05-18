#!/usr/bin/env bash

set -e
set -x

cat <<'EOF' > vmsingle-demo.yaml
apiVersion: operator.victoriametrics.com/v1beta1
kind: VMSingle
metadata:
  name: demo
  namespace: vm
EOF

kubectl -n vm apply -f vmsingle-demo.yaml;
kubectl -n vm wait --for=jsonpath='{.status.updateStatus}'=operational vmsingle/demo;
kubectl -n vm rollout status deployment vmsingle-demo  --watch=true;
