#!/usr/bin/env bash

set -e
set -x

cat <<'EOF' > demo-app-scrape.yaml
apiVersion: operator.victoriametrics.com/v1beta1
kind: VMServiceScrape
metadata:
  name: demo-app-service-scrape
spec:
  selector:
    matchLabels:
      app.kubernetes.io/name: demo-app
  endpoints:
  - port: http
    path: /metrics
EOF

kubectl apply -f demo-app-scrape.yaml;
kubectl wait --for=jsonpath='{.status.updateStatus}'=operational vmservicescrape/demo-app-service-scrape;