#!/usr/bin/env bash

set -e
set -x

cat <<'EOF' > demo-app-rule.yaml
apiVersion: operator.victoriametrics.com/v1beta1
kind: VMRule
metadata:
  name: demo
  namespace: default
spec:
  groups:
    - name: demo-app
      rules:
        - alert: DemoAlertFiring
          expr: 'sum(demo_alert_firing{job="demo-app",namespace="default"}) by (job,pod,namespace) > 0'
          for: 30s
          labels:
            job: '{{ $labels.job }}'
            pod: '{{ $labels.pod }}'
          annotations:
            description: 'demo-app pod {{ $labels.pod }} is firing demo alert'
EOF


kubectl apply -f demo-app-rule.yaml;
kubectl wait --for=jsonpath='{.status.updateStatus}'=operational vmrule/demo;