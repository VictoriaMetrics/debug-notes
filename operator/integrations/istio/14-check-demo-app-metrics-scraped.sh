#!/usr/bin/env bash

set -e
set -x

DEMO_APP_POD=$(kubectl get pod -n default -l "app.kubernetes.io/name=demo-app" -o jsonpath="{.items[0].metadata.name}");
DEMO_APP_POD_IP=$(kubectl get pod -n default ${DEMO_APP_POD} -o jsonpath='{.status.podIP}');
VMAGENT_POD_NAME=$(kubectl get pod -n vm -l "app.kubernetes.io/name=vmagent" -o jsonpath="{.items[0].metadata.name}");

kubectl -n vm logs -f "${VMAGENT_POD_NAME}" -c istio-proxy 2>&1 # | grep "${DEMO_APP_POD_IP}";