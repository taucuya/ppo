#!/bin/bash
OUT_DIR=constant_rps_out
IMG_DIR=img
RPS_VALUES=(560)
DURATION="300s"  
CONTAINER_NAME=ppo-api
NCORES=2.0
N=1

rm -rf $OUT_DIR
mkdir -p $OUT_DIR
mkdir -p $OUT_DIR/http
mkdir -p $IMG_DIR

echo "Benchmark: RPS constant (HTTP)"
echo "RPS levels: ${RPS_VALUES[@]}"

echo "Begin HTTP test"
echo "Starting benchmark with $N runs per RPS level"

for ((RUN=1; RUN<=N; RUN++)); do
    echo "Run $RUN/$N"
    
    METRICS_FILE="$OUT_DIR/http/metrics_n_${RUN}_cores_${NCORES}.csv"
    ./metric_scraper.sh $CONTAINER_NAME $METRICS_FILE 1 & 
    MONITOR_PID=$!
    echo "Scraper PID: $MONITOR_PID"

    for RPS in "${RPS_VALUES[@]}"; do        
        OUTPUT_FILE="$OUT_DIR/http/rps_${RPS}_n_${RUN}_cores_${NCORES}.csv"
        echo "Testing RPS: $RPS"
        
        CONCURRENT=$(( RPS / 50 ))  
        if [ $CONCURRENT -lt 10 ]; then
            CONCURRENT=10  
        fi
        TOTAL_REQUESTS=$(( RPS * 10 )) 
        
        echo "  Concurrent: $CONCURRENT, Total requests: $TOTAL_REQUESTS"
        
        hey -z $DURATION -q $RPS -c $CONCURRENT -n $TOTAL_REQUESTS -m POST \
            -T "application/json" \
            -d '{"name": "test", "date_of_birth": "2010-05-05", "email": "a@mail.com", "password": "111", "phone": "89164543280", "address": "a"}' \
            -o csv http://localhost:8080/api/v1/auth/signup > $OUTPUT_FILE
            
        python3 hey_to_csv.py $OUTPUT_FILE $RPS
        
        TOTAL_PROCESSED=$(wc -l < $OUTPUT_FILE)
        echo "  Processed requests: $((TOTAL_PROCESSED - 1))"  
    done
    
    kill $MONITOR_PID
    wait $MONITOR_PID 2>/dev/null
    pkill -f "metric_scraper.sh" 2>/dev/null
    echo "Run $RUN completed"
done