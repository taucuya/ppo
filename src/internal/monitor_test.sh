#!/bin/bash

STATS_FILE="test_stats.log"
echo "" > $STATS_FILE
echo "Время запуска: $(date)" >> $STATS_FILE

monitor_processes() {
    while true; do
        COUNT=$(ps aux | grep -E "go.test|go-test" | grep -v grep | wc -l)
        echo "[$(date +%H:%M:%S)] Процессов: $COUNT" >> $STATS_FILE
        sleep 0.1
    done
}

monitor_processes &
MONITOR_PID=$!

sleep 0.5

go test -v ./... -shuffle=on -p 3 -parallel 2
TEST_RESULT=$?

kill $MONITOR_PID 2>/dev/null
wait $MONITOR_PID 2>/dev/null

echo "Статистика процессов:"
cat $STATS_FILE

echo "Тесты завершены с кодом: $TEST_RESULT"
exit $TEST_RESULT