#!/usr/bin/env bash

set -e
set -x

curl -L https://istio.io/downloadIstio | ISTIO_VERSION=1.26.0 sh -;
./istio-1.26.0/bin/istioctl install -f ./istio-1.26.0/samples/bookinfo/demo-profile-no-gateways.yaml -y;