module github.com/net12labs/cirm/service-daemon //china-ip-routes-maker

go 1.21

require github.com/mattn/go-sqlite3 v1.14.32 // indirect

require github.com/net12labs/cirm/dali v0.0.0-20251115100224-b432eab83657

require github.com/net12labs/cirm/client-web/admin v0.0.0-20251114191024-95d4142052d4

require github.com/net12labs/cirm/client-web/user v0.0.0-20251114170627-3a58839cfb1a

require github.com/net12labs/cirm/client-web/provider v0.0.0-20251114191024-95d4142052d4

require github.com/net12labs/cirm/client-web/root v0.0.0-20251114191024-95d4142052d4

require (
	github.com/net12labs/cirm/api-web v0.0.0-00010101000000-000000000000
	github.com/net12labs/cirm/client-web/platform v0.0.0-20251114191024-95d4142052d4
)

require github.com/net12labs/cirm/mali v0.0.0-20251116070540-b170edf05994 // indirect

replace (
	github.com/net12labs/cirm/api-web => ../api-web
	github.com/net12labs/cirm/bin => ../bin
	github.com/net12labs/cirm/client-web/admin => ../client-web/admin
	github.com/net12labs/cirm/client-web/platform => ../client-web/platform
	github.com/net12labs/cirm/client-web/provider => ../client-web/provider
	github.com/net12labs/cirm/client-web/root => ../client-web/root
	github.com/net12labs/cirm/client-web/user => ../client-web/user
	github.com/net12labs/cirm/dali => ../dali
)
