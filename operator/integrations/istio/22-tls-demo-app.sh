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
    scheme: https
    tlsConfig:
      insecureSkipVerify: true
      caFile: /etc/istio-certs/root-cert.pem
      certFile: /etc/istio-certs/cert-chain.pem
      keyFile: /etc/istio-certs/key.pem
EOF

kubectl apply -f demo-app-scrape.yaml
