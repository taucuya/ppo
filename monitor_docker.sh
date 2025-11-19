#!/bin/sh

CONTAINER_NAME=$1 
OUTPUT_FILE=$2     
DURATION=${3:-15}

echo "Monitoring container: $CONTAINER_NAME for $DURATION seconds"

get_timestamp() {
    date +%s%3N
}

check_container() {
    docker ps --format "{{.Names}}" | grep -q "^${CONTAINER_NAME}$"
    return $?
}

echo "Waiting for container $CONTAINER_NAME..."
for i in $(seq 1 60); do
    if check_container; then
        echo "Container found!"
        break
    fi
    if [ $i -eq 60 ]; then
        echo "ERROR: Container $CONTAINER_NAME not found after 60 seconds"
        exit 1
    fi
    sleep 1
done

START_TIME=$(date +%s)
END_TIME=$(( START_TIME + DURATION ))
SAMPLE_COUNT=0

echo "Starting monitoring with 100ms interval..."
trap 'echo "Monitoring stopped"; exit 0' INT

while [ $(date +%s) -lt $END_TIME ]; do
    if check_container; then
        stats=$(docker stats --no-stream --format "{{.CPUPerc}},{{.MemUsage}},{{.MemPerc}}" $CONTAINER_NAME 2>/dev/null)
        
        if [ $? -eq 0 ] && [ ! -z "$stats" ]; then
            timestamp=$(get_timestamp)
            cpu_percent=$(echo $stats | cut -d',' -f1 | sed 's/%//')
            mem_usage=$(echo $stats | cut -d',' -f2 | awk '{print $1}')
            mem_limit=$(echo $stats | cut -d',' -f2 | awk '{print $3}')
            mem_percent=$(echo $stats | cut -d',' -f3 | sed 's/%//')
            
            echo "$timestamp,$cpu_percent,$mem_usage,$mem_limit,$mem_percent" >> $OUTPUT_FILE
            SAMPLE_COUNT=$((SAMPLE_COUNT + 1))
            
            if [ $((SAMPLE_COUNT % 100)) -eq 0 ]; then
                elapsed=$(( $(date +%s) - START_TIME ))
                echo "Progress: ${elapsed}s/${DURATION}s, Samples: $SAMPLE_COUNT"
            fi
        fi
    else
        echo "Container disappeared, stopping"
        break
    fi
    
done

echo "Monitoring completed. $SAMPLE_COUNT samples saved to $OUTPUT_FILE"
echo "Expected samples: $((DURATION * 10)), Actual: $SAMPLE_COUNT"