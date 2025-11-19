module github.com/net12labs/cirm/astro-site/ai-agent-web

go 1.24.1

require (
	github.com/net12labs/cirm/dali v0.0.0-20251116081312-f45cc7e2572e
	github.com/net12labs/cirm/ops v0.0.0-00010101000000-000000000000
)

require github.com/net12labs/cirm/mali v0.0.0-20251119101424-f011b0bcca3e // indirect

replace github.com/net12labs/cirm/dali => ../../dali

replace github.com/net12labs/cirm/ops => ../../ops
