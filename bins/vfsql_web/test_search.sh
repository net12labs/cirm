#!/bin/bash
# filepath: /home/lxk/Desktop/cirm/bins/vfsql_web/test_search.sh

echo "=== Testing VFSQL Search Functionality ==="
echo

# Test 1: Pattern search
echo "Test 1: Search by pattern (*.txt)"
curl -s "http://localhost:8080/api/search?pattern=*.txt&path=/&recursive=true" | jq '.'
echo

# Test 2: Search all files
echo "Test 2: Search all files"
curl -s "http://localhost:8080/api/search?pattern=*&path=/&recursive=true" | jq '.'
echo

# Test 3: Search by tags
echo "Test 3: Search by tags"
curl -s "http://localhost:8080/api/search?tags=test&path=/&recursive=true" | jq '.'
echo

# Test 4: Search by description
echo "Test 4: Search by description"
curl -s "http://localhost:8080/api/search?description=test&path=/&recursive=true" | jq '.'
echo

# Test 5: Non-recursive search
echo "Test 5: Non-recursive search"
curl -s "http://localhost:8080/api/search?pattern=*&path=/&recursive=false" | jq '.'
echo

echo "=== Tests Complete ==="
