#!/usr/bin/env bash

set -e
set -x

DEMO_APP_POD=$(kubectl get pod -n default -l "app.kubernetes.io/name=demo-app" -o jsonpath="{.items[0].metadata.name}");
kubectl exec -n default  ${DEMO_APP_POD} -- curl -s --url http://127.0.0.1:8080/alerting/receivedWebhooks;