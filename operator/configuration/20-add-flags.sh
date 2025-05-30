#!/usr/bin/env bash

set -e
set -x

mkdir -p add-operator-flag;

cat <<'EOF' > add-operator-flag/patch.yaml
- op: add
  path: /spec/template/spec/containers/0/args/-
  value: '-zap-log-level=debug'
EOF

cat <<'EOF' > add-operator-flag/kustomization.yaml
resources:
  - ../operator-and-crds.yaml

patches:
  - path: patch.yaml
    target:
      kind: Deployment
      name: vm-operator
EOF

kustomize build add-operator-flag -o operator-and-crds.yaml --load-restrictor=LoadRestrictionsNone;
cat operator-and-crds.yaml | grep "zap-log-level";
