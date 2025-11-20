#!/bin/bash
# Fix return statements in CreateFilesystem function

cd /home/lxk/Desktop/cirm/bins/vfsql

# Backup first
cp volume.go volume.go.backup

# Replace "return nil, " with "return " in the CreateFilesystem function (lines 175-223)
sed -i '175,223s/return nil, /return /' volume.go

echo "Fixed return statements in volume.go"
echo "Original backed up as volume.go.backup"
