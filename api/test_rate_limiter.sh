#!/bin/bash

url="http://localhost:8080/"

function make_requests {
    local token="$1"
    local description="$2"
    local max_requests="$3"
    local iterations=1
    local client_ip="192.168.0.1"

    echo "Testing scenario: $description"

    while true; do
        if [ -n "$token" ]; then
            response=$(curl -s -D - -H "API_KEY: $token" -H "X-Real-IP: $client_ip" $url)
        else
            response=$(curl -s -D - -H "X-Real-IP: $client_ip" $url)
        fi

        status_code=$(echo "$response" | grep HTTP | awk '{print $2}')
        body=$(echo "$response" | sed -n '/^\r$/,$p' | tail -n +2)

        echo "Request $iterations:"
        echo "$body"
        echo "Status Code: $status_code"
        echo "-------------------------"

        if [ "$status_code" -eq 429 ]; then
            echo "Received 429 Too Many Requests"
            break
        fi

        if [ "$iterations" -ge "$max_requests" ]; then
            echo "Max requests reached without 429"
            break
        fi

        iterations=$((iterations + 1))
    done
}

# Scenario 1: Known token
make_requests "token1" "Known token" 11

# Scenario 2: Unknown token
make_requests "unknown-token" "Unknown token" 3

# Scenario 3: No token
make_requests "" "No token" 3