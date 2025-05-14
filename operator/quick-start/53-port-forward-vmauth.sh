#!/usr/bin/env bash

set -e
set -x

kubectl -n vm port-forward svc/vmauth-demo 8427:8427