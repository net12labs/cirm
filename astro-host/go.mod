module github.com/net12labs/cirm/astro-host //china-ip-routes-maker

go 1.24.1

require (
	github.com/mattn/go-sqlite3 v1.14.32 // indirect
	github.com/net12labs/cirm/dolly v0.0.0-00010101000000-000000000000 // indirect
)

require github.com/net12labs/cirm/dali v0.0.0-20251116081312-f45cc7e2572e

require (
	github.com/net12labs/cirm/astro-dom v0.0.0-00010101000000-000000000000
	github.com/net12labs/cirm/astro-site/agent-client-web v0.0.0-00010101000000-000000000000
	github.com/net12labs/cirm/astro-site/agent-web v0.0.0-00010101000000-000000000000
	github.com/net12labs/cirm/astro-site/agent-web-api v0.0.0-00010101000000-000000000000
	github.com/net12labs/cirm/astro-site/ai-agent-web v0.0.0-00010101000000-000000000000
	github.com/net12labs/cirm/astro-site/ai-agent-web-api v0.0.0-00010101000000-000000000000
	github.com/net12labs/cirm/astro-site/ai-agent-web-page v0.0.0-00010101000000-000000000000
	github.com/net12labs/cirm/astro-site/website-web v0.0.0-00010101000000-000000000000
	github.com/net12labs/cirm/astro-site/website-web-api v0.0.0-00010101000000-000000000000
	github.com/net12labs/cirm/astro-site/website-web-page v0.0.0-00010101000000-000000000000
	github.com/net12labs/cirm/mali v0.0.0-20251119101424-f011b0bcca3e
	github.com/net12labs/cirm/ops v0.0.0-00010101000000-000000000000
)

replace (
	github.com/net12labs/cirm/astro-dom => ../astro-dom
	github.com/net12labs/cirm/astro-site/agent-client-web => ../astro-site/agent-client-web
	github.com/net12labs/cirm/astro-site/agent-web => ../astro-site/agent-web
	github.com/net12labs/cirm/astro-site/agent-web-api => ../astro-site/agent-web-api
	github.com/net12labs/cirm/astro-site/ai-agent-web => ../astro-site/ai-agent-web
	github.com/net12labs/cirm/astro-site/ai-agent-web-api => ../astro-site/ai-agent-web-api
	github.com/net12labs/cirm/astro-site/ai-agent-web-page => ../astro-site/ai-agent-web-page
	github.com/net12labs/cirm/astro-site/website-web => ../astro-site/website-web
	github.com/net12labs/cirm/astro-site/website-web-api => ../astro-site/website-web-api
	github.com/net12labs/cirm/astro-site/website-web-page => ../astro-site/website-web-page
	github.com/net12labs/cirm/dali => ../dali
	github.com/net12labs/cirm/dolly => ../dolly
	github.com/net12labs/cirm/mali => ../mali
	github.com/net12labs/cirm/ops => ../ops

)
