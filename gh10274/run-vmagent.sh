#!/usr/bin/env bash

set -e

TMP_DATA_PATH=$(mktemp -d)
echo "Using tmp data path: $TMP_DATA_PATH"

../../VictoriaMetrics/bin/vmagent \
  -promscrape.config=scrape.yaml \
  -remoteWrite.url=http://localhost:8428/api/v1/write \
  -remoteWrite.tmpDataPath="$TMP_DATA_PATH" \
  -httpListenAddr=:0 \
  -remoteWrite.disableOnDiskQueue=true
