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
  ingress:
    class_name: 'nginx'                 # <-- Change this to match your Ingress controller (e.g., 'traefik')
    host: victoriametrics.mycompany.com # <-- Change this to the domain name you’ll use
EOF

kubectl -n vm apply -f vmauth-demo.yaml