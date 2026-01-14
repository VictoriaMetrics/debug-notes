#!/usr/bin/env bash

set -e

TMP_DATA_PATH=$(mktemp -d)
echo "Using tmp data path: $TMP_DATA_PATH"

../../VictoriaMetrics/bin/victoria-metrics \
  -search.latencyOffset=1s \
  -dedup.minScrapeInterval=1s \
  "-storageDataPath=${TMP_DATA_PATH}"

