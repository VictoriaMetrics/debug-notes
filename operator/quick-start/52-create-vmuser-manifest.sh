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
        - "/api/v1/query"
    # vmalert
    - crd:
        kind: VMAlert
        name: demo
        namespace: vm
      paths:
        - "/vmalert.*"
        - "/api/v1/rules"
EOF

kubectl -n vm apply -f vmuser-demo.yaml
