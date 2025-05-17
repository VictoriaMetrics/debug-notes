#!/usr/bin/env bash

set -e
set -x

# install operator
export VM_OPERATOR_VERSION=$(basename $(curl -fs -o /dev/null -w %{redirect_url} \
  https://github.com/VictoriaMetrics/operator/releases/latest))
echo "VM_OPERATOR_VERSION=$VM_OPERATOR_VERSION"
wget -O operator-and-crds.yaml \
  "https://github.com/VictoriaMetrics/operator/releases/download/$VM_OPERATOR_VERSION/install-no-webhook.yaml"
kubectl apply -f operator-and-crds.yaml

VM_OPERATOR_VERSION=v0.57.0

# install vmsingle
cat <<'EOF' > vmsingle-demo.yaml
apiVersion: operator.victoriametrics.com/v1beta1
kind: VMSingle
metadata:
  name: demo
  namespace: vm
EOF
kubectl -n vm apply -f vmsingle-demo.yaml

# install vmagent
cat <<'EOF' > vmagent-demo.yaml
apiVersion: operator.victoriametrics.com/v1beta1
kind: VMAgent
metadata:
  name: demo
  namespace: vm
spec:
  selectAllByDefault: true
  remoteWrite:
    - url: "http://vmsingle-demo.vm.svc:8429/api/v1/write"
EOF
kubectl -n vm apply -f vmagent-demo.yaml

