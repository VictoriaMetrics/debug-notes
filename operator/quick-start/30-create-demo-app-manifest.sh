#!/usr/bin/env bash

set -e
set -x

cat <<EOF > demo-app.yaml
apiVersion: v1
kind: Pod
metadata:
  name: demo-app
  namespace: default
  labels:
    app.kubernetes.io/name: demo-app
spec:
  containers:
    - name: main
      image: docker.io/makasim/demo-app
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
    - port: 9100
      name: metrics
EOF