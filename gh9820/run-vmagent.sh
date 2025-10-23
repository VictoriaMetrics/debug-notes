#!/usr/bin/env bash

./vmagent -remoteWrite.url=http://127.0.0.1:8427/insert/0/prometheus/api/v1/write \
  -remoteWrite.tmpDataPath=vmagentdata \
  -remoteWrite.maxDiskUsagePerURL=2GB \
  -promscrape.config=prometheus.yml