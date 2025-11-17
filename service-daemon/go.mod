module github.com/net12labs/cirm/service-daemon //china-ip-routes-maker

go 1.24.1

require github.com/mattn/go-sqlite3 v1.14.32 // indirect

require github.com/net12labs/cirm/dali v0.0.0-20251116081312-f45cc7e2572e

require (
	github.com/net12labs/cirm/ai-agent-web v0.0.0-00010101000000-000000000000
	github.com/net12labs/cirm/ai-agent-web-api v0.0.0-00010101000000-000000000000
	github.com/net12labs/cirm/ai-agent-web-client v0.0.0-00010101000000-000000000000
	github.com/net12labs/cirm/site-client-web v0.0.0-00010101000000-000000000000
)

require (
	github.com/net12labs/cirm/agent-client-web v0.0.0-00010101000000-000000000000
	github.com/net12labs/cirm/agent-web v0.0.0-00010101000000-000000000000
	github.com/net12labs/cirm/agent-web-api v0.0.0-00010101000000-000000000000
	github.com/net12labs/cirm/dolly v0.0.0-00010101000000-000000000000
	github.com/net12labs/cirm/domain-web v0.0.0-00010101000000-000000000000
	github.com/net12labs/cirm/mali v0.0.0-20251117081457-1fef358e292e
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
	github.com/net12labs/cirm/dolly => ../dolly
	github.com/net12labs/cirm/domain-web => ../domain-web
	github.com/net12labs/cirm/mali => ../mali
	github.com/net12labs/cirm/site-client-web => ../site-client-web
	github.com/net12labs/cirm/site-web => ../site-web
	github.com/net12labs/cirm/site-web-api => ../site-web-api

)
