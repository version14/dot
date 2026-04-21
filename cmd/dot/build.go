package main

// buildVersion is set at build time via:
//
//	go build -ldflags "-X main.buildVersion=v0.1.0" ./cmd/dot
//
// Falls back to "dev" for local builds.
var buildVersion = "dev"
