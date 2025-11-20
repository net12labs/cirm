#!/bin/bash
# filepath: /home/lxk/Desktop/cirm/bins/vfsql_web/run.sh

# Build the application
echo "Building VFSQL Web Server..."
go build -o vfsql_web .

if [ $? -ne 0 ]; then
    echo "Build failed!"
    exit 1
fi

# Create web symlink if needed
if [ ! -e "web" ] && [ -d "../vfsql/web" ]; then
    echo "Creating symlink to web directory..."
    ln -sf ../vfsql/web web
fi

# Run the server
echo ""
echo "Starting VFSQL Web Server..."
echo "================================"
./vfsql_web "$@"
