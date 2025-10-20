#!/usr/bin/env bash

MAX_CONCURRENT_INSERTS=${1:-500}
BACKPRESSURE=${2:-true}
CAP_REROUTING=${3:-false}
DISABLE_REROUTING=$([ "$CAP_REROUTING" = "true" ] && echo "false" || echo "true")

cp vminsert-linux-arm64 vminsert-linux-arm64-1

./vminsert-linux-arm64-1 \
  -httpListenAddr=:8381 \
  -storageNode=192.168.50.102:8411,192.168.50.102:8412 \
  -maxConcurrentInserts="${MAX_CONCURRENT_INSERTS}" \
  -capacityRerouting="${CAP_REROUTING}" \
  -disableRerouting="${DISABLE_REROUTING}" \
  -backpressure="${BACKPRESSURE}" 2>&1 | tee vminsert-1.log