module github.com/net12labs/cirm/service-daemon //china-ip-routes-maker

go 1.21

require github.com/mattn/go-sqlite3 v1.14.32

require github.com/net12labs/cirm/dali v0.0.0-20251114182153-000f7e221ad5

require github.com/net12labs/cirm/client-web/admin v0.0.0-20251114191024-95d4142052d4

require github.com/net12labs/cirm/client-web/user v0.0.0-20251114170627-3a58839cfb1a

require github.com/net12labs/cirm/client-web/provider v0.0.0-20251114191024-95d4142052d4

require github.com/net12labs/cirm/client-web/root v0.0.0-20251114191024-95d4142052d4
require github.com/net12labs/cirm/client-web/platform v0.0.0-20251114191024-95d4142052d4

replace (
	github.com/net12labs/cirm/client-web/admin => ../client-web/admin
	github.com/net12labs/cirm/client-web/provider => ../client-web/provider
	github.com/net12labs/cirm/client-web/root => ../client-web/root
	github.com/net12labs/cirm/client-web/user => ../client-web/user
	github.com/net12labs/cirm/client-web/platform => ../client-web/platform
	github.com/net12labs/cirm/dali => ../dali
)
