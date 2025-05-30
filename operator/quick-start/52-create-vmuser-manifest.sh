#!/usr/bin/env bash

set -e
set -x

cat <<'EOF' > vmuser-demo.yaml
apiVersion: operator.victoriametrics.com/v1beta1
kind: VMUser
metadata:
  name: demo
spec:
  name: demo
  username: demo
  generatePassword: true
  targetRefs:
    # vmsingle
    - crd:
        kind: VMSingle
        name: demo
        namespace: vm
      paths:
        - "/vmui.*"
        - "/prometheus/.*"
    # vmalert
    - crd:
        kind: VMAlert
        name: demo
        namespace: vm
      paths:
        - "/vmalert.*"
        - "/api/v1/groups"
        - "/api/v1/alert"
        - "/api/v1/alerts"
EOF

kubectl -n vm apply -f vmuser-demo.yaml;
kubectl -n vm wait --for=jsonpath='{.status.updateStatus}'=operational vmuser/demo;
