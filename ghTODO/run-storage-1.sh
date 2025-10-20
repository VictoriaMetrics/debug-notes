#!/usr/bin/env bash

GOMAXPROCS=2 ./bin/vmstorage \
    -storageDataPath=vmstorage-1 \
    -vminsertAddr=:8411 \
    -vmselectAddr=:8401 \
    -httpListenAddr=:8483
