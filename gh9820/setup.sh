#!/usr/bin/env bash

set -e

(rm -rf ../../VictoriaMetrics/bin/*)
(cd ../../VictoriaMetrics && git checkout cluster)
(cd ../../VictoriaMetrics && git pull opensource cluster)
(cd ../../VictoriaMetrics && make vmstorage vminsert vmauth vmselect vmagent && \
  cp ./bin/vmstorage ../debug-notes/gh9820/vmstorage-legacy && \
  cp ./bin/vminsert ../debug-notes/gh9820/vminsert-legacy && \
  cp ./bin/vmselect ../debug-notes/gh9820/vmselect && \
  cp ./bin/vmagent ../debug-notes/gh9820/vmagent && \
  cp ./bin/vmauth ../debug-notes/gh9820/vmauth)

(rm -rf ../../VictoriaMetrics/bin/*)
(cd ../../VictoriaMetrics && git checkout vminsert-rpc)
(cd ../../VictoriaMetrics && git pull opensource vminsert-rpc)
(cd ../../VictoriaMetrics && make vmstorage vminsert && cp ./bin/vmstorage ../debug-notes/gh9820/vmstorage-rpc && cp ./bin/vminsert ../debug-notes/gh9820/vminsert-rpc)