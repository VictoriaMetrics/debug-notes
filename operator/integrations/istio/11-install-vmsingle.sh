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
