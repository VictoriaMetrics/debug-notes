#!/usr/bin/env bash

set -e
set -x

cat <<'EOF' > operator-patch.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: vm-operator
  namespace: vm
spec:
  template:
    spec:
      containers:
      - name: manager
        env:
        - name: VM_VMSINGLEDEFAULT_RESOURCE_LIMIT_MEM
          value: "3000Mi"
        - name: VM_VMSINGLEDEFAULT_RESOURCE_LIMIT_CPU
          value: "2400m"
EOF

cat <<'EOF' > kustomization.yaml
resources:
  - operator-and-crds.yaml

patches:
  - path: operator-patch.yaml
    target:
      kind: Deployment
      name: vm-operator
EOF

kustomize build -o operator-and-crds.yaml;