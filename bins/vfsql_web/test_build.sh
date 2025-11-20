#!/bin/bash
# filepath: /home/lxk/Desktop/cirm/bins/vfsql_web/test_build.sh

cd /home/lxk/Desktop/cirm/bins/vfsql_web

echo "=== Testing build ==="
go build -o vfsql_web . 2>&1

if [ $? -eq 0 ]; then
    echo "✓ Build successful"
    echo ""
    echo "Starting server..."
    ./vfsql_web
else
    echo "✗ Build failed"
    echo ""
    echo "Please fix the compilation errors above"
fi
