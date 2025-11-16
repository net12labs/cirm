module github.com/net12labs/cirm/ai-agent-web-api

go 1.21

toolchain go1.24.1

replace github.com/net12labs/cirm/dali => ../dali

replace github.com/net12labs/cirm/mali => ../mali

require github.com/net12labs/cirm/dali v0.0.0-00010101000000-000000000000

require github.com/net12labs/cirm/mali v0.0.0-20251116070540-b170edf05994 // indirect
