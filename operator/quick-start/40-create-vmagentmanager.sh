#!/usr/bin/env bash

set -e
set -x

cat <<EOF > vmalertmanager-demo.yaml
apiVersion: operator.victoriametrics.com/v1beta1
kind: VMAlertmanager
metadata:
  name: demo
  namespace: vm
spec:
  configRawYaml: |
    route:
      receiver: 'demo-app'
    receivers:
    - name: 'demo-app'
      webhook_configs:
      - url: 'http://demo-app.default.svc:8080/alerting/webhook'
EOF

kubectl apply -f vmalertmanager-demo.yaml;