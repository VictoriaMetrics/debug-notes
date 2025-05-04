#!/usr/bin/env bash

set -e
set -x

curl -i --url http://127.0.0.1:8429/api/v1/query?query=demo_counter_total