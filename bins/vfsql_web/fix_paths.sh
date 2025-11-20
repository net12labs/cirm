#!/bin/bash
# filepath: /home/lxk/Desktop/cirm/bins/vfsql_web/fix_paths.sh

# This script patches the server.go to normalize search result paths

echo "Fixing double slash issue in search results..."

# Find the handleSearch function and add path normalization
cd /home/lxk/Desktop/cirm/bins/vfsql/api

# Create a backup
cp server.go server.go.backup

# Add path normalization after search results
# Look for "s.jsonResponse(w, results)" in handleSearch and add normalization before it

echo "Manual fix needed:"
echo "In /home/lxk/Desktop/cirm/bins/vfsql/api/server.go"
echo "In the handleSearch function, before 's.jsonResponse(w, results)', add:"
echo ""
echo "    // Normalize paths - remove double slashes"
echo "    for i := range results {"
echo "        results[i] = strings.ReplaceAll(results[i], \"//\", \"/\")"
echo "    }"
echo ""
echo "Or the issue might be in VFS Search - check how paths are built there"
