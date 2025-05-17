// Package buildinfo provides variables that store build-time information
// such as version, commit hash, and build date. These variables are typically
// populated by LDFLAGS during the Go build process.
// file: internal/buildinfo/buildinfo.go
package buildinfo

// These variables are populated by LDFLAGS during the build process
var (
	Version    string = "dev"
	CommitHash string = "unknown"
	BuildDate  string = "unknown"
)
