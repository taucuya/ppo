#!/bin/bash
echo "=== Load Balancer Test ==="

backend_8080=0
backend_8081=0
backend_8082=0
unknown=0

echo "Testing 100 requests to /health endpoint:"
for i in {1..100}; do
    response=$(curl -s http://localhost/health 2>/dev/null)
    
    if echo "$response" | grep -q '"port":"8080"'; then
        backend="8080"
        ((backend_8080++))
        echo "Request $i: backend 8080"
    elif echo "$response" | grep -q '"port":"8081"'; then
        backend="8081"
        ((backend_8081++))
        echo "Request $i: backend 8081"
    elif echo "$response" | grep -q '"port":"8082"'; then
        backend="8082"
        ((backend_8082++))
        echo "Request $i: backend 8082"
    else
        backend="unknown"
        ((unknown++))
        echo "Request $i: unknown - $response"
    fi
    
    sleep 0.1
done

echo ""
echo "=== Results ==="
echo "Backend 8080 (weight 2): $backend_8080 requests"
echo "Backend 8081 (weight 1): $backend_8081 requests" 
echo "Backend 8082 (weight 1): $backend_8082 requests"
echo "Unknown: $unknown requests"

total=$(( backend_8080 + backend_8081 + backend_8082 ))

if [ $total -gt 0 ]; then
    echo ""
    echo "Expected ratio: 2:1:1 (50%:25%:25%)"
    pct_8080=$(( backend_8080 * 100 / total ))
    pct_8081=$(( backend_8081 * 100 / total ))
    pct_8082=$(( backend_8082 * 100 / total ))
    echo "Actual ratio: ${pct_8080}% : ${pct_8081}% : ${pct_8082}%"
fi