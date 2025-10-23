#!/usr/bin/env bash

./vmstorage-rpc \
    -storageDataPath=vmstorage-rpc-data \
    -vminsertAddr=:8412 \
    -vmselectAddr=:8402 \
    -httpListenAddr=:8484
