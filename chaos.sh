#!/bin/bash

TOXI_API="http://localhost:8474"

echo "=== Setting up Toxiproxy ==="

curl -s -X POST "$TOXI_API/proxies" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "primary-api",
    "listen": "0.0.0.0:8091",
    "upstream": "primary-api:8081",
    "enabled": true
  }'

echo ""
echo "=== Injecting chaos: 500ms latency + 20% packet loss ==="

curl -s -X POST "$TOXI_API/proxies/primary-api/toxics" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "latency-500ms",
    "type": "latency",
    "stream": "downstream",
    "toxicity": 1.0,
    "attributes": {
      "latency": 500,
      "jitter": 50
    }
  }'

curl -s -X POST "$TOXI_API/proxies/primary-api/toxics" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "packet-loss-20pct",
    "type": "bandwidth",
    "stream": "