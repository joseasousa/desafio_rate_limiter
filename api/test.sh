#!/bin/bash

url="http://localhost:8080/"
token="token1"
iterations=11

for ((i=1; i<=iterations; i++))
do
    response=$(curl -s -w "\n%{http_code}" -H "API_KEY: $token" $url)
    body=$(echo "$response" | sed '$d')
    status_code=$(echo "$response" | tail -n 1)
    echo "Request $i: Status Code: $status_code - ResponseBody: $body"
    sleep 0.1  # espera 100ms entre as requisições
done