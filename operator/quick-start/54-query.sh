#!/usr/bin/env bash

set -e
set -x

export DEMO_USERNAME="$(kubectl get secret -n vm vmuser-demo -o jsonpath="{.data.username}" | base64 --decode)"
export DEMO_PASSWORD="$(kubectl get secret -n vm vmuser-demo -o jsonpath="{.data.password}" | base64 --decode)"
echo "Username: $DEMO_USERNAME; Password: $DEMO_PASSWORD"


curl -i -u "${DEMO_USERNAME}:${DEMO_PASSWORD}123" -H "Host: vm-demo.k8s.orb.local" \
  --url http://127.0.0.1:8427/api/v1/query --url-query query=a_metric