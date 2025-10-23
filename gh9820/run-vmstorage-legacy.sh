#!/usr/bin/env bash

./vmstorage-legacy \
    -storageDataPath=vmstorage-legacy-data \
    -vminsertAddr=:8411 \
    -vmselectAddr=:8401 \
    -httpListenAddr=:8483
