#!/usr/bin/env bash

set -e
set -x

../quick-start/10-create-vmsingle-manifest.sh
../quick-start/11-install-vmsingle.sh
../quick-start/20-create-vmagent-manifest.sh
../quick-start/21-install-vmagent.sh
../quick-start/30-create-demo-app-manifest.sh
../quick-start/32-create-demo-app-service-scrape.sh