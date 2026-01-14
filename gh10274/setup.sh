#!/usr/bin/env bash

set -e

(cd ../../VictoriaMetrics && make victoria-metrics)
(cd ../../VictoriaMetrics && make vmagent)