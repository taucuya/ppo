#!/bin/bash

CORES=$1

python3 plotter_latency.py \
    processed/http_latency_cores_$CORES.csv \
    img/latency_cores_$CORES.png

python3 plotter_metrics.py \
    processed/http_metrics_cores_$CORES.csv \
    img/metrics_cores_$CORES.png  

python3 percentiles_hist.py \
    ./processed/rps_560_n_1_cores_2.0.csv