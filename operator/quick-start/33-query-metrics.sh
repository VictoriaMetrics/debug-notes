#!/usr/bin/env bash

set -e
set -x

curl -i --url http://127.0.0.1:8428/api/v1/query --url-query 'query=demo_counter_total{job="demo-app",namespace="default"}'