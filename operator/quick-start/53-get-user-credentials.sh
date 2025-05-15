#!/usr/bin/env bash

set -e
set -x

export DEMO_USERNAME="$(kubectl get secret -n vm vmuser-demo -o jsonpath="{.data.username}" | base64 --decode)";
export DEMO_PASSWORD="$(kubectl get secret -n vm vmuser-demo -o jsonpath="{.data.password}" | base64 --decode)";
echo "Username: $DEMO_USERNAME; Password: $DEMO_PASSWORD";