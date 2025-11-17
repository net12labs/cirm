module github.com/net12labs/cirm/service-daemon //china-ip-routes-maker

go 1.21

require github.com/mattn/go-sqlite3 v1.14.32 // indirect

require github.com/net12labs/cirm/dali v0.0.0-20251116081312-f45cc7e2572e

require github.com/net12labs/cirm/site-client-web/admin v0.0.0-20251114191024-95d4142052d4

require github.com/net12labs/cirm/site-client-web/provider v0.0.0-20251114191024-95d4142052d4

require (
	github.com/net12labs/cirm/ai-agent-web v0.0.0-00010101000000-000000000000
	github.com/net12labs/cirm/ai-agent-web-api v0.0.0-00010101000000-000000000000
	github.com/net12labs/cirm/ai-agent-web-client v0.0.0-00010101000000-000000000000
)

require (
	github.com/net12labs/cirm/agent-client-web v0.0.0-00010101000000-000000000000
	github.com/net12labs/cirm/agent-web v0.0.0-00010101000000-000000000000
	github.com/net12labs/cirm/agent-web-api v0.0.0-00010101000000-000000000000
	github.com/net12labs/cirm/mali v0.0.0-20251116215737-9700c55ab92c
	github.com/net12labs/cirm/site-client-web/consumer v0.0.0-00010101000000-000000000000
	github.com/net12labs/cirm/site-client-web/platform v0.0.0-00010101000000-000000000000
	github.com/net12labs/cirm/site-client-web/site v0.0.0-00010101000000-000000000000
	github.com/net12labs/cirm/site-web v0.0.0-00010101000000-000000000000
	github.com/net12labs/cirm/site-web-api v0.0.0-20251116173359-e56a45cf349c
)

replace (
	github.com/net12labs/cirm/agent-client-web => ../agent-client-web
	github.com/net12labs/cirm/agent-web => ../agent-web
	github.com/net12labs/cirm/agent-web-api => ../agent-web-api
	github.com/net12labs/cirm/ai-agent-web => ../ai-agent-web
	github.com/net12labs/cirm/ai-agent-web-api => ../ai-agent-web-api
	github.com/net12labs/cirm/ai-agent-web-client => ../ai-agent-web-client
	github.com/net12labs/cirm/dali => ../dali
	github.com/net12labs/cirm/mali => ../mali
	github.com/net12labs/cirm/site-client-web/admin => ../site-client-web/admin
	github.com/net12labs/cirm/site-client-web/consumer => ../site-client-web/consumer
	github.com/net12labs/cirm/site-client-web/platform => ../site-client-web/platform
	github.com/net12labs/cirm/site-client-web/provider => ../site-client-web/provider
	github.com/net12labs/cirm/site-client-web/site => ../site-client-web/site
	github.com/net12labs/cirm/site-web => ../site-web
	github.com/net12labs/cirm/site-web-api => ../site-web-api
)
