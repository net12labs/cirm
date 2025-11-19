module github.com/net12labs/cirm/astro-pack //china-ip-routes-maker

go 1.24.1

require github.com/mattn/go-sqlite3 v1.14.32 // indirect

require github.com/net12labs/cirm/dali v0.0.0-20251116081312-f45cc7e2572e

require (
	github.com/net12labs/cirm/agent-client-web v0.0.0-00010101000000-000000000000
	github.com/net12labs/cirm/agent-web v0.0.0-00010101000000-000000000000
	github.com/net12labs/cirm/agent-web-api v0.0.0-00010101000000-000000000000
	github.com/net12labs/cirm/ai-agent-web v0.0.0-00010101000000-000000000000
	github.com/net12labs/cirm/ai-agent-web-api v0.0.0-00010101000000-000000000000
	github.com/net12labs/cirm/ai-agent-web-page v0.0.0-00010101000000-000000000000
	github.com/net12labs/cirm/dolly v0.0.0-00010101000000-000000000000
	github.com/net12labs/cirm/mali v0.0.0-20251119101424-f011b0bcca3e
	github.com/net12labs/cirm/ops v0.0.0-00010101000000-000000000000
	github.com/net12labs/cirm/website-web v0.0.0-00010101000000-000000000000
	github.com/net12labs/cirm/website-web-api v0.0.0-20251116173359-e56a45cf349c
	github.com/net12labs/cirm/website-web-page v0.0.0-00010101000000-000000000000
)

replace (
	github.com/net12labs/cirm/agent-client-web => ../agent-client-web
	github.com/net12labs/cirm/agent-web => ../agent-web
	github.com/net12labs/cirm/agent-web-api => ../agent-web-api
	github.com/net12labs/cirm/ai-agent-web => ../ai-agent-web
	github.com/net12labs/cirm/ai-agent-web-api => ../ai-agent-web-api
	github.com/net12labs/cirm/ai-agent-web-page => ../ai-agent-web-page
	github.com/net12labs/cirm/astro-dom => ../astro-dom
	github.com/net12labs/cirm/dali => ../dali
	github.com/net12labs/cirm/dolly => ../dolly
	github.com/net12labs/cirm/mali => ../mali
	github.com/net12labs/cirm/ops => ../ops
	github.com/net12labs/cirm/website-web => ../website-web
	github.com/net12labs/cirm/website-web-api => ../website-web-api
	github.com/net12labs/cirm/website-web-page => ../website-web-page

	github.com/net12labs/cirm/astro-dom => ../astro-dom

)
