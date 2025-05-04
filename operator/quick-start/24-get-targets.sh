#!/usr/bin/env bash

set -e
set -x

VMAGENT_POD_NAME=$(kubectl get pod -n vm -l "app.kubernetes.io/name=vmagent" -o jsonpath="{.items[0].metadata.name}")
kubectl exec -n vm $VMAGENT_POD_NAME -c vmagent  -- wget -qO -  http://127.0.0.1:8429/api/v1/targets |
  jq -r '.data.activeTargets[].discoveredLabels.__meta_kubernetes_endpoint_address_target_name'