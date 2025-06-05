#!/usr/bin/env bash

set -e
set -x

cat <<'EOF' > operator-scrape.yaml
apiVersion: operator.victoriametrics.com/v1beta1
kind: VMServiceScrape
metadata:
  name: operator-service-scrape
  # You might need to change the namespace below
  namespace: vm
spec:
  selector:
    matchLabels:
      # You might need to change the labels below
      app.kubernetes.io/instance: default
      app.kubernetes.io/name: victoria-metrics-operator
  endpoints:
    - port: http
EOF

kubectl apply -f operator-scrape.yaml;
kubectl wait -n vm --for=jsonpath='{.status.updateStatus}'=operational vmservicescrape/operator-service-scrape;
