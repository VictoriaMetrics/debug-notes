#!/usr/bin/env bash

set -e
set -x

cat <<'EOF' > demo-app.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: demo-app
  namespace: default
  labels:
    app.kubernetes.io/name: demo-app
spec:
  replicas: 1
  selector:
    matchLabels:
      app.kubernetes.io/name: demo-app
  template:
    metadata:
      labels:
        app.kubernetes.io/name: demo-app
    spec:
      containers:
        - name: main
          image: docker.io/victoriametrics/demo-app:1.2
---
apiVersion: v1
kind: Service
metadata:
  name: demo-app
  namespace: default
  labels:
    app.kubernetes.io/name: demo-app
spec:
  selector:
    app.kubernetes.io/name: demo-app
  ports:
    - port: 8080
      name: http
EOF

kubectl -n default apply -f demo-app.yaml