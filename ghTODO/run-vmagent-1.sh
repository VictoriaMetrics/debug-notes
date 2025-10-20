#!/usr/bin/env bash

./../../VictoriaMetrics/bin/vmagent \
    -remoteWrite.queues=2000 \
    -remoteWrite.url=http://127.0.0.1:8427/insert/0/prometheus/api/v1/write \
    -maxConcurrentInserts=500 \
    -remoteWrite.tmpDataPath=vmagentdata \
    -remoteWrite.maxDiskUsagePerURL=2GB \
    -remoteWrite.basicAuth.username=foo \
    -remoteWrite.basicAuth.password=bar


#      -remoteWrite.url=http://vminsert-1:8480/insert/0/prometheus/api/v1/write"