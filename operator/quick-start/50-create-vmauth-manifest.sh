#!/usr/bin/env bash

set -e
set -x

cat <<'EOF' > vmauth-demo.yaml
apiVersion: operator.victoriametrics.com/v1beta1
kind: VMAuth
metadata:
  name: demo
spec:
  selectAllByDefault: true
  userNamespaceSelector: {}
  userSelector: {}
  ingress:
    class_name: nginx # <-- change this to your ingress-controller
    host: vm-demo.k8s.orb.local # <-- change this to your domain
EOF

kubectl -n vm apply -f vmauth-demo.yaml