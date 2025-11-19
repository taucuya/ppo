#!/bin/bash
echo "Running increase"
./http_rps_increase.sh
echo "Running constant"
./http_constant.sh
echo "Running recovery"
./http_recovery.sh
echo "Done"