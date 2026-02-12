#!/usr/bin/env bash
#
# Load test for the singleflight service.
# Sends concurrent requests in waves to demonstrate:
#   - Singleflight deduplication (many concurrent requests for the same ID)
#   - Cache hits vs misses (first wave = miss, subsequent = hit until TTL expires)
#   - HTTP latency under load
#
# Usage: ./loadtest.sh
# Requirements: curl, bash

BASE_URL="http://localhost:8080/template-details"
TEMPLATE_IDS=($(seq 1 100))
CONCURRENCY=200       # requests per burst
WAVES=25              # number of burst waves
PAUSE_BETWEEN=1       # seconds between waves

total_requests=0
total_errors=0

echo "============================================"
echo " Singleflight Load Test"
echo "============================================"
echo " Concurrency per wave : $CONCURRENCY"
echo " Waves                : $WAVES"
echo " Pause between waves  : ${PAUSE_BETWEEN}s"
echo " Template IDs         : ${TEMPLATE_IDS[*]}"
echo "============================================"
echo ""

for wave in $(seq 1 $WAVES); do
    # Pick a random template ID -- all concurrent requests hit the SAME ID
    # to maximize singleflight deduplication
    id=${TEMPLATE_IDS[$((RANDOM % ${#TEMPLATE_IDS[@]}))]}

    echo "--- Wave $wave/$WAVES: $CONCURRENCY concurrent requests for id=$id ---"

    pids=()
    errors=0
    start=$(date +%s%N)

    for i in $(seq 1 $CONCURRENCY); do
        (
            status=$(curl -s -o /dev/null -w "%{http_code}" "$BASE_URL?id=$id")
            if [ "$status" != "200" ]; then
                exit 1
            fi
        ) &
        pids+=($!)
    done

    # Wait and count errors
    for pid in "${pids[@]}"; do
        if ! wait "$pid"; then
            ((errors++))
        fi
    done

    end=$(date +%s%N)
    elapsed_ms=$(( (end - start) / 1000000 ))

    total_requests=$((total_requests + CONCURRENCY))
    total_errors=$((total_errors + errors))

    echo "  Completed in ${elapsed_ms}ms | Errors: $errors/$CONCURRENCY"

    # Also sprinkle in some requests for random IDs (spread across cache)
    for spread_id in "${TEMPLATE_IDS[@]}"; do
        curl -s -o /dev/null "$BASE_URL?id=$spread_id" &
    done
    wait

    total_requests=$((total_requests + ${#TEMPLATE_IDS[@]}))

    if [ "$wave" -lt "$WAVES" ]; then
        sleep $PAUSE_BETWEEN
    fi
done

echo ""
echo "============================================"
echo " Load Test Complete"
echo "============================================"
echo " Total requests : $total_requests"
echo " Total errors   : $total_errors"
echo " Success rate   : $(( (total_requests - total_errors) * 100 / total_requests ))%"
echo "============================================"
echo ""
echo "Check Grafana at http://localhost:3000 (admin/admin)"
echo "Set time range to 'Last 5 minutes' and refresh."
