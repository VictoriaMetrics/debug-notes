#!/usr/bin/env bash

set -e
set -x

cat <<'EOF' > vmagent-demo.yaml
apiVersion: operator.victoriametrics.com/v1beta1
kind: VMAgent
metadata:
  name: demo
  namespace: vm
spec:
  selectAllByDefault: true
  remoteWrite:
    - url: "http://vmsingle-demo.vm.svc:8429/api/v1/write"

  extraArgs:
    remoteWrite.tlsInsecureSkipVerify: "true"
  podMetadata:
    annotations:
      traffic.sidecar.istio.io/includeInboundPorts: ""   # do not intercept any inbound ports
      traffic.sidecar.istio.io/includeOutboundIPRanges: ""  # do not intercept any outbound traffic
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
kubectl -n vm apply -f vmagent-demo.yaml
