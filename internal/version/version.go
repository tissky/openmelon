// Package version exposes the build version of openmelon.
//
// The default value is the development sentinel "dev". Release builds
// override it with -ldflags:
//
//	go build -ldflags "-X github.com/eight-acres-lab/openmelon/internal/version.Version=v0.2.0" ./cmd/openmelon
//
// Release scripts read the tag from `git describe --tags` and pass it in.
package version

// Version is set at build time. "dev" indicates an untagged build.
var Version = "dev"
