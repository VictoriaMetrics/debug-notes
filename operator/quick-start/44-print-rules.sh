#!/usr/bin/env bash

set -e
set -x

VMALERT_POD_NAME=$(kubectl get pod -n vm -l "app.kubernetes.io/name=vmalert" -o jsonpath="{.items[0].metadata.name}");
kubectl exec -n vm $VMALERT_POD_NAME -c vmalert  -- wget -qO -  http://127.0.0.1:8080/api/v1/rules |
  jq -r '.data.groups[].rules[].name';