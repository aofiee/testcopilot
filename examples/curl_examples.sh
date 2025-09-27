#!/bin/bash

# HMAC Service API Examples
# Make sure the service is running on localhost:8080

BASE_URL="http://localhost:8080"

echo "=== HMAC Service API Examples ==="
echo

# Health check
echo "1. Health Check:"
curl -s -X GET "${BASE_URL}/health" | python3 -m json.tool
echo
echo

# Encrypt (Generate HMAC)
echo "2. Generate HMAC:"
HMAC_RESPONSE=$(curl -s -X POST "${BASE_URL}/hmac/encrypt" \
  -H "Content-Type: application/json" \
  -d '{"message":"hello world","key":"mysecret"}')

echo $HMAC_RESPONSE | python3 -m json.tool
echo

# Extract HMAC for verification
HMAC_VALUE=$(echo $HMAC_RESPONSE | python3 -c "import sys, json; print(json.load(sys.stdin)['hmac'])")
echo "Generated HMAC: $HMAC_VALUE"
echo

# Decrypt (Verify HMAC) - Valid case
echo "3. Verify HMAC (Valid):"
curl -s -X POST "${BASE_URL}/hmac/decrypt" \
  -H "Content-Type: application/json" \
  -d "{\"message\":\"hello world\",\"key\":\"mysecret\",\"hmac\":\"${HMAC_VALUE}\"}" | python3 -m json.tool
echo
echo

# Decrypt (Verify HMAC) - Invalid case
echo "4. Verify HMAC (Invalid):"
curl -s -X POST "${BASE_URL}/hmac/decrypt" \
  -H "Content-Type: application/json" \
  -d '{"message":"hello world","key":"mysecret","hmac":"invalid-hmac"}' | python3 -m json.tool
echo
echo

# Error cases
echo "5. Error Cases:"
echo

echo "5a. Missing message:"
curl -s -X POST "${BASE_URL}/hmac/encrypt" \
  -H "Content-Type: application/json" \
  -d '{"key":"mysecret"}' | python3 -m json.tool
echo
echo

echo "5b. Missing key:"
curl -s -X POST "${BASE_URL}/hmac/encrypt" \
  -H "Content-Type: application/json" \
  -d '{"message":"hello world"}' | python3 -m json.tool
echo
echo

echo "5c. Invalid JSON:"
curl -s -X POST "${BASE_URL}/hmac/encrypt" \
  -H "Content-Type: application/json" \
  -d 'invalid json' | python3 -m json.tool
echo
echo

echo "5d. Non-existent endpoint:"
curl -s -X GET "${BASE_URL}/nonexistent" | python3 -m json.tool
echo

echo "=== Examples Complete ==="