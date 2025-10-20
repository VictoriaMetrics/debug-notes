#!/usr/bin/env bash

./../../VictoriaMetrics/bin/vmauth \
    -auth.config=auth-vm-cluster.yml \
    -maxConcurrentPerUserRequests=10000 \
    -maxConcurrentRequests=10000 \
    -responseTimeout=2m
