#!/usr/bin/env bash

set -e
set -x

./02-install-operator.sh
./11-install-vmsingle.sh
./21-install-vmagent.sh
./30-create-demo-app-manifest.sh
./32-create-demo-app-service-scrape.sh
./40-create-vmagentmanager.sh
./41-create-vmalert.sh
./43-create-demo-app-rule.sh
./50-create-vmauth-manifest.sh
./52-create-vmuser-manifest.sh