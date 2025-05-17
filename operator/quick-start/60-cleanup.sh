#!/usr/bin/env bash

set -e
set -x

vm_resource_names=$(kubectl api-resources \
  --api-group=operator.victoriametrics.com \
  --namespaced=true \
  --output=name)

for vm_resource_name in $vm_resource_names; do
  echo "Deleting all resources of type $vm_resource_name"
  kubectl delete "$vm_resource_name" --all --all-namespaces
done

kubectl delete --ignore-not-found=true -f demo-app.yaml
kubectl delete --ignore-not-found=true -f operator-and-crds.yaml
