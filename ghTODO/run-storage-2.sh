#!/usr/bin/env bash

GOMAXPROCS=2 ./../../VictoriaMetrics/bin/vmstorage \
    -storageDataPath=vmstorage-2 \
    -vminsertAddr=:8412 \
    -vmselectAddr=:8402 \
    -httpListenAddr=:8484