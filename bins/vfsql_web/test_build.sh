#!/bin/bash
# filepath: /home/lxk/Desktop/cirm/bins/vfsql_web/test_build.sh

cd /home/lxk/Desktop/cirm/bins/vfsql_web

echo "=== Testing Build ==="
go build -o vfsql_web . 2>&1 | head -20

if [ $? -eq 0 ]; then
    echo ""
    echo "✓ Build successful!"
else
    echo ""
    echo "✗ Build failed - showing first 20 errors"
    echo ""
    echo "Common fixes:"
    echo "1. Check for duplicate function declarations"
    echo "2. Check for misplaced code (comments or partial functions)"
    echo "3. Look for 'illegal label' or 'missing comma' errors"
fi
