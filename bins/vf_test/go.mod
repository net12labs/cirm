module github.com/lxk/cirm/bins/vf_test

go 1.21

require github.com/lxk/cirm/bins/vfsql v0.0.0

require github.com/mattn/go-sqlite3 v1.14.18 // indirect

replace github.com/lxk/cirm/bins/vfsql => ../vfsql
