#!/usr/bin/env bash

set -e
set -x

mkdir -p demo-app-scrape-with-tls;

kubectl apply -f demo-app-scrape.yaml;
kubectl wait --for=jsonpath='{.status.updateStatus}'=operational vmservicescrape/demo-app-service-scrape;
