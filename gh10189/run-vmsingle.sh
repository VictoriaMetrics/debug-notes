#!/usr/bin/env bash

./../../VictoriaMetrics/bin/victoria-metrics \
    -promscrape.config=scrape.yaml \
    -search.latencyOffset=1s
