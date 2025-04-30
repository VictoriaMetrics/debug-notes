#!/usr/bin/env bash

set -e
set -x

cat > vmsingle-demo.yaml <<'EOF'
apiVersion: operator.victoriametrics.com/v1beta1
kind: VMSingle
metadata:
  name: demo
  namespace: vm
spec:
  # make it optional, remove from this file
  retentionPeriod: "4d"
EOF