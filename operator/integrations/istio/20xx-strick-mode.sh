#!/usr/bin/env bash

set -e
set -x

cat <<'EOF' > default-ns-peer-authentication.yaml
apiVersion: security.istio.io/v1
kind: PeerAuthentication
metadata:
  name: demo
  namespace: default
spec:
  mtls:
    mode: STRICT
EOF
kubectl -n default apply -f default-ns-peer-authentication.yaml

#kubectl create namespace vm || true
#kubectl label namespace vm istio-injection=enabled
#
#cat <<'EOF' > vm-ns-peer-authentication.yaml
#apiVersion: security.istio.io/v1
#kind: PeerAuthentication
#metadata:
#  name: demo
#  namespace: vm
#spec:
#  mtls:
#    mode: PERMISSIVE
#EOF
#kubectl -n vm apply -f vm-ns-peer-authentication.yaml