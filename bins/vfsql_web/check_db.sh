#!/bin/bash
# filepath: /home/lxk/Desktop/cirm/bins/vfsql_web/check_db.sh

# Check what's actually in the database

echo "=== Checking VFSQL Database ==="
echo

if [ ! -f "vfsql.db" ]; then
    echo "vfsql.db not found!"
    exit 1
fi

echo "Files in database:"
sqlite3 vfsql.db "SELECT name, size FROM inodes WHERE type='file';"

echo
echo "Storage entries:"
sqlite3 vfsql.db "SELECT inode_id, length(content) as content_length FROM storage;"

echo
echo "Files with storage:"
sqlite3 vfsql.db "
SELECT i.name, i.size as inode_size, length(s.content) as storage_size
FROM inodes i
LEFT JOIN storage s ON i.id = s.inode_id
WHERE i.type = 'file';
"
