#!/usr/bin/env bash

set -e
set -x

mkdir -p add-operator-envs;

cat <<'EOF' > add-operator-envs/patch.yaml
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

cat <<'EOF' > add-operator-envs/kustomization.yaml
resources:
  - ../operator-and-crds.yaml

patches:
  - path: patch.yaml
    target:
      kind: Deployment
      name: vm-operator
EOF

kustomize build add-operator-envs -o operator-and-crds.yaml --load-restrictor=LoadRestrictionsNone;
cat operator-and-crds.yaml | grep -E -A 1 "VM_VMSINGLEDEFAULT_RESOURCE_LIMIT_MEM|VM_VMSINGLEDEFAULT_RESOURCE_LIMIT_CPU";
