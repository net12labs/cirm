module github.com/net12labs/cirm/astro-host

go 1.24.1

require (
	github.com/mattn/go-sqlite3 v1.14.32 // indirect
	github.com/net12labs/cirm/dolly v0.0.0-00010101000000-000000000000 // indirect
)

require github.com/net12labs/cirm/dali v0.0.0-20251116081312-f45cc7e2572e

require (
	github.com/net12labs/cirm/astro-dom v0.0.0-00010101000000-000000000000
	github.com/net12labs/cirm/astro-site v0.0.0-00010101000000-000000000000
	github.com/net12labs/cirm/mali v0.0.0-20251119101424-f011b0bcca3e
	github.com/net12labs/cirm/ops v0.0.0-00010101000000-000000000000
)

replace (
	github.com/net12labs/cirm/astro-dom => ../astro-dom
	github.com/net12labs/cirm/astro-site => ../astro-site
	github.com/net12labs/cirm/dali => ../dali
	github.com/net12labs/cirm/dolly => ../dolly
	github.com/net12labs/cirm/mali => ../mali
	github.com/net12labs/cirm/ops => ../ops

)
