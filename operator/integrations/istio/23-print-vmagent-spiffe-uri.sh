#!/usr/bin/env bash

set -e
set -x

VMAGENT_POD_NAME=$(kubectl get pod -n vm -l "app.kubernetes.io/name=vmagent" -o jsonpath="{.items[0].metadata.name}");
kubectl -n vm exec -it $VMAGENT_POD_NAME -c istio-proxy -- openssl x509 -in /etc/istio-certs/cert-chain.pem -noout -text | grep -A 1 "Subject Alternative Name"