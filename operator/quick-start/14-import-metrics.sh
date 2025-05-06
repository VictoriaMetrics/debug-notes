#!/usr/bin/env bash

set -e
set -x

curl -i -X POST \
  --url http://127.0.0.1:8428/api/v1/import/prometheus \
  --header 'Content-Type: text/plain' \
  --data 'a_metric{foo="fooVal"} 1
'
