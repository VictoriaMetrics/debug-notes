#!/usr/bin/env bash

set -e
set -x

cat <<'EOF' > global-peer-authentication.yaml
apiVersion: security.istio.io/v1
kind: PeerAuthentication
metadata:
  name: default
  namespace: istio-system
spec:
  mtls:
    mode: STRICT
EOF

cat <<'EOF' > vm-ns-peer-authentication.yaml
apiVersion: security.istio.io/v1
kind: PeerAuthentication
metadata:
  name: default
  namespace: vm
spec:
  mtls:
    mode: PERMISSIVE
EOF

kubectl -n istio-system apply -f global-peer-authentication.yaml;
kubectl -n vm apply -f default-ns-peer-authentication.yaml
