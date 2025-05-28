#!/usr/bin/env bash

set -e
set -x

mkdir -p demo-app-scrape-with-tls;

cat <<'EOF' > demo-app-scrape-with-tls/patch.yaml
- op: replace
  path: /spec/endpoints/0/scheme
  value: https
- op: add
  path: /spec/endpoints/0/tlsConfig
  value:
    caFile: /etc/istio-certs/root-cert.pem
    certFile: /etc/istio-certs/cert-chain.pem
    insecureSkipVerify: true
    keyFile: /etc/istio-certs/key.pem

EOF

cat <<'EOF' > demo-app-scrape-with-tls/kustomization.yaml
resources:
  - ../demo-app-scrape.yaml

patches:
  - path: patch.yaml
    target:
      kind: VMServiceScrape
      name: demo-app-service-scrape
EOF

kustomize build demo-app-scrape-with-tls -o demo-app-scrape.yaml --load-restrictor=LoadRestrictionsNone;
cat demo-app-scrape.yaml;