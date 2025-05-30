#!/usr/bin/env bash

set -e
set -x

mkdir -p vmagent-with-istio-certs;

cat <<'EOF' > vmagent-with-istio-certs/patch.yaml
apiVersion: operator.victoriametrics.com/v1beta1
kind: VMAgent
metadata:
  name: demo
  namespace: vm
spec:
  podMetadata:
    annotations:
      proxy.istio.io/config: |  # configure an env variable `OUTPUT_CERTS` to write certificates to the given folder
        proxyMetadata:
          OUTPUT_CERTS: /etc/istio-certs
      sidecar.istio.io/userVolumeMount: '[{"name": "istio-certs", "mountPath": "/etc/istio-certs"}]'
  volumeMounts:
    - mountPath: /etc/istio-certs/
      name: istio-certs
      readOnly: true
  volumes:
    - emptyDir:
        medium: Memory
      name: istio-certs
EOF

cat <<'EOF' > vmagent-with-istio-certs/kustomization.yaml
resources:
  - ../vmagent-demo.yaml

patches:
  - path: patch.yaml
    target:
      kind: VMAgent
      name: demo
EOF

kustomize build vmagent-with-istio-certs -o vmagent-demo.yaml --load-restrictor=LoadRestrictionsNone;
cat vmagent-demo.yaml;
