#!/bin/bash

echo "=== Bitespeed Identity Reconciliation API Tests ==="
echo ""

BASE_URL="http://localhost:8080"

# Test 1: Health check
echo "1. Testing health endpoint..."
curl -s -o /dev/null -w "Status: %{http_code}\n" $BASE_URL/health
echo ""

# Test 2: Create first contact
echo "2. Creating first contact..."
curl -X POST $BASE_URL/identify \
  -H "Content-Type: application/json" \
  -d '{
    "email": "lorraine@hillvalley.edu",
    "phoneNumber": "123456"
  }' | jq '.'
echo ""

# Test 3: Create contact with same email, different phone
echo "3. Creating contact with same email, different phone..."
curl -X POST $BASE_URL/identify \
  -H "Content-Type: application/json" \
  -d '{
    "email": "lorraine@hillvalley.edu",
    "phoneNumber": "789012"
  }' | jq '.'
echo ""

# Test 4: Create contact with new email, same phone as first
echo "4. Creating contact with new email, same phone as first..."
curl -X POST $BASE_URL/identify \
  -H "Content-Type: application/json" \
  -d '{
    "email": "mcfly@hillvalley.edu",
    "phoneNumber": "123456"
  }' | jq '.'
echo ""

# Test 5: Merge separate primary contacts
echo "5. Testing merge of separate primary contacts..."
curl -X POST $BASE_URL/identify \
  -H "Content-Type: application/json" \
  -d '{
    "email": "george@hillvalley.edu",
    "phoneNumber": "919191"
  }' | jq '.'
echo ""

curl -X POST $BASE_URL/identify \
  -H "Content-Type: application/json" \
  -d '{
    "email": "biffsucks@hillvalley.edu",
    "phoneNumber": "717171"
  }' | jq '.'
echo ""

curl -X POST $BASE_URL/identify \
  -H "Content-Type: application/json" \
  -d '{
    "email": "george@hillvalley.edu",
    "phoneNumber": "717171"
  }' | jq '.'
echo ""

# Test 6: Error cases
echo "6. Testing error cases..."
echo "6a. Empty request:"
curl -X POST $BASE_URL/identify \
  -H "Content-Type: application/json" \
  -d '{}' | jq '.'
echo ""

echo "6b. Invalid JSON:"
curl -X POST $BASE_URL/identify \
  -H "Content-Type: application/json" \
  -d '{invalid json}' 
echo ""

echo "6c. Wrong HTTP method:"
curl -X GET $BASE_URL/identify
echo ""

echo "=== Tests completed ==="